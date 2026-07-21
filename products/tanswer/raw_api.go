package tanswer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func newRawAPICommand(opts *RootOptions) *cobra.Command {
	var queryText string
	var bodyText string
	cmd := &cobra.Command{
		Use:   "api <METHOD> <PATH>",
		Short: "调用用户已知且已授权的 Open API",
		Long:  "Open API 兜底调用。用于访问语义快捷层未覆盖、但用户已知且已授权的全悉 Open API。该命令不提供专属参数说明或专属摘要，只结构化返回状态码和原始响应。",
		Example: "  chaitin-cli tanswer api POST /rpc --body '{\"jsonrpc\":\"2.0\",\"method\":\"OpsService.GetBaseInfo\",\"params\":{},\"id\":\"1\"}'\n" +
			"  chaitin-cli tanswer api GET /api/openapi/rpc/openapi.json",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := LoadConfig(ConfigOptions{
				Address:            opts.Address,
				Token:              opts.Token,
				Timeout:            opts.Timeout,
				Format:             opts.Format,
				InsecureSkipVerify: opts.InsecureSkipVerify,
			})
			if err != nil {
				return err
			}
			query, err := parseQueryJSON(queryText)
			if err != nil {
				return err
			}
			body, err := parseBodyJSON(bodyText)
			if err != nil {
				return err
			}
			client := NewClient(cfg)
			resp, err := client.DoJSON(cmd.Context(), strings.ToUpper(args[0]), args[1], query, body)
			if err != nil {
				return err
			}
			raw, err := RenderJSON(SuccessEnvelope{
				Success: resp.StatusCode >= 200 && resp.StatusCode < 300,
				Task:    "Open API 兜底调用",
				Command: "chaitin-cli tanswer api",
				Data: map[string]any{
					"status_code": resp.StatusCode,
					"raw":         json.RawMessage(resp.Body),
				},
			})
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(raw))
			return nil
		},
	}
	cmd.Flags().StringVar(&queryText, "query", "", "query parameters as JSON object")
	cmd.Flags().StringVar(&bodyText, "body", "", "request body as JSON object or @file path")
	return cmd
}

func parseQueryJSON(text string) (map[string]string, error) {
	if strings.TrimSpace(text) == "" {
		return nil, nil
	}
	raw := map[string]any{}
	if err := json.Unmarshal([]byte(text), &raw); err != nil {
		return nil, err
	}
	out := map[string]string{}
	for key, value := range raw {
		out[key] = fmt.Sprint(value)
	}
	return out, nil
}

func parseBodyJSON(text string) (any, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, nil
	}
	if strings.HasPrefix(text, "@") {
		raw, err := os.ReadFile(strings.TrimPrefix(text, "@"))
		if err != nil {
			return nil, err
		}
		text = string(raw)
	}
	var body any
	if err := json.Unmarshal([]byte(text), &body); err != nil {
		return nil, err
	}
	return body, nil
}
