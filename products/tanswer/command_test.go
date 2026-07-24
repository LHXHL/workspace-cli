package tanswer

import (
	"bytes"
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
