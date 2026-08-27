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
	"github.com/otterscale/otterscale/internal/pki"
	"github.com/otterscale/otterscale/internal/transport"
	"github.com/otterscale/otterscale/internal/transport/http"
	"github.com/otterscale/otterscale/internal/transport/pipe"
	"github.com/otterscale/otterscale/internal/transport/tunnel"
)

// Config holds the runtime parameters for an Agent. The Prometheus
// proxy target is not listed here: it is read straight from the
// application config by NewHandler.
type Config struct {
	Cluster         string
	ServerURL       string
	TunnelServerURL string
}

// Agent binds a local HTTP reverse-proxy to a dynamically allocated
// port and exposes it to the control-plane via a chisel tunnel.
type Agent struct {
	handler *Handler
	tunnel  core.TunnelConsumer
}

// NewAgent returns an Agent wired to the given handler and tunnel
// consumer.
func NewAgent(handler *Handler, tunnel core.TunnelConsumer) *Agent {
	return &Agent{handler: handler, tunnel: tunnel}
}

// Run starts the agent. It creates an in-memory pipe listener for the
// HTTP server, a TCP bridge for chisel to forward to, and a tunnel
// client, then blocks until ctx is canceled.
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
		// Every request here is proxied to kube-apiserver, including
		// exec, attach, port-forward, log follow and watch, whose
		// duration is unbounded.
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

// register wraps the TunnelConsumer so that it returns a
// RegisterResult containing mTLS credentials and derived auth.
func (a *Agent) register() tunnel.RegisterFunc {
	return func(ctx context.Context, serverURL, cluster string) (*tunnel.RegisterResult, error) {
		reg, err := a.tunnel.Register(ctx, serverURL, cluster)
		if err != nil {
			return nil, err
		}

		// Derive the chisel auth string from the signed
		// certificate. This must match the password the server
		// computed when it signed the same certificate.
		auth, err := pki.DeriveAuth(reg.AgentID, reg.Certificate)
		if err != nil {
			return nil, fmt.Errorf("derive auth: %w", err)
		}

		return &tunnel.RegisterResult{
			Endpoint:  reg.Endpoint,
			Auth:      auth,
			CACertPEM: reg.CACertificate,
			CertPEM:   reg.Certificate,
			KeyPEM:    reg.PrivateKeyPEM,
		}, nil
	}
}

// ProvideEnrolmentToken reads the token this agent presents when it
// registers. It fails when none is configured: without one the server
// rejects every registration, and failing here says why instead of
// leaving the agent to retry an unauthenticated call forever.
func ProvideEnrolmentToken(conf *config.Config) (core.EnrolmentToken, error) {
	token, err := conf.AgentEnrolmentToken()
	if err != nil {
		return "", err
	}
	if token == "" {
		return "", errors.New(
			"enrolment token is required but not configured; " +
				"set --enrolment-token, --enrolment-token-file or OTTERSCALE_AGENT_ENROLMENT_TOKEN " +
				"to the value of `otterscale enrolment-token --cluster <name>`",
		)
	}
	return core.EnrolmentToken(token), nil
}

// warnInsecureServerURL warns when registration would send the
// enrolment token over plaintext HTTP to a remote host. Anything on
// that path can read the token and register clusters of its own.
//
// This warns rather than refuses: a service mesh may terminate TLS
// outside this process, which makes plain HTTP on the wire legitimate,
// and only the operator knows whether that is the case.
func warnInsecureServerURL(serverURL string) {
	u, err := url.Parse(serverURL)
	if err != nil || u.Scheme != "http" || isLoopback(u.Hostname()) {
		return
	}

	slog.Warn("registering over plaintext HTTP: the enrolment token is exposed to anything on the network path",
		"server_url", serverURL,
		"advice", "use https, unless TLS is terminated for this process by a service mesh",
	)
}

// isLoopback reports whether host addresses this machine, in which case
// the request never reaches a network.
func isLoopback(host string) bool {
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
