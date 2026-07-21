package commands

import "github.com/spf13/cobra"

func NewAssetCommand() *cobra.Command {
	return domainCommand("asset", "资产配置读写",
		"list", "detail", "group-tree", "download-template", "export",
		"create", "update", "delete", "batch-maintain", "batch-tag",
		"group-create", "group-rename", "group-delete", "tree-move", "import",
	)
}
