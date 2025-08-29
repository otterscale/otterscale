package cmd

import (
	"github.com/spf13/cobra"

	"github.com/otterscale/otterscale/internal/config"
)

func NewInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Initialize a new Otterscale configuration",
		Long:    "Initialize a new Otterscale configuration by printing the default configuration to stdout. This outputs the default configuration that can be redirected to a file.",
		Example: "otterscale init > config.yaml\notterscale init > /path/to/config.yaml",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return config.PrintDefaultConfig()
		},
	}
	return cmd
}
