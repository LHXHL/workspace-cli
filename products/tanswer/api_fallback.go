package tanswer

import "github.com/spf13/cobra"

func NewAPIFallbackCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "api <METHOD> <PATH>",
		Short: "调用语义命令未覆盖的 Open API",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}
