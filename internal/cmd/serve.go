package cmd

import (
	"log/slog"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/spf13/cobra"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/mux"
	"github.com/otterscale/otterscale/internal/mux/interceptor/impersonation"
)

func NewServe(conf *config.Config, serve *mux.Serve) *cobra.Command {
	var address, configPath string

	cmd := &cobra.Command{
		Use:     "serve",
		Short:   "Start the main API server",
		Long:    "Start the OtterScale API server that provides gRPC and HTTP endpoints for the core services (used after bootstrap)",
		Example: "otterscale serve --address=:8299 --config=otterscale.yaml",
		RunE: func(_ *cobra.Command, _ []string) error {
			if os.Getenv(containerEnvVar) != "" {
				address = defaultContainerAddress
				configPath = defaultContainerConfigPath
				slog.Info("Container environment detected, using default configuration", "address", address, "config", configPath)
			}

			slog.Info("Loading configuration file", "path", configPath)
			if err := conf.Load(configPath); err != nil {
				return err
			}

			opts, err := newServeOptions(conf)
			if err != nil {
				return err
			}
			return startHTTPServer(address, serve, opts...)
		},
	}

	cmd.Flags().StringVarP(
		&address,
		"address",
		"a",
		":0",
		"Address for server to listen on",
	)

	cmd.Flags().StringVarP(
		&configPath,
		"config",
		"c",
		"otterscale.yaml",
		"Config path for server to load",
	)

	return cmd
}

func newServeOptions(conf *config.Config) ([]connect.HandlerOption, error) {
	openTelemetryInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return nil, err
	}

	impersonationInterceptor, err := impersonation.NewInterceptor(conf)
	if err != nil {
		return nil, err
	}

	return []connect.HandlerOption{
		connect.WithInterceptors(openTelemetryInterceptor, impersonationInterceptor),
	}, nil
}
