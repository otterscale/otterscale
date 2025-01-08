package main

import (
	"flag"
	"strings"
	"time"

	"github.com/google/wire"

	"github.com/openhdc/openhdc"
	"github.com/openhdc/openhdc/api/property/v1"
	"github.com/openhdc/openhdc/connector/jira/client"

	_ "go.uber.org/automaxprocs"
)

const name = "jira"

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
	syncMode = flag.String("sync_mode", "", "sync mode, such as full_overwrite, full_append, incremental_append, incremental_append_dedupe")
	cursor   = flag.String("cursor", "", "incremental cursor")

	// Jira setting
	server    = flag.String("server", "", "connection string, such as 'postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10'")
	username  = flag.String("username", "", "User name.")
	password  = flag.String("password", "", "User Password.")
	token     = flag.String("token", "", "Jira API Token. See the https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/ for more information on how to generate token key.")
	projects  = flag.String("projects", "", "List of project keys or empty if you want to sync for all projects.")
	startDate = flag.String("start_date", "", "The date you want to sync data from jira, use the format YYYY-MM-DDT00:00:00Z.")

	// write
	batchSize      = flag.Int64("batch_size", defaultBatchSize, "default batch size of rows is 10,000 if not specified")
	batchSizeBytes = flag.Int64("batch_size_bytes", defaultBatchSizeBytes, "default batch size of bytes is 10,000,000 bytes if not specified")
	batchTimeout   = flag.Duration("batch_timeout", defaultBatchTimeout, "default batch timeout is 60s if not specified")
	createIndex    = flag.Bool("create_index", true, "create an index to improve performance")
)

var ProviderSet = wire.NewSet(openhdc.NewServer, openhdc.NewService, client.NewConnector)

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
		client.WithName(name),
		client.WithServer(*server),
		client.WithUsername(*username),
		client.WithPassword(*password),
		client.WithToken(*token),
		client.WithProjects(strings.Split(*projects, ",")),
		client.WithStartDate(*startDate),
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
