package tanswer

import (
	"time"

	tanswercommands "github.com/chaitin/chaitin-cli/products/tanswer/commands"
	"github.com/spf13/cobra"
)

// NewCommand creates the T-Answer product command. The AI-readable semantic
// commands are the only first-class entry points; old firewall/rules aliases
// are intentionally not registered.
func NewCommand() *cobra.Command {
	cfg := DefaultRuntimeConfig()

	cmd := &cobra.Command{
		Use:   "tanswer",
		Short: "全悉 AI 可读 CLI",
		Long:  "全悉面向人类操作者和 AI Agent 的语义命令入口，写操作必须 preview/confirm。",
	}

	cmd.PersistentFlags().StringVar(&cfg.URL, "url", cfg.URL, "T-Answer console URL")
	cmd.PersistentFlags().StringVar(&cfg.APIKey, "api-key", cfg.APIKey, "Open API key for T-Answer")
	cmd.PersistentFlags().DurationVar(&cfg.Timeout, "timeout", cfg.Timeout, "Request timeout")
	cmd.PersistentFlags().BoolVar(&cfg.Insecure, "insecure", cfg.Insecure, "Skip TLS certificate verification")
	cmd.PersistentFlags().StringVar(&cfg.Output, "output", cfg.Output, "Output format, currently json")

	cmd.AddCommand(
		tanswercommands.NewAuthCommand(),
		NewManifestCommand(),
		tanswercommands.NewSystemCommand(),
		tanswercommands.NewAlarmCommand(),
		tanswercommands.NewFileAlarmCommand(),
		tanswercommands.NewAssetCommand(),
		tanswercommands.NewMetadataCommand(),
		tanswercommands.NewPolicyCommand(),
		tanswercommands.NewResponseCommand(),
		NewAPIFallbackCommand(),
	)

	return cmd
}

func durationOrDefault(v time.Duration, fallback time.Duration) time.Duration {
	if v <= 0 {
		return fallback
	}
	return v
}
