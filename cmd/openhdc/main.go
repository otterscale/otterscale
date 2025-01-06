package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/openhdc/openhdc/internal/cmd"
)

var version = "devel"

// start and wait for stop signal
func run() error {
	signals := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM}
	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()
	return cmd.NewCmdRoot(version).ExecuteContext(ctx)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
