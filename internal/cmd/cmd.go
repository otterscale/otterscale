package cmd

import (
	"github.com/spf13/cobra"

	"github.com/openhdc/otterscale/internal/app"
)

func New(version string, app *app.ApplicationService, config *app.ConfigurationService, facility *app.FacilityService, general *app.GeneralService, machine *app.MachineService, network *app.NetworkService, scope *app.ScopeService, tag *app.TagService) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "openhdc",
		Short:        "",
		Long:         "",
		Version:      version,
		SilenceUsage: true,
	}
	cmd.AddCommand(
		NewCmdInit(),
		NewCmdServe(app, config, facility, general, machine, network, scope, tag),
	)
	return cmd
}
