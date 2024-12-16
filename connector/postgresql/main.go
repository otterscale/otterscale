package main

import (
	"flag"
	"time"

	"github.com/google/wire"

	"github.com/openhdc/openhdc"
	"github.com/openhdc/openhdc/api/workload/v1"
	"github.com/openhdc/openhdc/connector/postgresql/client"

	_ "go.uber.org/automaxprocs"
)

const (
	defaultBatchSize      = 10000
	defaultBatchSizeBytes = 100000000
	defaultBatchTimeout   = 60 * time.Second
)

var (
	kind    = flag.String("kind", "", "")
	network = flag.String("network", "tcp", "")
	address = flag.String("address", ":0", "")

	batchSize      = flag.Int64("batch_size", defaultBatchSize, "")
	batchSizeBytes = flag.Int64("batch_size_bytes", defaultBatchSizeBytes, "")
	batchTimeout   = flag.Duration("batch_timeout", defaultBatchTimeout, "")
	createIndex    = flag.Bool("create_index", true, "")
	connString     = flag.String("connection_string", "", "")
	namespace      = flag.String("namespace", "", "")
)

var ProviderSet = wire.NewSet(openhdc.NewServer, openhdc.NewService, client.NewConnector)

func newApp(srv *openhdc.Server) *openhdc.App {
	return openhdc.New(
		openhdc.WithServers(srv),
		openhdc.WithKind(workload.ParseKind(*kind)),
	)
}

func main() {
	flag.Parse()

	serverOpts := []openhdc.ServerOption{
		openhdc.WithNetwork(*network),
		openhdc.WithAddress(*address),
	}

	clientOpts := []client.Option{
		client.WithBatchSize(*batchSize),
		client.WithBatchSizeBytes(*batchSizeBytes),
		client.WithBatchTimeout(*batchTimeout),
		client.WithCreateIndex(*createIndex),
		client.WithConnString(*connString),
		client.WithNamespace(*namespace),
	}

	// wire app
	app, cleanup, err := wireApp(serverOpts, clientOpts)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
