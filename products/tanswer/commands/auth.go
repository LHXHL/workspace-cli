package commands

import "github.com/spf13/cobra"

func NewAuthCommand() *cobra.Command {
	return domainCommand("auth", "认证和连接检查", "status", "check")
}
