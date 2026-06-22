package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/chaitin/chaitin-cli/products/insight/client"
	"github.com/spf13/cobra"
)

func NewAssetCmd(getClient func(cmd *cobra.Command) *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asset",
		Short: "Manage Assets",
	}

	ipCmd := &cobra.Command{
		Use:   "ip",
		Short: "List IP assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient(cmd)

			count, _ := cmd.Flags().GetInt64("count")
			offset, _ := cmd.Flags().GetUint64("offset")

			payload := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      "0",
				"method":  "AssetMgrService.IpAssetList",
				"params": map[string]interface{}{
					"filter": map[string]interface{}{},
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
		Short: "List Web assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient(cmd)

			count, _ := cmd.Flags().GetInt64("count")
			offset, _ := cmd.Flags().GetUint64("offset")

			payload := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      "0",
				"method":  "AssetMgrService.WebsiteAssetList",
				"params": map[string]interface{}{
					"filter": map[string]interface{}{},
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

	softwareCmd := &cobra.Command{
		Use:   "software",
		Short: "List Software assets overview",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient(cmd)

			count, _ := cmd.Flags().GetInt64("count")
			offset, _ := cmd.Flags().GetUint64("offset")

			payload := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      "0",
				"method":  "AssetMgrService.SoftwareAssetOverviewList",
				"params": map[string]interface{}{
					"filter": map[string]interface{}{},
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
	softwareCmd.Flags().Int64("count", 20, "Number of items to return")
	softwareCmd.Flags().Uint64("offset", 0, "Number of items to skip")
	cmd.AddCommand(softwareCmd)

	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "List Asset Tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient(cmd)

			count, _ := cmd.Flags().GetInt64("count")
			offset, _ := cmd.Flags().GetUint64("offset")

			payload := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      "0",
				"method":  "AssetMgrService.AssetTagList",
				"params": map[string]interface{}{
					"filter": map[string]interface{}{},
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
	tagCmd.Flags().Int64("count", 20, "Number of items to return")
	tagCmd.Flags().Uint64("offset", 0, "Number of items to skip")
	cmd.AddCommand(tagCmd)

	return cmd
}
