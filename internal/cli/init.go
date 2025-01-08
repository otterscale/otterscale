package cli

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func NewCmdInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "",
		Long:    "",
		Example: "",
		Args:    cobra.MinimumNArgs(1),
		RunE:    cmdInit,
	}
	return cmd
}

func cmdInit(cmd *cobra.Command, args []string) error {
	slog.Warn("not implemented")
	return nil
}
