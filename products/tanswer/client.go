package tanswer

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type JSONRPCRequest struct {
	ID      string         `json:"id"`
	JSONRPC string         `json:"jsonrpc"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params"`
}

func NewClient(cfg RuntimeConfig) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: cfg.Insecure}

	return &Client{
		baseURL: strings.TrimRight(cfg.URL, "/"),
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout:   durationOrDefault(cfg.Timeout, defaultTimeout),
			Transport: transport,
		},
	}
}

func (c *Client) NewRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	}
	if c.apiKey != "" {
		req.Header.Set("API-Token", c.apiKey)
	}
	return req, nil
}

func (c *Client) NewJSONRPCRequest(ctx context.Context, method string, params map[string]any) (*http.Request, error) {
	if params == nil {
		params = map[string]any{}
	}
	return c.NewRequest(ctx, http.MethodPost, "/rpc", JSONRPCRequest{
		ID:      uuid.NewString(),
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	})
}
