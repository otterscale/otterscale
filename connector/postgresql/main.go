package main

import (
	"github.com/google/wire"

	"github.com/openhdc/openhdc/connector/postgresql/client"
	"github.com/openhdc/openhdc/pkg/app"
	"github.com/openhdc/openhdc/pkg/connector"
	"github.com/openhdc/openhdc/pkg/transport"

	_ "go.uber.org/automaxprocs"
)

var ProviderSet = wire.NewSet(
	client.NewCodec,
	client.NewAdapter,
	connector.New,
)

func newApp(c *connector.Connector) *app.App {
	return app.New(
		app.Servers(c.Server),
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

	serverOpts := []transport.ServerOption{
		transport.Network("unix"),
		transport.Address(":0"),
	}

	connectorOpts := []connector.Option{}

	// wire app
	app, cleanup, err := wireApp(clientOpts, serverOpts, connectorOpts)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
