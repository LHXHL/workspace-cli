package tanswer

import "github.com/spf13/cobra"

func NewManifestCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "manifest",
		Short: "输出全悉 AI 可读命令清单",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}
