package commands

import "github.com/spf13/cobra"

func NewSystemCommand() *cobra.Command {
	return domainCommand("system", "系统状态查询", "+status")
}
