// Package agent implements the agent-side runtime that reverse-proxies
// Kubernetes API requests received through a chisel tunnel.
package agent

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/url"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/transport"
	"github.com/otterscale/otterscale/internal/transport/http"
	"github.com/otterscale/otterscale/internal/transport/pipe"
	"github.com/otterscale/otterscale/internal/transport/tunnel"
)

// Config holds the agent's runtime parameters. The Prometheus proxy target is
// absent: NewHandler reads it straight from the application config.
type Config struct {
	Cluster         string
	ServerURL       string
	TunnelServerURL string
}

// Agent binds a local HTTP reverse proxy to a dynamically allocated port and
// exposes it to the control plane through a chisel tunnel.
type Agent struct {
	handler *Handler
	tunnel  core.TunnelConsumer
}

func NewAgent(handler *Handler, tunnel core.TunnelConsumer) *Agent {
	return &Agent{handler: handler, tunnel: tunnel}
}

// Run creates the in-memory pipe listener for the HTTP server, the TCP bridge
// chisel forwards to, and the tunnel client, then blocks until ctx is
// canceled.
func (a *Agent) Run(ctx context.Context, cfg *Config) error {
	warnInsecureServerURL(cfg.ServerURL)

	pl := pipe.NewListener()

	bridge, err := tunnel.NewBridge(ctx, pl)
	if err != nil {
		return fmt.Errorf("failed to create tunnel bridge: %w", err)
	}

	httpSrv, err := http.NewServer(
		ctx,
		http.WithListener(pl),
		// Every request is proxied to kube-apiserver, exec, attach,
		// port-forward, log follow and watch included, whose duration is
		// unbounded.
		http.WithoutRequestTimeouts(),
		http.WithMount(a.handler.Mount),
	)
	if err != nil {
		return fmt.Errorf("failed to create HTTP server: %w", err)
	}

	tunnelClt, err := tunnel.NewClient(
		tunnel.WithServerURL(cfg.ServerURL),
		tunnel.WithTunnelServerURL(cfg.TunnelServerURL),
		tunnel.WithCluster(cfg.Cluster),
		tunnel.WithLocalPort(bridge.Port()),
		tunnel.WithRegister(a.register()),
	)
	if err != nil {
		return fmt.Errorf("failed to create tunnel client: %w", err)
	}
	return transport.Serve(ctx, httpSrv, bridge, tunnelClt)
}

// register adapts the TunnelConsumer to a tunnel.RegisterResult.
func (a *Agent) register() tunnel.RegisterFunc {
	return func(ctx context.Context, serverURL, cluster string) (*tunnel.RegisterResult, error) {
		reg, err := a.tunnel.Register(ctx, serverURL, cluster)
		if err != nil {
			return nil, err
		}

		// The server issues the credential. Deriving it locally forced the
		// agent to reproduce the server's scheme exactly, and keyed it on the
		// agent's hostname — which is not unique across clusters.
		if reg.TunnelUser == "" || reg.TunnelPassword == "" {
			return nil, fmt.Errorf("server returned no tunnel credential; it is older than this agent")
		}

		return &tunnel.RegisterResult{
			Endpoint:  reg.Endpoint,
			Auth:      reg.TunnelUser + ":" + reg.TunnelPassword,
			CACertPEM: reg.CACertificate,
			CertPEM:   reg.Certificate,
			KeyPEM:    reg.PrivateKeyPEM,
		}, nil
	}
}

// ProvideJoinToken fails when no token is configured: the server would
// reject every registration anyway, and failing here says why instead of
// leaving the agent to retry an unauthenticated call forever.
func ProvideJoinToken(conf *config.Config) (core.JoinToken, error) {
	token, err := conf.AgentJoinToken()
	if err != nil {
		return "", err
	}
	if token == "" {
		return "", errors.New(
			"join token is required but not configured; " +
				"set --join-token, --join-token-file or OTTERSCALE_AGENT_JOIN_TOKEN " +
				"to the value of `otterscale join token --cluster <name>`",
		)
	}
	return core.JoinToken(token), nil
}

// warnInsecureServerURL fires when registration would send the join token
// over plaintext HTTP to a remote host, where anything on the path could read
// it and register clusters of its own.
//
// It warns rather than refuses: a service mesh may terminate TLS outside this
// process, and only the operator knows whether that is the case.
func warnInsecureServerURL(serverURL string) {
	u, err := url.Parse(serverURL)
	if err != nil || u.Scheme != "http" || isLoopback(u.Hostname()) {
		return
	}

	slog.Warn("registering over plaintext HTTP: the join token is exposed to anything on the network path",
		"server_url", serverURL,
		"advice", "use https, unless TLS is terminated for this process by a service mesh",
	)
}

// isLoopback means the request never reaches a network.
func isLoopback(host string) bool {
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
