package commands

import "github.com/spf13/cobra"

func NewFileAlarmCommand() *cobra.Command {
	return domainCommand("file-alarm", "文件告警查询", "+overview", "+malicious", "+webshell", "+sandbox", "+detail")
}
