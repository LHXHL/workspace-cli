package commands

import "github.com/spf13/cobra"

func NewResponseCommand() *cobra.Command {
	return domainCommand("response", "响应处置读写",
		"+block-policies", "+block-policy-create", "+block-policy-update",
		"+block-policy-enable", "+block-policy-disable", "+block-policy-delete",
		"+block-records", "+whitelist", "+whitelist-create", "+whitelist-update",
		"+whitelist-enable", "+whitelist-disable", "+whitelist-delete",
		"+block-policy-from-alarm", "+whitelist-from-alarm", "+devices",
		"+device-records", "+auto-policies", "+auto-list",
	)
}
