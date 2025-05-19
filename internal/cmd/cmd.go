package cmd

import (
	"github.com/google/wire"
	"github.com/spf13/cobra"

	"github.com/openhdc/otterscale/internal/app"
)

var ProviderSet = wire.NewSet(New)

func New(version string, na *app.NexusApp) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "openhdc",
		Short:        "",
		Long:         "",
		Version:      version,
		SilenceUsage: true,
	}
	cmd.AddCommand(
		NewCmdInit(),
		NewCmdServe(na),
	)
	return cmd
}
