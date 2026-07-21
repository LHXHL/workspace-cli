package tanswer

import (
	"bytes"
	"strings"
	"testing"
)

func TestRawAPIHelpExplainsFallback(t *testing.T) {
	var out bytes.Buffer
	cmd := NewRootCommand(RootOptions{Out: &out, ErrOut: &out})
	cmd.SetArgs([]string{"tanswer", "api", "--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	text := out.String()
	for _, want := range []string{"Open API 兜底调用", "METHOD", "PATH", "--query", "--body"} {
		if !strings.Contains(text, want) {
			t.Fatalf("help missing %q:\n%s", want, text)
		}
	}
}
