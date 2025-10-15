package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/otterscale/otterscale/internal/mux"
)

func NewBootstrap(bootstrap *mux.Bootstrap) *cobra.Command {
	var address string

	cmd := &cobra.Command{
		Use:     "bootstrap",
		Short:   "Start the bootstrap API server",
		Long:    "Start the OtterScale API server that provides gRPC and HTTP endpoints for bootstrap service",
		Example: "otterscale bootstrap --address=:8299",
		RunE: func(_ *cobra.Command, _ []string) error {
			if os.Getenv(containerEnvVar) != "" {
				address = defaultContainerAddress
				slog.Info("Container environment detected, using default configuration", "address", address)
			}
			return startHTTPServer(address, bootstrap)
		},
	}

	cmd.Flags().StringVarP(
		&address,
		"address",
		"a",
		":0",
		"address of service",
	)

	return cmd
}
