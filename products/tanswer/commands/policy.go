package commands

import "github.com/spf13/cobra"

func NewPolicyCommand() *cobra.Command {
	return domainCommand("policy", "安全策略读写",
		"detection-whitelist", "detection-whitelist-create", "detection-whitelist-update",
		"detection-whitelist-enable", "detection-whitelist-disable", "detection-whitelist-delete",
		"detection-whitelist-from-alarm", "detection-whitelist-export", "detection-whitelist-import",
		"custom-intelligence", "custom-intelligence-create", "custom-intelligence-update",
		"custom-intelligence-enable", "custom-intelligence-disable", "custom-intelligence-delete",
		"custom-intelligence-export", "custom-intelligence-import",
	)
}
