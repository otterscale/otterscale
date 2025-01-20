package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/openhdc/openhdc"
	"github.com/openhdc/openhdc/api/property/v1"
	"github.com/openhdc/openhdc/connectors/github/client"

	_ "go.uber.org/automaxprocs"
)

var (
	name    = "csv"
	version = "devel"
)

const (
	defaultBatchSize      = 10000
	defaultBatchSizeBytes = 100000000
	defaultBatchTimeout   = 60 * time.Second
)

var (
	// print
	printVersion = flag.Bool("v", false, "print version")

	// general
	kind    = flag.String("kind", "", "connector type, such as source or destination")
	network = flag.String("network", "tcp", "network of grpc server")
	address = flag.String("address", ":0", "address of grpc server")

	// read
	owner = flag.String("owner", "", "github code owner")
	repo  = flag.String("repo", "", "github code repositories")
	opts  = flag.String("opts", "", "functional options")

	// write
	batchSize      = flag.Int64("batch_size", defaultBatchSize, "default batch size of rows is 10,000 if not specified")
	batchSizeBytes = flag.Int64("batch_size_bytes", defaultBatchSizeBytes, "default batch size of bytes is 10,000,000 bytes if not specified")
	batchTimeout   = flag.Duration("batch_timeout", defaultBatchTimeout, "default batch timeout is 60s if not specified")
	createIndex    = flag.Bool("create_index", true, "create an index to improve performance")
)

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

	// version
	if *printVersion {
		fmt.Println(name, version)
		return
	}

	// options
	serverOpts := []openhdc.ServerOption{
		openhdc.WithNetwork(*network),
		openhdc.WithAddress(*address),
	}

	clientOpts := []client.Option{
		client.WithName(name),
		client.WithOwner(*owner),
		client.WithRepo(*repo),
		client.WithOpts(*opts),
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
