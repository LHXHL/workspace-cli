package agentcompose

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"connectrpc.com/connect"
	agentcomposev2 "github.com/chaitin/chaitin-cli/products/agentcompose/gen/agentcompose/v2"
	agentcomposev2connect "github.com/chaitin/chaitin-cli/products/agentcompose/gen/agentcompose/v2/agentcomposev2connect"
	"google.golang.org/protobuf/proto"
)

func TestValidateEventLogTarget(t *testing.T) {
	if err := validateEventLogTarget("evt_abc123", false); err != nil {
		t.Fatalf("validateEventLogTarget(full id) returned error: %v", err)
	}
	for _, value := range []string{"evt_", "evt", "run-123", "abc"} {
		err := validateEventLogTarget(value, false)
		if err == nil {
			t.Fatalf("validateEventLogTarget(%q) = nil, want usage error", value)
		}
		cliErr, ok := err.(*CLIError)
		if !ok || cliErr.ExitCode() != exitUsage {
			t.Fatalf("validateEventLogTarget(%q) = %#v, want usage error", value, err)
		}
	}
}

func TestEventLogReplayRunsNarrowsByAgent(t *testing.T) {
	runs := []*agentcomposev2.RunSummary{
		{RunId: "run-reviewer", AgentName: "reviewer"},
		{RunId: "run-writer", AgentName: "writer"},
	}
	if got := eventLogReplayRuns(runs, ""); len(got) != 2 {
		t.Fatalf("eventLogReplayRuns without agent = %d runs, want 2", len(got))
	}
	narrowed := eventLogReplayRuns(runs, "writer")
	if len(narrowed) != 1 || narrowed[0].GetRunId() != "run-writer" {
		t.Fatalf("eventLogReplayRuns(writer) = %#v", narrowed)
	}
	if got := eventLogReplayRuns(runs, "missing"); len(got) != 0 {
		t.Fatalf("eventLogReplayRuns(missing) = %#v, want empty", got)
	}
}

// schedulerRunStub filters ListRuns by the scheduler_run_id filter, unlike the
// shared runStub, which ignores request filters.
type schedulerRunStub struct {
	agentcomposev2connect.UnimplementedRunServiceHandler
	mu               sync.Mutex
	runsByScheduler  map[string][]*agentcomposev2.RunSummary
	listRequests     []*agentcomposev2.ListRunsRequest
	followRunOutputs map[string]string
}

func (s *schedulerRunStub) ListRuns(_ context.Context, req *connect.Request[agentcomposev2.ListRunsRequest]) (*connect.Response[agentcomposev2.ListRunsResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.listRequests = append(s.listRequests, proto.Clone(req.Msg).(*agentcomposev2.ListRunsRequest))
	runs := s.runsByScheduler[req.Msg.GetSchedulerRunId()]
	if int(req.Msg.GetLimit()) < len(runs) {
		runs = runs[:req.Msg.GetLimit()]
	}
	return connect.NewResponse(&agentcomposev2.ListRunsResponse{Runs: runs, Total: uint32(len(runs))}), nil
}

func (s *schedulerRunStub) FollowRunLogs(_ context.Context, req *connect.Request[agentcomposev2.FollowRunLogsRequest], stream *connect.ServerStream[agentcomposev2.RunLogChunk]) error {
	output := s.followRunOutputs[req.Msg.GetRunId()]
	for _, line := range strings.SplitAfter(output, "\n") {
		if line == "" {
			continue
		}
		if err := stream.Send(&agentcomposev2.RunLogChunk{Data: line, RunStatus: agentcomposev2.RunStatus_RUN_STATUS_SUCCEEDED}); err != nil {
			return err
		}
	}
	return stream.Send(&agentcomposev2.RunLogChunk{RunStatus: agentcomposev2.RunStatus_RUN_STATUS_SUCCEEDED, IsFinal: true})
}

func newEventLogsTestServer(t *testing.T, run *schedulerRunStub, trace func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	path, handler := agentcomposev2connect.NewProjectServiceHandler(&projectStub{project: fixtureProject()})
	mux.Handle(path, handler)
	path, handler = agentcomposev2connect.NewRunServiceHandler(run)
	mux.Handle(path, handler)
	mux.HandleFunc("/api/events/", trace)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization = %q", got)
		}
		mux.ServeHTTP(w, r)
	}))
}

func eventTraceResponse(t *testing.T, eventID string, schedulerRunIDs ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/events/"+eventID+"/trace" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		runs := make([]map[string]any, 0, len(schedulerRunIDs))
		for _, schedulerRunID := range schedulerRunIDs {
			runs = append(runs, map[string]any{"delivery": map[string]any{"run_id": schedulerRunID, "status": "run_succeeded"}})
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"event": map[string]any{"id": eventID}, "runs": runs}); err != nil {
			t.Fatalf("encode trace response: %v", err)
		}
	}
}

func TestExecuteLogsForEventReplaysRunsAndValidatesSelectors(t *testing.T) {
	run := &schedulerRunStub{
		runsByScheduler: map[string][]*agentcomposev2.RunSummary{
			"sched-run-reviewer": {{RunId: "run-event-reviewer", ProjectId: "project-aaaaaaaaaaaaaaaa", AgentName: "agent", SandboxId: "sandbox-1"}},
			"sched-run-empty":    {},
		},
		followRunOutputs: map[string]string{"run-event-reviewer": "event output\n"},
	}
	server := newEventLogsTestServer(t, run, eventTraceResponse(t, "evt_event_test", "sched-run-reviewer", "sched-run-empty", "sched-run-reviewer"))

	stdout, _, err := executeCommand(t, server.URL, false, "logs", "--event", "evt_event_test")
	if err != nil {
		t.Fatalf("logs --event returned error: %v", err)
	}
	if !strings.Contains(stdout, "run-event-re") || !strings.Contains(stdout, "event output") {
		t.Fatalf("logs --event output = %q", stdout)
	}
	if len(run.listRequests) != 2 {
		t.Fatalf("ListRuns calls = %d, want 2 (deduplicated scheduler runs)", len(run.listRequests))
	}
	for _, request := range run.listRequests {
		if request.GetSchedulerRunId() == "" {
			t.Fatalf("ListRuns without scheduler_run_id filter: %+v", request)
		}
		if request.GetProjectId() != "project-aaaaaaaaaaaaaaaa" {
			t.Fatalf("ListRuns project id = %q", request.GetProjectId())
		}
	}

	for _, testCase := range []struct {
		args    []string
		message string
	}{
		{[]string{"logs", "--event", "evt_event_test", "--run", "run-x"}, "--run and --event are mutually exclusive"},
		{[]string{"logs", "--event", "evt_event_test", "--sandbox", "sandbox-x"}, "--sandbox and --event are mutually exclusive"},
		{[]string{"logs", "agent", "--event", "evt_event_test"}, "a positional target and --event are mutually exclusive"},
		{[]string{"logs", "--event", "evt_"}, "full event id"},
	} {
		_, _, err := executeCommand(t, server.URL, false, testCase.args...)
		cliErr, ok := err.(*CLIError)
		if !ok || cliErr.ExitCode() != exitUsage {
			t.Fatalf("%v error = %#v, want usage error", testCase.args, err)
		}
		if !strings.Contains(cliErr.Message, testCase.message) {
			t.Fatalf("%v message = %q, want %q", testCase.args, cliErr.Message, testCase.message)
		}
	}
}

func TestExecuteLogsForEventResolvesAgentReference(t *testing.T) {
	run := &schedulerRunStub{
		runsByScheduler: map[string][]*agentcomposev2.RunSummary{
			"sched-run-agent": {
				{RunId: "run-event-agent", ProjectId: "project-aaaaaaaaaaaaaaaa", AgentName: "agent", SandboxId: "sandbox-1"},
				{RunId: "run-event-other", ProjectId: "project-aaaaaaaaaaaaaaaa", AgentName: "other", SandboxId: "sandbox-2"},
			},
		},
		followRunOutputs: map[string]string{"run-event-agent": "agent output\n", "run-event-other": "other output\n"},
	}
	server := newEventLogsTestServer(t, run, eventTraceResponse(t, "evt_agent_ref", "sched-run-agent"))

	stdout, stderr, err := executeCommand(t, server.URL, false, "logs", "--event", "evt_agent_ref", "--agent", "agent-bbbbbbbbbbbbbbbb")
	if err != nil {
		t.Fatalf("logs --event --agent <managed-id> returned error: %v\nstderr=%s", err, stderr)
	}
	if !strings.Contains(stdout, "agent output") || strings.Contains(stdout, "other output") {
		t.Fatalf("logs --event --agent <managed-id> output = %q", stdout)
	}
}

func TestExecuteLogsForEventJSONOutputAndEmptyNotice(t *testing.T) {
	run := &schedulerRunStub{
		runsByScheduler: map[string][]*agentcomposev2.RunSummary{
			"sched-run-reviewer": {{RunId: "run-event-reviewer", ProjectId: "project-aaaaaaaaaaaaaaaa", AgentName: "agent", SandboxId: "sandbox-1"}},
		},
		followRunOutputs: map[string]string{"run-event-reviewer": "json output\n"},
	}
	server := newEventLogsTestServer(t, run, eventTraceResponse(t, "evt_json_test", "sched-run-reviewer"))

	stdout, _, err := executeCommand(t, server.URL, false, "--json", "logs", "--event", "evt_json_test")
	if err != nil {
		t.Fatalf("logs --event --json returned error: %v", err)
	}
	var decoded struct {
		AgentName string `json:"agent_name"`
		RunID     string `json:"run_id"`
		Offset    uint64 `json:"offset"`
		IsFinal   bool   `json:"is_final"`
		RunStatus string `json:"run_status"`
		Content   string `json:"content"`
	}
	decoder := json.NewDecoder(strings.NewReader(stdout))
	var sawContent, sawFinal bool
	for decoder.More() {
		if err := decoder.Decode(&decoded); err != nil {
			t.Fatalf("decode NDJSON line: %v\n%s", err, stdout)
		}
		if decoded.RunID != "run-event-reviewer" {
			t.Fatalf("NDJSON run_id = %q", decoded.RunID)
		}
		if decoded.Content == "json output\n" {
			sawContent = true
		}
		if decoded.IsFinal {
			sawFinal = true
		}
	}
	if !sawContent || !sawFinal {
		t.Fatalf("NDJSON output missing content or final chunk: %s", stdout)
	}

	empty := &schedulerRunStub{runsByScheduler: map[string][]*agentcomposev2.RunSummary{}}
	emptyServer := newEventLogsTestServer(t, empty, eventTraceResponse(t, "evt_empty_test"))
	jsonOut, errOut, err := executeCommand(t, emptyServer.URL, false, "--json", "logs", "--event", "evt_empty_test")
	if err != nil {
		t.Fatalf("logs --event empty --json returned error: %v\nstderr=%s", err, errOut)
	}
	if strings.TrimSpace(jsonOut) != `{"runs":[]}` {
		t.Fatalf("logs --event empty --json output = %q", jsonOut)
	}

	textOut, textErr, err := executeCommand(t, emptyServer.URL, false, "logs", "--event", "evt_empty_test")
	if err != nil {
		t.Fatalf("logs --event empty returned error: %v", err)
	}
	if strings.TrimSpace(textOut) != "" || !strings.Contains(textErr, "No runs are associated with event evt_empty_test") {
		t.Fatalf("logs --event empty stdout/stderr = %q / %q", textOut, textErr)
	}
}

func TestResolveEventLogSchedulerRunIDsNotFoundAndStatuses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/events/evt_missing/trace":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"event not found"}`))
		case "/api/events/evt_broken/trace":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"failed to trace event"}`))
		case "/api/events/evt_bad/trace":
			_, _ = w.Write([]byte("not json"))
		default:
			t.Fatalf("unexpected trace path %q", r.URL.Path)
		}
	}))
	defer server.Close()
	state := &commandState{options: runtimeOptions{URL: server.URL, Token: "test-token"}}

	_, err := resolveEventLogSchedulerRunIDs(context.Background(), state, "evt_missing")
	cliErr, ok := err.(*CLIError)
	if !ok || cliErr.ExitCode() != exitNotFound || !strings.Contains(cliErr.Message, "evt_missing not found") {
		t.Fatalf("not-found error = %#v", err)
	}

	_, err = resolveEventLogSchedulerRunIDs(context.Background(), state, "evt_broken")
	cliErr, ok = err.(*CLIError)
	if !ok || cliErr.ExitCode() != exitNetwork {
		t.Fatalf("server error = %#v", err)
	}

	_, err = resolveEventLogSchedulerRunIDs(context.Background(), state, "evt_bad")
	cliErr, ok = err.(*CLIError)
	if !ok || cliErr.ExitCode() != exitGeneral || !strings.Contains(cliErr.Message, "parse event") {
		t.Fatalf("malformed json error = %#v", err)
	}
}

func TestResolveEventLogSchedulerRunIDsDeduplicates(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		runs := []map[string]any{
			{"delivery": map[string]any{"run_id": "sched-run-1"}},
			{"delivery": map[string]any{"run_id": "sched-run-2"}},
			{"delivery": map[string]any{"run_id": "sched-run-1"}},
			{"delivery": map[string]any{"run_id": ""}},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"event": map[string]any{"id": "evt_x"}, "runs": runs})
	}))
	defer server.Close()
	state := &commandState{options: runtimeOptions{URL: server.URL, Token: "test-token"}}

	runIDs, err := resolveEventLogSchedulerRunIDs(context.Background(), state, "evt_x")
	if err != nil {
		t.Fatalf("resolveEventLogSchedulerRunIDs returned error: %v", err)
	}
	if gotPath != "/api/events/evt_x/trace" {
		t.Fatalf("trace path = %q", gotPath)
	}
	if len(runIDs) != 2 || runIDs[0] != "sched-run-1" || runIDs[1] != "sched-run-2" {
		t.Fatalf("scheduler run ids = %#v", runIDs)
	}
}
