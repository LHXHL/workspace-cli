package main

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
)

// version is overridden at build time with
// `-ldflags "-X main.version=v1.2.3"`. Local builds that do not inject it fall
// back to the metadata the Go toolchain records in the binary.
var version string

const devVersion = "dev"

// readBuildInfo is a seam so tests can control the fallback path.
var readBuildInfo = debug.ReadBuildInfo

func resolveVersion() string {
	if injected := strings.TrimSpace(version); injected != "" {
		return injected
	}
	return versionFromBuildInfo(readBuildInfo)
}

func versionFromBuildInfo(read func() (*debug.BuildInfo, bool)) string {
	info, ok := read()
	if !ok || info == nil {
		return devVersion
	}
	if moduleVersion := strings.TrimSpace(info.Main.Version); moduleVersion != "" && moduleVersion != "(devel)" {
		return moduleVersion
	}
	for _, setting := range info.Settings {
		if setting.Key != "vcs.revision" {
			continue
		}
		if revision := strings.TrimSpace(setting.Value); revision != "" {
			return devVersion + "+" + shortRevision(revision)
		}
	}
	return devVersion
}

func shortRevision(revision string) string {
	const maxLength = 12
	if len(revision) > maxLength {
		return revision[:maxLength]
	}
	return revision
}

func newVersionCommand(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the chaitin-cli version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), root.Version)
			return nil
		},
	}
}
