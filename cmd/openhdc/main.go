package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/openhdc/openhdc/internal/cli"
)

var version = "devel"

func main() {
	// start and wait for stop signal
	signals := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM}
	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()
	_ = cli.NewCmdRoot(version).ExecuteContext(ctx)
}
