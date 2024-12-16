package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	cmd "github.com/openhdc/openhdc/internal/cli"
)

var version = "devel"

func main() {
	// new logger

	// set global

	// wire app

	// start and wait for stop signal
	signals := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM}
	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()
	if err := cmd.NewCmdRoot(version).ExecuteContext(ctx); err != nil {
		panic(err)
	}
}
