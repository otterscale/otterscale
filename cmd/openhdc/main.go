package main

import (
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"github.com/openhdc/openhdc/internal/cli"
	_ "github.com/openhdc/openhdc/internal/migrations"
)

var version = "devel"

func newApp(fns []func(se *core.ServeEvent) error) *pocketbase.PocketBase {
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

	// set functions
	for _, fn := range fns {
		app.OnServe().BindFunc(fn)
		app.OnRecordAfterCreateSuccess()
	}

	return app
}

func main() {
	// wire app
	app, cleanup, err := wireApp()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Start(); err != nil {
		panic(err)
	}
}
