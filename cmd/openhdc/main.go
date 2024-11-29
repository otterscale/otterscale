package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/openhdc/openhdc/internal/cmd"
)

var version = "devel"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	defer stop()
	if err := cmd.NewCmdRoot(version).ExecuteContext(ctx); err != nil {
		panic(err)
	}
}
