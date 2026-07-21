package commands

import "github.com/spf13/cobra"

func NewMetadataCommand() *cobra.Command {
	return domainCommand("metadata", "流量元数据查询和配置", "+protocol", "+search", "+detail", "+near-alarm", "+config", "+config-update")
}
