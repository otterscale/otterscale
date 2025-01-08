package main

import (
	"flag"

	"github.com/openhdc/openhdc"

	_ "go.uber.org/automaxprocs"
)

const (
	name    = "openhdc-server"
	version = "devel"
)

var (
	network = flag.String("network", "tcp", "network of grpc server")
	address = flag.String("address", ":0", "address of grpc server")
)

func newApp(srv *openhdc.Server) *openhdc.App {
	return openhdc.New(
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

	// wire app
	app, cleanup, err := wireApp(serverOpts)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
