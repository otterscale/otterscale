package main

import (
	"github.com/google/wire"

	"github.com/openhdc/openhdc/pkg/app"
	"github.com/openhdc/openhdc/pkg/connector"
	"github.com/openhdc/openhdc/pkg/transport"

	_ "go.uber.org/automaxprocs"
)

var ProviderSet = wire.NewSet(
	newConnector,
	connector.NewDestinationAdapter,
	connector.NewTransportServer,
)

func newApp(srv *transport.Server) *app.App {
	return app.New(
		app.Servers(srv),
	)
}

func main() {
	// // create zap logger
	// zlog := zap.NewServiceZap(id, name, version, env, key, host, ip)
	// logger := zap.NewLogger(zlog)
	// defer func() { _ = logger.Sync() }()

	// // set global logger
	// log.SetLogger(logger)

	// // get config from consul
	// cfg := get.NewConfig(key)
	// defer func() { _ = cfg.Close() }()

	// wire app
	app, cleanup, err := wireApp()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
