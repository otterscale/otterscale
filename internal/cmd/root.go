package cmd

import "github.com/spf13/cobra"

func NewCmdRoot(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "openhdc",
		Short:        "",
		Long:         "",
		Version:      version,
		SilenceUsage: true,
	}
	cmd.AddCommand(
		NewCmdInit(),
		NewCmdInspect(),
		NewCmdSync(),
	)
	return cmd
}
