package tanswer

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

type RootOptions struct {
	Address            string
	Token              string
	Timeout            string
	Format             string
	InsecureSkipVerify bool
	Out                io.Writer
	ErrOut             io.Writer
}

// NewCommand creates the T-Answer product command for chaitin-cli.
func NewCommand() *cobra.Command {
	opts := RootOptions{}
	return newTAnswerCommand(&opts, true)
}

// NewRootCommand is kept for product-level tests that exercise the full
// chaitin-cli tanswer command path.
func NewRootCommand(opts RootOptions) *cobra.Command {
	if opts.Out == nil {
		opts.Out = io.Discard
	}
	if opts.ErrOut == nil {
		opts.ErrOut = io.Discard
	}

	root := &cobra.Command{
		Use:           "chaitin-cli",
		Short:         "Chaitin command line tools",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.SetOut(opts.Out)
	root.SetErr(opts.ErrOut)

	root.AddCommand(newTAnswerCommand(&opts, true))
	return root
}

func newTAnswerCommand(opts *RootOptions, registerFlags bool) *cobra.Command {
	if opts == nil {
		opts = &RootOptions{}
	}

	cmd := &cobra.Command{
		Use:   "tanswer",
		Short: "全悉 AI 可读 CLI",
		Long:  "全悉面向人类操作者和 AI Agent 的语义命令入口。优先使用语义命令；仅在目标能力已开放且已授权时使用 api fallback。",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			syncOptionsFromFlags(cmd, opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	if registerFlags {
		cmd.PersistentFlags().StringVar(&opts.Address, "url", "", "T-Answer console URL, or TANSWER_URL")
		cmd.PersistentFlags().StringVar(&opts.Token, "api-key", "", "OpenAPI Token, or TANSWER_API_KEY")
		cmd.PersistentFlags().StringVar(&opts.Timeout, "timeout", "", "request timeout, default 30s")
		cmd.PersistentFlags().StringVar(&opts.Format, "output", "", "output format, default json")
		cmd.PersistentFlags().BoolVar(&opts.InsecureSkipVerify, "insecure", false, "skip TLS certificate verification, or TANSWER_INSECURE")
	}

	cmd.AddCommand(newAuthCommand(opts))
	cmd.AddCommand(newRawAPICommand(opts))
	cmd.AddCommand(newAlarmCommand(opts))
	cmd.AddCommand(newFileAlarmCommand(opts))
	cmd.AddCommand(newAssetCommand(opts))
	cmd.AddCommand(newSystemCommand(opts))
	cmd.AddCommand(newMetadataCommand(opts))
	cmd.AddCommand(newPolicyCommand(opts))
	cmd.AddCommand(newResponseCommand(opts))
	cmd.AddCommand(newManifestCommand(opts))
	return cmd
}

func syncOptionsFromFlags(cmd *cobra.Command, opts *RootOptions) {
	if cmd == nil || opts == nil {
		return
	}
	if value := inheritedFlagValue(cmd, "url"); value != "" {
		opts.Address = value
	}
	if value := inheritedFlagValue(cmd, "api-key"); value != "" {
		opts.Token = value
	}
	if value := inheritedFlagValue(cmd, "timeout"); value != "" {
		opts.Timeout = value
	}
	if value := inheritedFlagValue(cmd, "output"); value != "" {
		opts.Format = value
	}
	if flag := cmd.Flags().Lookup("insecure"); flag != nil && flag.Value.String() == "true" {
		opts.InsecureSkipVerify = true
	}
}

func inheritedFlagValue(cmd *cobra.Command, name string) string {
	if flag := cmd.Flags().Lookup(name); flag != nil {
		return flag.Value.String()
	}
	if flag := cmd.InheritedFlags().Lookup(name); flag != nil {
		return flag.Value.String()
	}
	return ""
}

func Execute(opts RootOptions) error {
	cmd := NewRootCommand(opts)
	if err := cmd.Execute(); err != nil {
		return fmt.Errorf("chaitin-cli failed: %w", err)
	}
	return nil
}
