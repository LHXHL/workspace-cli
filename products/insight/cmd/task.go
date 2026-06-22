package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/chaitin/chaitin-cli/products/insight/client"
	"github.com/spf13/cobra"
)

func NewTaskCmd(getClient func(cmd *cobra.Command) *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Manage Exposure tasks",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient(cmd)

			id, _ := cmd.Flags().GetString("id")

			query := ""
			if id != "" {
				query = "?id=" + id
			}

			resp, err := c.Request("GET", "/exposure/api/task"+query, nil)
			if err != nil {
				return err
			}

			format, _ := cmd.Flags().GetString("output")
			renderer := client.NewRenderer(format, cmd.OutOrStdout())
			return renderer.Render(resp)
		},
	}

	listCmd.Flags().String("id", "", "task id to filter")
	cmd.AddCommand(listCmd)
	
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a task re-execution",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient(cmd)

			id, _ := cmd.Flags().GetString("id")

			payload := map[string]interface{}{
				"id": id,
			}

			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				return err
			}

			resp, err := c.Request("POST", "/exposure/api/task/reexecute", bytes.NewReader(payloadBytes))
			if err != nil {
				return err
			}

			format, _ := cmd.Flags().GetString("output")
			renderer := client.NewRenderer(format, cmd.OutOrStdout())
			return renderer.Render(resp)
		},
	}
	startCmd.Flags().String("id", "", "task id to start")
	startCmd.MarkFlagRequired("id")
	cmd.AddCommand(startCmd)

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a task execution",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient(cmd)

			id, _ := cmd.Flags().GetString("id")

			payload := map[string]interface{}{
				"id": id,
			}

			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				return err
			}

			resp, err := c.Request("POST", "/exposure/api/task/stop", bytes.NewReader(payloadBytes))
			if err != nil {
				return err
			}

			format, _ := cmd.Flags().GetString("output")
			renderer := client.NewRenderer(format, cmd.OutOrStdout())
			return renderer.Render(resp)
		},
	}
	stopCmd.Flags().String("id", "", "task execution id to stop")
	stopCmd.MarkFlagRequired("id")
	cmd.AddCommand(stopCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Check task execution status",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient(cmd)

			id, _ := cmd.Flags().GetString("exec-id")

			query := "?id=" + id

			resp, err := c.Request("GET", "/exposure/api/task/execution"+query, nil)
			if err != nil {
				return err
			}

			format, _ := cmd.Flags().GetString("output")
			renderer := client.NewRenderer(format, cmd.OutOrStdout())
			return renderer.Render(resp)
		},
	}
	statusCmd.Flags().String("exec-id", "", "execution id to check status")
	statusCmd.MarkFlagRequired("exec-id")
	cmd.AddCommand(statusCmd)

	return cmd
}
