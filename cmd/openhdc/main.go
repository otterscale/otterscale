package main

import (
	"flag"

	"github.com/openhdc/openhdc"
)

var version = "devel"

var (

	// general
	network = flag.String("network", "tcp", "network of grpc server")
	address = flag.String("address", ":0", "address of grpc server")
)

func main() {
	// options
	serverOpts := []openhdc.ServerOption{
		openhdc.WithNetwork(*network),
		openhdc.WithAddress(*address),
	}

	// wire app
	app, cleanup, err := wireApp(version, serverOpts)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Execute(); err != nil {
		panic(err)
	}
}
