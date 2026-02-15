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

// BuildTunnelListener generates a server TLS certificate for the
// given host, writes the mTLS materials to a temporary directory,
// and returns a fully configured tunnel transport.Listener. The
// caller is responsible for starting the listener via transport.Serve.
// The temporary certificate files are cleaned up when the listener
// stops.
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

	if err := os.WriteFile(caFile, s.ca.CertPEM(), 0o600); err != nil {
		os.RemoveAll(certDir)
		return nil, fmt.Errorf("write CA cert: %w", err)
	}
	if err := os.WriteFile(certFile, serverCert, 0o600); err != nil {
		os.RemoveAll(certDir)
		return nil, fmt.Errorf("write server cert: %w", err)
	}
	if err := os.WriteFile(keyFile, serverKey, 0o600); err != nil {
		os.RemoveAll(certDir)
		return nil, fmt.Errorf("write server key: %w", err)
	}

	slog.Info("tunnel CA initialized", "subject", "otterscale-ca")

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

// tunnelListenerWithCleanup wraps a transport.Listener and removes
// the temporary TLS certificate directory when stopped.
type tunnelListenerWithCleanup struct {
	transport.Listener
	certDir string
}

func (l *tunnelListenerWithCleanup) Stop(ctx context.Context) error {
	err := l.Listener.Stop(ctx)
	os.RemoveAll(l.certDir)
	return err
}

// BuildHealthListener returns a transport.Listener that periodically
// health-checks registered tunnel endpoints and deregisters
// disconnected clusters.
func (s *Service) BuildHealthListener() transport.Listener {
	return NewHealthCheckListener(s)
}
