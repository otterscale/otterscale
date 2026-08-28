// Package server implements the control-plane runtime that serves the
// public gRPC/HTTP API and manages the chisel tunnel listener.
package server

import (
	"context"
	"fmt"
	"net"
	"net/url"

	linkv1 "github.com/otterscale/otterscale/api/link/v1"

	"github.com/otterscale/otterscale/internal/transport"
	"github.com/otterscale/otterscale/internal/transport/http"
)

type Config struct {
	Address           string
	AllowedOrigins    []string
	TunnelAddress     string
	ExternalTunnelURL string
	KeycloakRealmURL  string
	KeycloakClientID  string
}

// BackgroundListeners run in the managed lifecycle alongside the HTTP and
// tunnel servers. The named type lets Wire inject several background tasks
// (session reaper, cache evictor) without the Server knowing their types.
type BackgroundListeners []transport.Listener

// Server runs the HTTP server (gRPC + REST) and the chisel tunnel listener in
// parallel via transport.Serve.
type Server struct {
	handler    *Handler
	tunnel     transport.TunnelService
	background BackgroundListeners
}

// NewServer takes tunnel as an interface, keeping the concrete tunnel
// implementation behind the boundary.
func NewServer(handler *Handler, tunnel transport.TunnelService, background BackgroundListeners) *Server {
	return &Server{handler: handler, tunnel: tunnel, background: background}
}

// Run blocks until ctx is canceled or an unrecoverable error occurs. Health,
// reflection, and link-registration endpoints are public (no auth).
func (s *Server) Run(ctx context.Context, cfg *Config) error {
	if cfg.KeycloakRealmURL == "" {
		return fmt.Errorf("keycloak realm URL is required but not configured")
	}

	tunnelHost, err := resolveTunnelHost(cfg)
	if err != nil {
		return err
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
		http.WithLongRunningPaths(s.handler.LongRunningPaths()),
		http.WithMount(s.handler.Mount),
	)
	if err != nil {
		return fmt.Errorf("failed to create HTTP server: %w", err)
	}

	// TunnelService keeps certificate generation and file I/O behind the
	// interface.
	tunnelSrv, err := s.tunnel.BuildTunnelListener(cfg.TunnelAddress, tunnelHost)
	if err != nil {
		return fmt.Errorf("failed to create tunnel server: %w", err)
	}

	// Detects disconnected clients and removes stale registrations.
	healthChecker := s.tunnel.BuildHealthListener()

	listeners := []transport.Listener{httpSrv, tunnelSrv, healthChecker}
	listeners = append(listeners, s.background...)

	return transport.Serve(ctx, listeners...)
}

// resolveTunnelHost returns the SAN for the tunnel server's TLS certificate.
// Agents dial with the CA pinned and full hostname verification, so this must
// be the name they actually connect to — never a wildcard listen address.
//
// The external tunnel URL wins; the listen address is a fallback for setups
// bound to a concrete address. Both are validated at startup, so a mismatch
// surfaces here rather than as an opaque handshake failure on every agent.
func resolveTunnelHost(cfg *Config) (string, error) {
	if cfg.ExternalTunnelURL != "" {
		u, err := url.Parse(cfg.ExternalTunnelURL)
		if err != nil {
			return "", fmt.Errorf("parse external tunnel URL %q: %w", cfg.ExternalTunnelURL, err)
		}
		host := u.Hostname()
		if host == "" {
			return "", fmt.Errorf(
				"external tunnel URL %q has no host; expected a form like https://tunnel.example.com:8300",
				cfg.ExternalTunnelURL,
			)
		}
		return host, nil
	}

	host, _, err := net.SplitHostPort(cfg.TunnelAddress)
	if err != nil {
		return "", fmt.Errorf("parse tunnel address %q: %w", cfg.TunnelAddress, err)
	}
	if isWildcardHost(host) {
		return "", fmt.Errorf(
			"tunnel address %q binds to a wildcard address, which cannot be used as the tunnel certificate host: "+
				"set --external-tunnel-url to the URL agents dial (e.g. https://tunnel.example.com:8300)",
			cfg.TunnelAddress,
		)
	}
	return host, nil
}

// isWildcardHost covers empty and unspecified addresses (0.0.0.0, ::), which
// would produce a certificate no agent dialing a real address can verify.
func isWildcardHost(host string) bool {
	if host == "" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsUnspecified()
}
