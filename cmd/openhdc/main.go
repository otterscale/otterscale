package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/openhdc/openhdc/internal/cmd"
)

var version = "devel"

func main() {
	signals := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM}
	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()
	if err := cmd.NewCmdRoot(version).ExecuteContext(ctx); err != nil {
		panic(err)
	}
}
