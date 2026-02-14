// Package agent implements the agent-side runtime that reverse-proxies
// Kubernetes API requests received through a chisel tunnel.
package agent

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"k8s.io/client-go/rest"

	"github.com/otterscale/otterscale/internal/bootstrap"
	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/pki"
	"github.com/otterscale/otterscale/internal/transport"
	"github.com/otterscale/otterscale/internal/transport/http"
	"github.com/otterscale/otterscale/internal/transport/pipe"
	"github.com/otterscale/otterscale/internal/transport/tunnel"
)

// Config holds the runtime parameters for an Agent.
type Config struct {
	Cluster         string
	ServerURL       string
	TunnelServerURL string
	Bootstrap       bool
}

// SelfUpdater abstracts the self-update mechanism so it can be
// injected via DI and mocked in tests.
type SelfUpdater interface {
	Patch(ctx context.Context, version string) error
}

// Agent binds a local HTTP reverse-proxy to a dynamically allocated
// port and exposes it to the control-plane via a chisel tunnel.
type Agent struct {
	cfg          *rest.Config
	handler      *Handler
	tunnel       core.TunnelConsumer
	version      core.Version
	bootstrapper *bootstrap.Bootstrapper
	updater      SelfUpdater
}

// NewAgent returns an Agent wired to the given handler, tunnel
// consumer, bootstrapper, and self-updater. version is injected via
// DI and used for version-mismatch detection during registration.
func NewAgent(cfg *rest.Config, handler *Handler, tunnel core.TunnelConsumer, version core.Version, bootstrapper *bootstrap.Bootstrapper, updater SelfUpdater) *Agent {
	return &Agent{cfg: cfg, handler: handler, tunnel: tunnel, version: version, bootstrapper: bootstrapper, updater: updater}
}

// Run starts the agent. When bootstrap is enabled, it first applies
// embedded infrastructure manifests (FluxCD, Module CRD) to the local
// cluster. It then creates an in-memory pipe listener for the HTTP
// server, a TCP bridge for chisel to forward to, and a tunnel client,
// then blocks until ctx is cancelled.
func (a *Agent) Run(ctx context.Context, cfg Config) error {
	if cfg.Bootstrap {
		if err := a.bootstrapper.Run(ctx); err != nil {
			return fmt.Errorf("bootstrap: %w", err)
		}
	}

	pl := pipe.NewListener()

	bridge, err := tunnel.NewBridge(pl)
	if err != nil {
		return fmt.Errorf("failed to create tunnel bridge: %w", err)
	}

	httpSrv, err := http.NewServer(
		http.WithListener(pl),
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
		tunnel.WithKeepAlive(30*time.Second),
		tunnel.WithMaxRetryCount(6),
		tunnel.WithMaxRetryInterval(10*time.Second),
		tunnel.WithRegister(a.register()),
	)
	if err != nil {
		return fmt.Errorf("failed to create tunnel client: %w", err)
	}
	return transport.Serve(ctx, httpSrv, bridge, tunnelClt)
}

// register wraps the TunnelConsumer so that it returns a
// RegisterResult containing mTLS credentials and derived auth.
// After a successful registration it checks whether the server
// version diverges from the agent version and, if so, triggers a
// self-update by patching its own Deployment image.
func (a *Agent) register() tunnel.RegisterFunc {
	return func(ctx context.Context, serverURL, cluster string) (*tunnel.RegisterResult, error) {
		reg, err := a.tunnel.Register(ctx, serverURL, cluster)
		if err != nil {
			return nil, err
		}

		// Check version and trigger self-update if needed.
		a.checkVersion(ctx, reg)

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

// checkVersion compares the agent and server versions. When they
// differ and a self-updater is configured, the agent patches its own
// Deployment image to trigger a rolling update. Errors are logged but
// do not prevent the tunnel from connecting — the agent continues to
// serve with the current version.
func (a *Agent) checkVersion(ctx context.Context, reg core.Registration) {
	log := slog.Default().With("component", "version-check")

	if reg.ServerVersion == "" {
		log.Debug("server did not report a version, skipping check")
		return
	}

	agentVersion := string(a.version)

	if reg.ServerVersion == agentVersion {
		log.Info("version match", "version", agentVersion)
		return
	}

	log.Warn("version mismatch detected",
		"agent_version", agentVersion,
		"server_version", reg.ServerVersion,
	)

	if err := a.updater.Patch(ctx, reg.ServerVersion); err != nil {
		log.Error("self-update failed", "error", err)
	}
}
