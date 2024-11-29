package cmd

import "github.com/spf13/cobra"

func NewCmdInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "",
		Long:    "",
		Example: "",
		Args:    cobra.MinimumNArgs(1),
		RunE:    _init,
	}
	return cmd
}

func _init(cmd *cobra.Command, args []string) error {
	return nil
}
