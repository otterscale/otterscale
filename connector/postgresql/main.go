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
	name    = "postgresql"
	version = "devel"
)

const (
	defaultBatchSize      = 10000
	defaultBatchSizeBytes = 100000000
	defaultBatchTimeout   = 60 * time.Second
)

var (
	// general
	kind    = flag.String("kind", "", "connector type, such as source or destination")
	network = flag.String("network", "tcp", "network of grpc server")
	address = flag.String("address", ":0", "address of grpc server")

	// read
	connString = flag.String("connection_string", "", "connection string, such as 'postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10'")
	namespace  = flag.String("namespace", "", "namespace of database, such as 'public'")

	// write
	batchSize      = flag.Int64("batch_size", defaultBatchSize, "default batch size of rows is 10,000 if not specified")
	batchSizeBytes = flag.Int64("batch_size_bytes", defaultBatchSizeBytes, "default batch size of bytes is 10,000,000 bytes if not specified")
	batchTimeout   = flag.Duration("batch_timeout", defaultBatchTimeout, "default batch timeout is 60s if not specified")
	createIndex    = flag.Bool("create_index", true, "create an index to improve performance")
)

var ProviderSet = wire.NewSet(openhdc.NewServer, openhdc.NewService, client.NewConnector, openhdc.NewDefaultCodec)

func newApp(srv *openhdc.Server) *openhdc.App {
	return openhdc.New(
		openhdc.WithKind(property.ParseWorkloadKind(*kind)),
		openhdc.WithName(name),
		openhdc.WithVersion(version),
		openhdc.WithServers(srv),
	)
}

func main() {
	flag.Parse()

	serverOpts := []openhdc.ServerOption{
		openhdc.WithNetwork(*network),
		openhdc.WithAddress(*address),
	}

	clientOpts := []client.Option{
		client.WithName(name),
		client.WithConnString(*connString),
		client.WithNamespace(*namespace),
		client.WithBatchSize(*batchSize),
		client.WithBatchSizeBytes(*batchSizeBytes),
		client.WithBatchTimeout(*batchTimeout),
		client.WithCreateIndex(*createIndex),
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
