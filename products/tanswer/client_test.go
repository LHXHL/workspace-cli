package tanswer

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestClientBuildsJSONRPCRequest(t *testing.T) {
	client := NewClient(RuntimeConfig{
		URL:      "https://quanxi.example.com/",
		APIKey:   "token-placeholder",
		Timeout:  5 * time.Second,
		Insecure: true,
	})

	req, err := client.NewJSONRPCRequest(context.Background(), "OpsService.GetBaseInfo", map[string]any{"verbose": true})
	if err != nil {
		t.Fatal(err)
	}
	if req.Method != http.MethodPost {
		t.Fatalf("method = %s", req.Method)
	}
	if req.URL.String() != "https://quanxi.example.com/rpc" {
		t.Fatalf("url = %s", req.URL.String())
	}
	if got := req.Header.Get("API-Token"); got != "token-placeholder" {
		t.Fatalf("API-Token header = %q", got)
	}
	if got := req.Header.Get("Content-Type"); !strings.Contains(got, "application/json") {
		t.Fatalf("Content-Type = %q", got)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["jsonrpc"] != "2.0" || payload["method"] != "OpsService.GetBaseInfo" {
		t.Fatalf("bad JSON-RPC payload: %#v", payload)
	}
	params := payload["params"].(map[string]any)
	if params["verbose"] != true {
		t.Fatalf("params = %#v", params)
	}
}
