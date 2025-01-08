package cli

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func NewCmdInspectConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "",
		Long:    "",
		Example: "",
		Args:    cobra.MinimumNArgs(1),
		RunE:    cmdInspectConfig,
	}
	return cmd
}

func cmdInspectConfig(cmd *cobra.Command, args []string) error {
	slog.Warn("not implemented")
	return nil
}
