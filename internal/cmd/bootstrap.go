package cmd

import (
	"github.com/spf13/cobra"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/mux"
)

func NewBootstrap(bootstrap *mux.Bootstrap) *cobra.Command {
	cmd := &cobra.Command{
		Use: "bootstrap",
		// Short:   "Initialize a new OtterScale configuration",
		// Long:    "Initialize a new OtterScale configuration by printing the default configuration to stdout. This outputs the default configuration that can be redirected to a file.",
		// Example: "otterscale init > config.yaml\notterscale init > /path/to/config.yaml",
		RunE: func(_ *cobra.Command, _ []string) error {
			return config.PrintDefaultConfig()
		},
	}
	return cmd
}
