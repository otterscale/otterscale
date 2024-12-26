package cmd

import "github.com/spf13/cobra"

func NewCmdInspectConnection() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "connection",
		Short:   "",
		Long:    "",
		Example: "",
		Args:    cobra.MinimumNArgs(1),
		RunE:    inspectConnection,
	}
	return cmd
}

func inspectConnection(cmd *cobra.Command, args []string) error {
	return nil
}
