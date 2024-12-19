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
	// start and wait for stop signal
	signals := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM}
	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()
	_ = cmd.NewCmdRoot(version).ExecuteContext(ctx)
}
