package agentcompose

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"connectrpc.com/connect"
	agentcomposev2 "github.com/chaitin/chaitin-cli/products/agentcompose/gen/agentcompose/v2"
	"github.com/spf13/cobra"
)

const eventLogPrefix = "evt_"

// validateEventLogTarget rejects values that are not full event-bus event
// ids. Prefix matching is intentionally unsupported for --event.
func validateEventLogTarget(eventID string, jsonOutput bool) error {
	if !strings.HasPrefix(eventID, eventLogPrefix) || strings.TrimSpace(strings.TrimPrefix(eventID, eventLogPrefix)) == "" {
		return usageError(fmt.Sprintf("--event requires a full event id starting with %s; prefix matching is not supported", eventLogPrefix), jsonOutput)
	}
	return nil
}

// resolveEventLogSchedulerRunIDs fetches GET /api/events/{id}/trace and
// returns the delivery scheduler run ids in delivery order, deduplicated.
// The trace endpoint is a daemon REST route rather than a Connect procedure,
// so the shared bearer-token HTTP client is used directly.
func resolveEventLogSchedulerRunIDs(ctx context.Context, state *commandState, eventID string) ([]string, error) {
	eventID = strings.TrimSpace(eventID)
	endpoint := strings.TrimSuffix(state.options.URL, "/") + "/api/events/" + url.PathEscape(eventID) + "/trace"
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, usageError(fmt.Sprintf("build event trace request for %s: %v", eventID, err), state.options.JSON)
	}
	response, err := state.httpClient().Do(request)
	if err != nil {
		return nil, mapConnectError(err, state.options.URL, state.options.JSON)
	}
	defer func() { _ = response.Body.Close() }()
	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return nil, mapConnectError(readErr, state.options.URL, state.options.JSON)
	}
	switch {
	case response.StatusCode == http.StatusNotFound:
		return nil, newError("not_found", fmt.Sprintf("event %s not found", eventID), exitNotFound, state.options.JSON)
	case response.StatusCode >= 300:
		return nil, newError("service_unavailable", fmt.Sprintf("event trace for %s returned status %s", eventID, response.Status), exitNetwork, state.options.JSON)
	}
	var trace struct {
		Runs []struct {
			Delivery struct {
				RunID string `json:"run_id"`
			} `json:"delivery"`
		} `json:"runs"`
	}
	if err := json.Unmarshal(body, &trace); err != nil {
		return nil, newError("execution_failed", fmt.Sprintf("parse event %s trace response: %v", eventID, err), exitGeneral, state.options.JSON)
	}
	runIDs := make([]string, 0, len(trace.Runs))
	seen := make(map[string]struct{}, len(trace.Runs))
	for _, run := range trace.Runs {
		runID := strings.TrimSpace(run.Delivery.RunID)
		if runID == "" {
			continue
		}
		if _, ok := seen[runID]; ok {
			continue
		}
		seen[runID] = struct{}{}
		runIDs = append(runIDs, runID)
	}
	return runIDs, nil
}

// resolveEventLogRuns expands scheduler run ids into the project's agent runs
// that recorded them, preserving scheduler run order.
func resolveEventLogRuns(ctx context.Context, state *commandState, projectID string, schedulerRunIDs []string) ([]*agentcomposev2.RunSummary, error) {
	client := state.clients().run
	runs := make([]*agentcomposev2.RunSummary, 0, len(schedulerRunIDs))
	seen := make(map[string]struct{}, len(schedulerRunIDs))
	for _, schedulerRunID := range schedulerRunIDs {
		resp, err := client.ListRuns(ctx, connect.NewRequest(&agentcomposev2.ListRunsRequest{
			ProjectId:      strings.TrimSpace(projectID),
			SchedulerRunId: schedulerRunID,
			Limit:          200,
		}))
		if err != nil {
			return nil, mapConnectError(err, state.options.URL, state.options.JSON)
		}
		for _, run := range resp.Msg.GetRuns() {
			runID := run.GetRunId()
			if _, ok := seen[runID]; ok {
				continue
			}
			seen[runID] = struct{}{}
			runs = append(runs, run)
		}
	}
	return runs, nil
}

// eventLogReplayRuns narrows resolved runs to the requested agent, if any.
func eventLogReplayRuns(runs []*agentcomposev2.RunSummary, agentName string) []*agentcomposev2.RunSummary {
	agentName = strings.TrimSpace(agentName)
	if agentName == "" {
		return runs
	}
	narrowed := make([]*agentcomposev2.RunSummary, 0, len(runs))
	for _, run := range runs {
		if run.GetAgentName() == agentName {
			narrowed = append(narrowed, run)
		}
	}
	return narrowed
}

// executeLogsForEvent resolves the --event target and replays the associated
// agent runs. An event without runs prints a notice (or an empty JSON runs
// array) and succeeds.
func executeLogsForEvent(cmd *cobra.Command, state *commandState, projectID string, options logsOptions) error {
	schedulerRunIDs, err := resolveEventLogSchedulerRunIDs(cmd.Context(), state, options.Event)
	if err != nil {
		return err
	}
	runs := make([]*agentcomposev2.RunSummary, 0)
	if len(schedulerRunIDs) > 0 {
		runs, err = resolveEventLogRuns(cmd.Context(), state, projectID, schedulerRunIDs)
		if err != nil {
			return err
		}
	}
	runs = eventLogReplayRuns(runs, options.Agent)
	if len(runs) == 0 {
		if state.options.JSON {
			return writeJSON(cmd.OutOrStdout(), struct {
				Runs []*agentcomposev2.RunSummary `json:"runs"`
			}{Runs: []*agentcomposev2.RunSummary{}})
		}
		_, err := fmt.Fprintf(cmd.ErrOrStderr(), "No runs are associated with event %s\n", strings.TrimSpace(options.Event))
		return err
	}
	for _, run := range runs {
		if err := followOneRun(cmd, state, projectID, run, options); err != nil {
			return err
		}
	}
	return nil
}
