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

// Config holds the runtime parameters for a Server.
type Config struct {
	Address           string
	AllowedOrigins    []string
	TunnelAddress     string
	ExternalTunnelURL string
	KeycloakRealmURL  string
	KeycloakClientID  string
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

// resolveTunnelHost determines the host embedded as the SAN of the
// tunnel server's TLS certificate. Agents dial the tunnel over mTLS
// with the CA pinned and full hostname verification enabled, so this
// host must be the name agents actually connect to — never a wildcard
// listen address.
//
// The external tunnel URL is preferred; the local listen address is
// only a fallback for setups that bind to a concrete address. Both
// paths are validated at startup so a mismatch surfaces here instead
// of as an opaque TLS handshake failure on every agent.
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

// isWildcardHost reports whether host is empty or an unspecified IP
// address (0.0.0.0, ::). Such a host would produce a certificate no
// agent can verify, since agents dial a real hostname or address.
func isWildcardHost(host string) bool {
	if host == "" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsUnspecified()
}
