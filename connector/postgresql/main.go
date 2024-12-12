package main

import (
	"flag"
	"time"

	"github.com/google/wire"

	"github.com/openhdc/openhdc/connector/postgresql/client"
	"github.com/openhdc/openhdc/connector/postgresql/pgarrow"
	"github.com/openhdc/openhdc/internal/app"
	"github.com/openhdc/openhdc/internal/connector"

	_ "go.uber.org/automaxprocs"
)

var (
	kind    = flag.String("kind", "", "")
	network = flag.String("network", "tcp", "")
	address = flag.String("address", ":0", "")

	batchSize      = flag.Int64("batch_size", 10000, "")
	batchSizeBytes = flag.Int64("batch_size_bytes", 100000000, "")
	batchTimeout   = flag.Duration("batch_timeout", 60*time.Second, "")
	createIndex    = flag.Bool("create_index", true, "")
	connString     = flag.String("connection_string", "", "")
)

var ProviderSet = wire.NewSet(connector.NewServer, connector.NewService, client.NewConnector, pgarrow.NewCodec)

func newApp(srv *connector.Server) *app.App {
	return app.New(
		app.WithServers(srv),
	)
}

func main() {
	flag.Parse()

	serverOpts := []connector.ServerOption{
		connector.WithNetwork(*network),
		connector.WithAddress(*address),
	}

	connectorOpts := []connector.Option{
		connector.WithKind(connector.Kind(*kind)),
	}

	clientOpts := []client.Option{
		client.WithBatchSize(*batchSize),
		client.WithBatchSizeBytes(*batchSizeBytes),
		client.WithBatchTimeout(*batchTimeout),
		client.WithCreateIndex(*createIndex),
		client.WithConnString(*connString),
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
