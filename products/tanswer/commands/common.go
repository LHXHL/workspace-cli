package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func domainCommand(use, short string, children ...string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
	}
	for _, child := range children {
		cmd.AddCommand(leafCommand(child))
	}
	return cmd
}

func leafCommand(use string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: fmt.Sprintf("%s semantic command", use),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}
