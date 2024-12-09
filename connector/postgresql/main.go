package main

import (
	"github.com/google/wire"

	"github.com/openhdc/openhdc/connector/postgresql/client"
	"github.com/openhdc/openhdc/connector/postgresql/pgarrow"
	"github.com/openhdc/openhdc/internal/app"
	"github.com/openhdc/openhdc/internal/connector"

	_ "go.uber.org/automaxprocs"
)

var ProviderSet = wire.NewSet(pgarrow.NewCodec, client.NewConnector, connector.NewService, connector.NewServer)

func newApp(srv *connector.Server) *app.App {
	return app.New(
		app.WithServers(srv),
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

	// TODO: OPTIONS FROM FLAG
	clientOpts := []client.Option{}

	connectorOpts := []connector.Option{}

	serverOpts := []connector.ServerOption{
		connector.WithNetwork("unix"),
		connector.WithAddress(":0"),
	}

	// wire app
	app, cleanup, err := wireApp(clientOpts, connectorOpts, serverOpts)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
