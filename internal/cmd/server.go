package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/otterscale/otterscale/internal/cmd/server"
	"github.com/otterscale/otterscale/internal/config"
)

// ServerInjector is the Wire-generated Server factory, with its cleanup.
type ServerInjector func() (*server.Server, func(), error)

// NewServerCommand calls the injector lazily inside RunE, so expensive setup
// such as OIDC provider discovery only happens when the command runs.
func NewServerCommand(conf *config.Config, newServer ServerInjector) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start server that provides gRPC and HTTP endpoints for the core services",
		Long: "Start server that provides gRPC and HTTP endpoints for the core services.\n\n" +
			"Tunnel state is held in memory, so run exactly one replica: a second one would " +
			"keep its own registry and its own CA, and agents registered against one replica " +
			"are unreachable through the other.\n\n" +
			"The tunnel CA is generated at startup and never persisted. Agent certificates " +
			"issued before a restart stop being trusted; agents re-register automatically, so " +
			"expect their clusters to be briefly unreachable after every restart.",
		Example: "otterscale server --address=:8299 --tunnel-address=127.0.0.1:8300",
		RunE: func(cmd *cobra.Command, _ []string) error {
			srv, cleanup, err := newServer()
			if err != nil {
				return fmt.Errorf("failed to initialize server: %w", err)
			}
			defer cleanup()

			cfg := &server.Config{
				Address:           conf.ServerAddress(),
				AllowedOrigins:    conf.ServerAllowedOrigins(),
				TunnelAddress:     conf.ServerTunnelAddress(),
				ExternalTunnelURL: conf.ServerExternalTunnelURL(),
				KeycloakRealmURL:  conf.ServerKeycloakRealmURL(),
				KeycloakClientID:  conf.ServerKeycloakClientID(),
			}

			return srv.Run(cmd.Context(), cfg)
		},
	}

	if err := conf.BindFlags(cmd.Flags(), config.ServerOptions); err != nil {
		return nil, err
	}

	return cmd, nil
}
