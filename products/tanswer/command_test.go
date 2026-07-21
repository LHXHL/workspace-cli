package tanswer

import (
	"bytes"
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
