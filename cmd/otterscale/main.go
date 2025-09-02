package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	oscmd "github.com/otterscale/otterscale/internal/cmd"
	"github.com/otterscale/otterscale/internal/config"
)

var version = "devel"

func newCmd(conf *config.Config, mux *http.ServeMux) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "otterscale",
		Short:         "Open-source platform for managing server infrastructure and GPU farms",
		Long:          "OtterScale is an open-source platform that integrates MAAS, Juju, Kubernetes, and Ceph to provide comprehensive management of servers and GPU farms. It simplifies hardware provisioning, resource allocation, and infrastructure orchestration for data centers and compute clusters.",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(
		oscmd.NewInit(),
		oscmd.NewServe(conf, mux),
	)
	return cmd
}

func run() error {
	// options
	grpcHelper := true

	// wire cmd
	cmd, cleanup, err := wireCmd(grpcHelper)
	if err != nil {
		return err
	}
	defer cleanup()

	// start and wait for stop signal
	return cmd.Execute()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
