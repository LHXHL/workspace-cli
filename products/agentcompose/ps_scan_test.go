package agentcompose

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	agentcomposev2 "github.com/chaitin/agent-compose/proto/agentcompose/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// hourlySchedulerRuns returns count Scheduler Runs one hour apart, newest first,
// matching the daemon's started_at DESC order.
func hourlySchedulerRuns(newest time.Time, count int) []*agentcomposev2.SchedulerRun {
	runs := make([]*agentcomposev2.SchedulerRun, count)
	for i := range runs {
		runs[i] = &agentcomposev2.SchedulerRun{
			RunId:      fmt.Sprintf("scheduler-run-%04d", i),
			StartedAt:  timestamppb.New(newest.Add(-time.Duration(i) * time.Hour)),
			SandboxIds: []string{fmt.Sprintf("other-sandbox-%04d", i)},
		}
	}
	return runs
}

func psSandbox(id string, createdAt *timestamppb.Timestamp) *agentcomposev2.Sandbox {
	return &agentcomposev2.Sandbox{SandboxId: id, ProjectId: "project-aaaaaaaaaaaaaaaa", Status: agentcomposev2.SandboxStatus_SANDBOX_STATUS_RUNNING, CreatedAt: createdAt}
}

func psRunIDs(t *testing.T, out string) map[string]string {
	t.Helper()
	var payload struct {
		Sandboxes []sandboxDTO `json:"sandboxes"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("output = %s: %v", out, err)
	}
	runIDs := make(map[string]string, len(payload.Sandboxes))
	for _, sandbox := range payload.Sandboxes {
		runIDs[sandbox.SandboxID] = sandbox.RunID
	}
	return runIDs
}

func TestPSStopsSchedulerRunScanOnceEverySandboxIsMatched(t *testing.T) {
	newest := time.Date(2026, 10, 15, 12, 0, 0, 0, time.UTC)
	runs := hourlySchedulerRuns(newest, 1000)
	runs[150].SandboxIds = []string{"target"}
	project := &projectStub{project: fixtureProject(), schedulerRuns: runs}
	// No creation time means no cutoff, so only the matched-everything stop can
	// end the scan on the second page.
	sandbox := &sandboxStub{sandboxes: []*agentcomposev2.Sandbox{psSandbox("target", nil)}}
	server, _ := newTestServer(t, project, &runStub{}, sandbox)
	out, _, err := executeCommand(t, server.URL, false, "--json", "ps")
	if err != nil {
		t.Fatal(err)
	}
	if got := psRunIDs(t, out)["target"]; got != "scheduler-run-0150" {
		t.Fatalf("target run_id = %q", got)
	}
	if len(project.schedulerRunRequests) != 2 {
		t.Fatalf("ListSchedulerRuns calls = %d, want 2", len(project.schedulerRunRequests))
	}
}

func TestPSStopsSchedulerRunScanAtCutoff(t *testing.T) {
	newest := time.Date(2026, 10, 15, 12, 0, 0, 0, time.UTC)
	project := &projectStub{project: fixtureProject(), schedulerRuns: hourlySchedulerRuns(newest, 1000)}
	// Created 10h before the newest run, so the cutoff is 34h back: runs 0..34
	// are in range and the scan ends with the page that crosses it.
	sandbox := &sandboxStub{sandboxes: []*agentcomposev2.Sandbox{psSandbox("orphan", timestamppb.New(newest.Add(-10*time.Hour)))}}
	server, _ := newTestServer(t, project, &runStub{}, sandbox)
	out, _, err := executeCommand(t, server.URL, false, "--json", "ps")
	if err != nil {
		t.Fatal(err)
	}
	if got := psRunIDs(t, out)["orphan"]; got != "" {
		t.Fatalf("orphan run_id = %q, want none", got)
	}
	if len(project.schedulerRunRequests) != 1 {
		t.Fatalf("ListSchedulerRuns calls = %d, want 1 of 10 pages", len(project.schedulerRunRequests))
	}
}

func TestPSFindsSchedulerRunStartedBeforeSandboxWithinMargin(t *testing.T) {
	newest := time.Date(2026, 10, 15, 12, 0, 0, 0, time.UTC)
	runs := hourlySchedulerRuns(newest, 1000)
	// The run started 20h before it created the Sandbox, inside the 24h margin,
	// and sits on the second page.
	runs[120].SandboxIds = []string{"late"}
	project := &projectStub{project: fixtureProject(), schedulerRuns: runs}
	sandbox := &sandboxStub{sandboxes: []*agentcomposev2.Sandbox{psSandbox("late", timestamppb.New(newest.Add(-100*time.Hour)))}}
	server, _ := newTestServer(t, project, &runStub{}, sandbox)
	out, _, err := executeCommand(t, server.URL, false, "--json", "ps")
	if err != nil {
		t.Fatal(err)
	}
	if got := psRunIDs(t, out)["late"]; got != "scheduler-run-0120" {
		t.Fatalf("late run_id = %q", got)
	}
}

func TestPSScansAllSchedulerRunsWhenSandboxCreationTimeIsUnknown(t *testing.T) {
	newest := time.Date(2026, 10, 15, 12, 0, 0, 0, time.UTC)
	project := &projectStub{project: fixtureProject(), schedulerRuns: hourlySchedulerRuns(newest, 1000)}
	sandbox := &sandboxStub{sandboxes: []*agentcomposev2.Sandbox{
		psSandbox("recent", timestamppb.New(newest)),
		psSandbox("unknown", nil),
	}}
	server, _ := newTestServer(t, project, &runStub{}, sandbox)
	if _, _, err := executeCommand(t, server.URL, false, "--json", "ps"); err != nil {
		t.Fatal(err)
	}
	if len(project.schedulerRunRequests) != 10 {
		t.Fatalf("ListSchedulerRuns calls = %d, want all 10 pages", len(project.schedulerRunRequests))
	}
}

func TestPSSkipsSchedulerRunsWhenRunsCoverEverySandbox(t *testing.T) {
	newest := time.Date(2026, 10, 15, 12, 0, 0, 0, time.UTC)
	project := &projectStub{project: fixtureProject(), schedulerRuns: hourlySchedulerRuns(newest, 1000)}
	run := &runStub{runs: []*agentcomposev2.RunSummary{{RunId: "run-1", ProjectId: "project-aaaaaaaaaaaaaaaa", SandboxId: "covered"}}}
	sandbox := &sandboxStub{sandboxes: []*agentcomposev2.Sandbox{psSandbox("covered", timestamppb.New(newest))}}
	server, _ := newTestServer(t, project, run, sandbox)
	out, _, err := executeCommand(t, server.URL, false, "--json", "ps")
	if err != nil {
		t.Fatal(err)
	}
	if got := psRunIDs(t, out)["covered"]; got != "run-1" {
		t.Fatalf("covered run_id = %q", got)
	}
	if len(project.schedulerRunRequests) != 0 {
		t.Fatalf("ListSchedulerRuns calls = %d, want 0", len(project.schedulerRunRequests))
	}
}

func TestPSLooksUpSchedulerRunsByAgentFirst(t *testing.T) {
	newest := time.Date(2026, 10, 15, 12, 0, 0, 0, time.UTC)
	// A noisy agent owns the 1000 newest runs; the quiet agent's run that
	// created the Sandbox is older than all of them.
	runs := hourlySchedulerRuns(newest, 1000)
	for _, run := range runs {
		run.AgentName = "noisy"
	}
	created := newest.Add(-1000 * time.Hour)
	runs = append(runs, &agentcomposev2.SchedulerRun{RunId: "quiet-run", AgentName: "quiet", StartedAt: timestamppb.New(created), SandboxIds: []string{"quiet-sandbox"}})
	project := &projectStub{project: fixtureProject(), schedulerRuns: runs}
	target := psSandbox("quiet-sandbox", timestamppb.New(created))
	target.AgentName = "quiet"
	server, _ := newTestServer(t, project, &runStub{}, &sandboxStub{sandboxes: []*agentcomposev2.Sandbox{target}})
	out, _, err := executeCommand(t, server.URL, false, "--json", "ps")
	if err != nil {
		t.Fatal(err)
	}
	if got := psRunIDs(t, out)["quiet-sandbox"]; got != "quiet-run" {
		t.Fatalf("quiet-sandbox run_id = %q", got)
	}
	if len(project.schedulerRunRequests) != 1 || project.schedulerRunRequests[0].GetAgentName() != "quiet" {
		t.Fatalf("requests = %+v, want one filtered by agent", project.schedulerRunRequests)
	}
}

func TestPSFallsBackToAllSchedulerRunsWhenAgentDiffers(t *testing.T) {
	newest := time.Date(2026, 10, 15, 12, 0, 0, 0, time.UTC)
	runs := hourlySchedulerRuns(newest, 300)
	for _, run := range runs {
		run.AgentName = "other"
	}
	runs[5].SandboxIds = []string{"cross"}
	project := &projectStub{project: fixtureProject(), schedulerRuns: runs}
	sandbox := psSandbox("cross", timestamppb.New(newest.Add(-5*time.Hour)))
	sandbox.AgentName = "mismatched"
	server, _ := newTestServer(t, project, &runStub{}, &sandboxStub{sandboxes: []*agentcomposev2.Sandbox{sandbox}})
	out, _, err := executeCommand(t, server.URL, false, "--json", "ps")
	if err != nil {
		t.Fatal(err)
	}
	if got := psRunIDs(t, out)["cross"]; got != "scheduler-run-0005" {
		t.Fatalf("cross run_id = %q", got)
	}
	if len(project.schedulerRunRequests) != 2 || project.schedulerRunRequests[0].GetAgentName() != "mismatched" || project.schedulerRunRequests[1].GetAgentName() != "" {
		t.Fatalf("requests = %+v, want an agent scan then an unfiltered fallback", project.schedulerRunRequests)
	}
}
