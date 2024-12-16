package cli

import "github.com/spf13/cobra"

func NewCmdInspectConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "",
		Long:    "",
		Example: "",
		Args:    cobra.MinimumNArgs(1),
		RunE:    inspectConfig,
	}
	return cmd
}

func inspectConfig(cmd *cobra.Command, args []string) error {
	return nil
}
