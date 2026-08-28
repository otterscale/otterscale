package chisel

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/otterscale/otterscale/internal/transport"
	"github.com/otterscale/otterscale/internal/transport/tunnel"
)

// BuildTunnelListener issues a server certificate for host and writes the mTLS
// materials to a temporary directory, cleaned up when the listener stops. The
// caller starts the listener via transport.Serve.
func (s *Service) BuildTunnelListener(address, host string) (transport.Listener, error) {
	serverCert, serverKey, err := s.ca.GenerateServerCert(host)
	if err != nil {
		return nil, fmt.Errorf("generate server cert: %w", err)
	}

	certDir, err := os.MkdirTemp("", "otterscale-tls-server-*")
	if err != nil {
		return nil, fmt.Errorf("create cert dir: %w", err)
	}

	caFile := filepath.Join(certDir, "ca.pem")
	certFile := filepath.Join(certDir, "cert.pem")
	keyFile := filepath.Join(certDir, "key.pem")

	const secretFilePerm = 0o600 // owner-only read/write for TLS files
	if err := os.WriteFile(caFile, s.ca.CertPEM(), secretFilePerm); err != nil {
		os.RemoveAll(certDir)
		return nil, fmt.Errorf("write CA cert: %w", err)
	}
	if err := os.WriteFile(certFile, serverCert, secretFilePerm); err != nil {
		os.RemoveAll(certDir)
		return nil, fmt.Errorf("write server cert: %w", err)
	}
	if err := os.WriteFile(keyFile, serverKey, secretFilePerm); err != nil {
		os.RemoveAll(certDir)
		return nil, fmt.Errorf("write server key: %w", err)
	}

	// The CA is per process and never persisted, so a restart invalidates every
	// certificate handed out so far. Say so at startup: the re-registration
	// storm that follows is expected, not a fault to chase.
	slog.Info("tunnel CA initialized",
		"subject", "otterscale-ca",
		"ephemeral", true,
		"detail", "agents re-register after a server restart",
	)

	tunnelSrv, err := tunnel.NewServer(
		tunnel.WithAddress(address),
		tunnel.WithTLSCert(certFile),
		tunnel.WithTLSKey(keyFile),
		tunnel.WithTLSCA(caFile),
		tunnel.WithServer(s.ServerRef()),
	)
	if err != nil {
		os.RemoveAll(certDir)
		return nil, fmt.Errorf("create tunnel server: %w", err)
	}

	return &tunnelListenerWithCleanup{
		Listener: tunnelSrv,
		certDir:  certDir,
	}, nil
}

// tunnelListenerWithCleanup removes the temporary certificate directory on
// Stop.
type tunnelListenerWithCleanup struct {
	transport.Listener
	certDir string
}

func (l *tunnelListenerWithCleanup) Stop(ctx context.Context) error {
	err := l.Listener.Stop(ctx)
	os.RemoveAll(l.certDir)
	return err
}

// BuildHealthListener probes registered endpoints and deregisters disconnected
// clusters.
func (s *Service) BuildHealthListener() transport.Listener {
	return NewHealthCheckListener(s)
}
