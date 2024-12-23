package main

import (
	"flag"
	"time"

	"github.com/google/wire"

	"github.com/openhdc/openhdc"
	"github.com/openhdc/openhdc/api/property/v1"
	"github.com/openhdc/openhdc/connector/postgresql/client"

	_ "go.uber.org/automaxprocs"
)

const (
	defaultBatchSize      = 10000
	defaultBatchSizeBytes = 100000000
	defaultBatchTimeout   = 60 * time.Second
)

var (
	kind    = flag.String("kind", "", "connector type, such as source or destination")
	network = flag.String("network", "tcp", "network of grpc server")
	address = flag.String("address", ":0", "address of grpc server")

	syncMode       = flag.String("sync_mode", "", "sync mode, such as full_overwrite, full_append, incremental_append or incremental_append_dedupe")
	batchSize      = flag.Int64("batch_size", defaultBatchSize, "")            // TODO: USAGE
	batchSizeBytes = flag.Int64("batch_size_bytes", defaultBatchSizeBytes, "") // TODO: USAGE
	batchTimeout   = flag.Duration("batch_timeout", defaultBatchTimeout, "")   // TODO: USAGE

	connString = flag.String("connection_string", "", "connection string, such as 'postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10'")
	namespace  = flag.String("namespace", "", "namespace of database, such as 'public'")

	createIndex = flag.Bool("create_index", true, "") // TODO: USAGE
)

var ProviderSet = wire.NewSet(openhdc.NewServer, openhdc.NewService, client.NewConnector, client.NewWriter, openhdc.NewDefaultCodec)

func newApp(srv *openhdc.Server) *openhdc.App {
	return openhdc.New(
		openhdc.WithServers(srv),
		openhdc.WithKind(property.ParseWorkloadKind(*kind)),
	)
}

func main() {
	flag.Parse()

	serverOpts := []openhdc.ServerOption{
		openhdc.WithNetwork(*network),
		openhdc.WithAddress(*address),
	}

	clientOpts := []client.Option{
		client.WithSyncMode(property.ParseSyncMode(*syncMode)),
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
