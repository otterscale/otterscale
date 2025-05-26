package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/openhdc/otterscale/internal/config"
)

func NewInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Initialize a new Otterscale configuration",
		Long:    "Initialize a new Otterscale configuration in the specified directory. This creates default configuration files and directory structure required for Otterscale to operate.",
		Example: "otterscale init /path/to/config/dir/my-config.yaml\notterscale init ./my-config.yaml",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Info("Initializing configuration file", "path", args[0])
			return config.InitFile(args[0])
		},
	}
	return cmd
}
