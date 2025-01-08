package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/openhdc/openhdc/internal/cli"
)

var (
	name    = "openhdc-cli"
	version = "devel"
)

// print
var printVersion = flag.Bool("v", false, "print version")

func run() error {
	signals := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM}
	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()
	return cli.NewCmdRoot(version).ExecuteContext(ctx)
}

func main() {
	flag.Parse()

	// version
	if *printVersion {
		fmt.Println(name, version)
		return
	}

	// start and wait for stop signal
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
