package tanswer

import (
	"bytes"
	"strings"
	"testing"
)

func TestAuthStatusHelpMentionsPurpose(t *testing.T) {
	var out bytes.Buffer
	cmd := NewRootCommand(RootOptions{Out: &out, ErrOut: &out})
	cmd.SetArgs([]string{"tanswer", "auth", "status", "--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	text := out.String()
	for _, want := range []string{"查看当前连接状态", "TANSWER_URL", "OpenAPI Token"} {
		if !strings.Contains(text, want) {
			t.Fatalf("help missing %q:\n%s", want, text)
		}
	}
}

func TestAuthCheckHelpMentionsTokenValidation(t *testing.T) {
	var out bytes.Buffer
	cmd := NewRootCommand(RootOptions{Out: &out, ErrOut: &out})
	cmd.SetArgs([]string{"tanswer", "auth", "check", "--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if !strings.Contains(out.String(), "校验 OpenAPI Token") {
		t.Fatalf("help missing token validation purpose:\n%s", out.String())
	}
}
