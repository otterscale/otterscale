package cmd

import (
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/otterscale/otterscale/internal/config"
)

const (
	containerEnvVar            = "OTTERSCALE_CONTAINER"
	defaultContainerAddress    = ":8299"
	defaultContainerConfigPath = "/etc/app/otterscale.yaml"
)

func NewServe(conf *config.Config, mux *http.ServeMux) *cobra.Command {
	var address, configPath string

	cmd := &cobra.Command{
		Use:     "serve",
		Short:   "Start the OtterScale API server",
		Long:    "Start the OtterScale API server that provides gRPC and HTTP endpoints for all services",
		Example: "otterscale serve --address=:8299 --config=otterscale.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := maxprocs.Set(maxprocs.Logger(log.Printf))
			if err != nil {
				slog.Error("Error setting GOMAXPROCS", "err", err)
			}

			// Check if running in container and override address
			if os.Getenv(containerEnvVar) != "" {
				address = defaultContainerAddress
				configPath = defaultContainerConfigPath
				slog.Info("Container environment detected, using default configuration", "address", address, "config", configPath)
			}

			slog.Info("Loading configuration file", "path", configPath)
			if err := conf.Load(configPath); err != nil {
				return err
			}

			srv := &http.Server{
				Addr: address,
				Handler: h2c.NewHandler(
					cors.AllowAll().Handler(mux),
					&http2.Server{},
				),
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       5 * time.Minute,
				WriteTimeout:      5 * time.Minute,
				MaxHeaderBytes:    8 * 1024, // 8KiB
			}

			listener, err := net.Listen("tcp", address)
			if err != nil {
				return err
			}

			slog.Info("Server starting on", "address", listener.Addr().String())
			return srv.Serve(listener)
		},
	}

	cmd.Flags().StringVarP(
		&address,
		"address",
		"a",
		":0",
		"address of service",
	)

	cmd.Flags().StringVarP(
		&configPath,
		"config",
		"c",
		"otterscale.yaml",
		"config path",
	)

	return cmd
}
