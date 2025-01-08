package cli

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func NewCmdInspectConnection() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "connection",
		Short:   "",
		Long:    "",
		Example: "",
		Args:    cobra.MinimumNArgs(1),
		RunE:    cmdInspectConnection,
	}
	return cmd
}

func cmdInspectConnection(cmd *cobra.Command, args []string) error {
	slog.Warn("not implemented")
	return nil
}
