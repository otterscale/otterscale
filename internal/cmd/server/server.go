// Package server implements the control-plane runtime that serves the
// public gRPC/HTTP API and manages the chisel tunnel listener.
package server

import (
	"context"
	"fmt"
	"net"

	linkv1 "github.com/otterscale/api/link/v1"

	"github.com/otterscale/otterscale/internal/transport"
	"github.com/otterscale/otterscale/internal/transport/http"
)

// Config holds the runtime parameters for a Server.
type Config struct {
	Address          string
	AllowedOrigins   []string
	TunnelAddress    string
	KeycloakRealmURL string
	KeycloakClientID string
}

// BackgroundListeners is a slice of transport.Listener that
// participate in the managed lifecycle alongside the HTTP and tunnel
// servers. This named type exists to enable Wire injection of
// multiple background tasks (session reaper, cache evictor, etc.)
// without the Server depending on their concrete types.
type BackgroundListeners []transport.Listener

// Server binds an HTTP server (gRPC + REST) and a chisel tunnel
// listener, running them in parallel via transport.Serve.
type Server struct {
	handler    *Handler
	tunnel     transport.TunnelService
	background BackgroundListeners
}

// NewServer returns a Server wired to the given handler, tunnel
// service, and background listeners. The TunnelService interface
// decouples the server from concrete tunnel implementations, keeping
// infrastructure details behind the interface boundary.
func NewServer(handler *Handler, tunnel transport.TunnelService, background BackgroundListeners) *Server {
	return &Server{handler: handler, tunnel: tunnel, background: background}
}

// Run starts both the HTTP and tunnel servers. It blocks until ctx
// is canceled or an unrecoverable error occurs. Health, reflection,
// and link-registration endpoints are marked as public (no auth).
func (s *Server) Run(ctx context.Context, cfg *Config) error {
	if cfg.KeycloakRealmURL == "" {
		return fmt.Errorf("keycloak realm URL is required but not configured")
	}

	// Parse the tunnel address to extract the host for the TLS
	// certificate SAN.
	tunnelHost, _, err := net.SplitHostPort(cfg.TunnelAddress)
	if err != nil {
		return fmt.Errorf("parse tunnel address %q: %w", cfg.TunnelAddress, err)
	}

	oidc, err := http.NewOIDC(cfg.KeycloakRealmURL, cfg.KeycloakClientID)
	if err != nil {
		return fmt.Errorf("failed to create OIDC middleware: %w", err)
	}

	httpSrv, err := http.NewServer(
		ctx,
		http.WithAddress(cfg.Address),
		http.WithAllowedOrigins(cfg.AllowedOrigins),
		http.WithAuthMiddleware(oidc),
		http.WithPublicPaths([]string{
			"/grpc.health.v1.Health/Check",
			"/grpc.health.v1.Health/Watch",
			"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo",
			linkv1.LinkServiceRegisterProcedure,
		}),
		http.WithPublicPathPrefixes([]string{
			"/link/manifest/",
		}),
		http.WithMount(s.handler.Mount),
	)
	if err != nil {
		return fmt.Errorf("failed to create HTTP server: %w", err)
	}

	// Build the tunnel server listener with mTLS via the injected
	// TunnelService. Certificate generation and file I/O are
	// encapsulated behind the interface.
	tunnelSrv, err := s.tunnel.BuildTunnelListener(cfg.TunnelAddress, tunnelHost)
	if err != nil {
		return fmt.Errorf("failed to create tunnel server: %w", err)
	}

	// Detect disconnected tunnel clients and remove stale
	// registrations.
	healthChecker := s.tunnel.BuildHealthListener()

	listeners := []transport.Listener{httpSrv, tunnelSrv, healthChecker}
	listeners = append(listeners, s.background...)

	return transport.Serve(ctx, listeners...)
}
