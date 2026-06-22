package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/chaitin/chaitin-cli/products/insight/client"
	"github.com/spf13/cobra"
)

func NewVulnCmd(getClient func(cmd *cobra.Command) *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vuln",
		Short: "Manage Vulnerabilities",
	}

	ipCmd := &cobra.Command{
		Use:   "ip",
		Short: "List IP vulnerabilities",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient(cmd)

			count, _ := cmd.Flags().GetInt64("count")
			offset, _ := cmd.Flags().GetUint64("offset")

			// Simple payload based on the JSON-RPC V2 specification for Insight
			payload := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      "0",
				"method":  "ScanVulnIpService.SearchScanVulnIpList",
				"params": map[string]interface{}{
					"count":  count,
					"offset": offset,
				},
			}

			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				return err
			}

			resp, err := c.Request("POST", "/pedestal/rpc", bytes.NewReader(payloadBytes))
			if err != nil {
				return err
			}

			format, _ := cmd.Flags().GetString("output")
			renderer := client.NewRenderer(format, cmd.OutOrStdout())
			return renderer.Render(resp)
		},
	}
	ipCmd.Flags().Int64("count", 20, "Number of items to return")
	ipCmd.Flags().Uint64("offset", 0, "Number of items to skip")
	cmd.AddCommand(ipCmd)

	webCmd := &cobra.Command{
		Use:   "web",
		Short: "List Web vulnerabilities",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient(cmd)

			count, _ := cmd.Flags().GetInt64("count")
			offset, _ := cmd.Flags().GetUint64("offset")

			payload := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      "0",
				"method":  "ScanVulnIpService.SearchScanVulnWebList",
				"params": map[string]interface{}{
					"count":  count,
					"offset": offset,
				},
			}

			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				return err
			}

			resp, err := c.Request("POST", "/pedestal/rpc", bytes.NewReader(payloadBytes))
			if err != nil {
				return err
			}

			format, _ := cmd.Flags().GetString("output")
			renderer := client.NewRenderer(format, cmd.OutOrStdout())
			return renderer.Render(resp)
		},
	}
	webCmd.Flags().Int64("count", 20, "Number of items to return")
	webCmd.Flags().Uint64("offset", 0, "Number of items to skip")
	cmd.AddCommand(webCmd)

	return cmd
}
