package tanswer

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestRootCommandRegistersAgentReadableCommands(t *testing.T) {
	cmd := NewCommand()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})

	expected := []string{"auth", "manifest", "system", "alarm", "file-alarm", "asset", "metadata", "policy", "response", "api"}
	for _, name := range expected {
		found := false
		for _, child := range cmd.Commands() {
			if child.Name() == name {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected child command %q", name)
		}
	}
}

func TestTAnswerHelpGuidesUsersWithoutProductDocuments(t *testing.T) {
	var out bytes.Buffer
	cmd := NewCommand()
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	for _, want := range []string{
		"TANSWER_URL", "TANSWER_API_KEY", "auth check", "manifest",
		"--help", "preview", "confirm", "semantic commands",
	} {
		if !strings.Contains(strings.ToLower(out.String()), strings.ToLower(want)) {
			t.Fatalf("help missing %q:\n%s", want, out.String())
		}
	}
}

func TestManifestCommandsAreDiscoverableFromHelp(t *testing.T) {
	for _, entry := range BuildCommandManifest().Commands {
		var out bytes.Buffer
		root := NewRootCommand(RootOptions{Out: &out, ErrOut: &out})
		args := strings.Fields(strings.TrimPrefix(entry.FullCommand, "chaitin-cli "))
		command, _, err := root.Find(args)
		if err != nil || command == nil {
			t.Fatalf("manifest command %s is not registered: %v", entry.FullCommand, err)
		}
		if strings.TrimSpace(command.Long) == "" {
			t.Fatalf("manifest command %s has no detailed help", entry.FullCommand)
		}
		args = append(args, "--help")
		root.SetArgs(args)

		if err := root.Execute(); err != nil {
			t.Fatalf("%s help returned error: %v", entry.FullCommand, err)
		}
		help := out.String()
		if !strings.Contains(help, command.Long) {
			t.Fatalf("%s did not render detailed help:\n%s", entry.FullCommand, help)
		}
		if entry.RequiresConfirmation && (!strings.Contains(strings.ToLower(help), "preview") || !strings.Contains(strings.ToLower(help), "confirm")) {
			t.Fatalf("protected command %s help must mention preview and confirm:\n%s", entry.FullCommand, help)
		}
	}
}

func TestLegacyGuidanceCoverageMapsEveryTAnswerDomainToRuntimeGuidance(t *testing.T) {
	raw, err := os.ReadFile("testdata/migration-baseline/guidance-coverage.json")
	if err != nil {
		t.Fatalf("read guidance coverage: %v", err)
	}
	var coverage struct {
		Rules []struct {
			ID              string   `json:"id"`
			Source          []string `json:"source"`
			RuntimeSurfaces []string `json:"runtime_surfaces"`
			AutomatedTest   string   `json:"automated_test"`
			Scenario        string   `json:"scenario"`
		} `json:"rules"`
	}
	if err := json.Unmarshal(raw, &coverage); err != nil {
		t.Fatalf("unmarshal guidance coverage: %v", err)
	}
	want := map[string]bool{
		"shared": false, "auth": false, "alarm": false, "file-alarm": false,
		"asset": false, "system": false, "metadata": false, "policy": false,
		"response": false, "api-fallback": false,
	}
	for _, rule := range coverage.Rules {
		if _, ok := want[rule.ID]; ok && len(rule.Source) > 0 && len(rule.RuntimeSurfaces) > 0 && rule.AutomatedTest != "" && rule.Scenario != "" {
			want[rule.ID] = true
		}
	}
	for id, covered := range want {
		if !covered {
			t.Fatalf("legacy guidance coverage missing complete rule for %q", id)
		}
	}
}
