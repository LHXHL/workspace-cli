package tanswer

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestManifestCommandOutputsAIReadableCommandMetadata(t *testing.T) {
	var out bytes.Buffer
	cmd := NewRootCommand(RootOptions{Out: &out, ErrOut: &out})
	cmd.SetArgs([]string{"tanswer", "manifest"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	var env SuccessEnvelope
	if err := json.Unmarshal(out.Bytes(), &env); err != nil {
		t.Fatalf("manifest output is not envelope json: %v\n%s", err, out.String())
	}
	if env.Task != "查看 CLI 命令清单" {
		t.Fatalf("Task = %q", env.Task)
	}

	raw, err := json.Marshal(env.Data)
	if err != nil {
		t.Fatalf("marshal Data: %v", err)
	}
	var manifest CommandManifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatalf("manifest data shape mismatch: %v\n%s", err, string(raw))
	}

	if manifest.Product != "tanswer" {
		t.Fatalf("Product = %q", manifest.Product)
	}
	requireManifestCommand(t, manifest, "tanswer alarm overview", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer alarm timeline", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer alarm list", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer alarm high-priority", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer alarm detail", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer alarm by-attacker", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer alarm by-victim", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer alarm by-threat", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer alarm important-assets", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer alarm attacker-rank", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer alarm victim-rank", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer alarm phase-distribution", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer alarm related", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer file-alarm overview", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer file-alarm malicious", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer file-alarm webshell", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer file-alarm sandbox", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer file-alarm detail", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer system status", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer asset list", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer asset detail", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer asset group-tree", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer asset download-template", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer asset export", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer asset create", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer asset update", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer asset delete", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer asset batch-maintain", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer asset batch-tag", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer asset group-create", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer asset group-rename", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer asset group-delete", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer asset tree-move", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer asset import", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer metadata protocol", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer metadata search", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer metadata detail", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer metadata near-alarm", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer metadata config", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer metadata config-update", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer policy detection-whitelist", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer policy detection-whitelist-create", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer policy detection-whitelist-update", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer policy detection-whitelist-enable", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer policy detection-whitelist-disable", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer policy detection-whitelist-delete", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer policy detection-whitelist-from-alarm", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer policy detection-whitelist-export", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer policy detection-whitelist-import", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer policy custom-intelligence", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer policy custom-intelligence-create", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer policy custom-intelligence-update", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer policy custom-intelligence-enable", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer policy custom-intelligence-disable", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer policy custom-intelligence-delete", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer policy custom-intelligence-export", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer policy custom-intelligence-import", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer response block-policies", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer response block-policy-create", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer response block-policy-update", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer response block-policy-enable", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer response block-policy-disable", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer response block-policy-delete", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer response block-records", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer response whitelist", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer response whitelist-create", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer response whitelist-update", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer response whitelist-enable", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer response whitelist-disable", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer response whitelist-delete", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer response block-policy-from-alarm", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer response whitelist-from-alarm", "semantic_shortcut", "write_high")
	requireManifestCommand(t, manifest, "tanswer response devices", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer response device-records", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer response auto-policies", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer response auto-list", "semantic_shortcut", "read")
	requireManifestCommand(t, manifest, "tanswer api", "openapi_fallback", "read_write")
	requireManifestCommand(t, manifest, "tanswer auth check", "foundation", "read")
}

func TestManifestHelpExplainsPurpose(t *testing.T) {
	var out bytes.Buffer
	cmd := NewRootCommand(RootOptions{Out: &out, ErrOut: &out})
	cmd.SetArgs([]string{"tanswer", "manifest", "--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	text := out.String()
	for _, want := range []string{"查看 CLI 命令清单", "AI Agent", "risk", "output"} {
		if !strings.Contains(text, want) {
			t.Fatalf("help missing %q:\n%s", want, text)
		}
	}
}

func requireManifestCommand(t *testing.T, manifest CommandManifest, name string, layer string, risk string) {
	t.Helper()
	for _, cmd := range manifest.Commands {
		if cmd.Name == name {
			if cmd.Layer != layer {
				t.Fatalf("%s layer = %q", name, cmd.Layer)
			}
			if cmd.RiskLevel != risk {
				t.Fatalf("%s risk = %q", name, cmd.RiskLevel)
			}
			if risk == "write_high" && !cmd.RequiresConfirmation {
				t.Fatalf("%s requires_confirmation = false", name)
			}
			if cmd.OutputType == "" {
				t.Fatalf("%s output type is empty", name)
			}
			if len(cmd.Examples) == 0 {
				t.Fatalf("%s examples are empty", name)
			}
			return
		}
	}
	t.Fatalf("manifest missing command %q", name)
}
