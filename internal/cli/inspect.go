package cli

import "github.com/spf13/cobra"

func NewCmdInspect() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inspect",
		Short:   "",
		Long:    "",
		Example: "",
		Args:    cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(
		NewCmdInspectConfig(),
		NewCmdInspectConnection(),
	)
	return cmd
}
