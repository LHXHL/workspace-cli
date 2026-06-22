package cmd

import (
	"github.com/chaitin/chaitin-cli/products/insight/client"
	"github.com/spf13/cobra"
)

func NewSnapshotCmd(getClient func(cmd *cobra.Command) *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Asset snapshot management",
	}

	assetCmd := &cobra.Command{
		Use:   "asset",
		Short: "Get managed asset snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient(cmd)

			resp, err := c.Request("GET", "/exposure/api/snapshot/asset", nil)
			if err != nil {
				return err
			}

			format, _ := cmd.Flags().GetString("output")
			renderer := client.NewRenderer(format, cmd.OutOrStdout())
			return renderer.Render(resp)
		},
	}
	cmd.AddCommand(assetCmd)

	return cmd
}
