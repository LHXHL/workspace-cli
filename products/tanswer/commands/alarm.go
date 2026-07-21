package commands

import "github.com/spf13/cobra"

func NewAlarmCommand() *cobra.Command {
	return domainCommand("alarm", "威胁告警查询",
		"+overview", "+timeline", "+list", "+high-priority", "+detail",
		"+by-attacker", "+by-victim", "+by-threat", "+important-assets",
		"+attacker-rank", "+victim-rank", "+phase-distribution", "+related",
	)
}
