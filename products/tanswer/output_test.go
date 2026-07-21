package tanswer

import (
	"encoding/json"
	"testing"
)

func TestSuccessEnvelopeContract(t *testing.T) {
	out := NewSuccessEnvelope("查看威胁告警概览", "chaitin-cli tanswer alarm overview", map[string]any{"time": "today"}, map[string]any{"total": 1}, nil)
	b, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if got["success"] != true {
		t.Fatalf("success = %v", got["success"])
	}
	if got["task"] == "" || got["command"] == "" || got["data"] == nil {
		t.Fatalf("missing stable envelope fields: %#v", got)
	}
}
