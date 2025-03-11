package service

import (
	"os"
	"strings"

	"github.com/google/wire"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"github.com/openhdc/openhdc/internal/cli"
	"github.com/openhdc/openhdc/internal/service/app"
)

var ProviderSet = wire.NewSet(NewPocketBase)

func NewPocketBase(version string, pa *app.PipelineApp) *pocketbase.PocketBase {
	// initialize app
	app := pocketbase.New()

	// set version
	app.RootCmd.Version = version

	// set commands
	app.RootCmd.AddCommand(
		cli.NewCmdInit(),
		cli.NewCmdInspect(),
		cli.NewCmdSync(),
	)

	// set migration command
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})

	// bind functions
	pa.Bind(app)

	return app
}
