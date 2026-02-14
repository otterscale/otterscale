package tunnel

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync/atomic"

	"github.com/google/uuid"
	chserver "github.com/jpillora/chisel/server"
)

// ServerOption configures a Server.
type ServerOption func(*Server)

// Server manages a chisel reverse-tunnel listener with mTLS
// certificate authentication and automatic user provisioning.
type Server struct {
	serverRef *atomic.Pointer[chserver.Server] // shared with TunnelProvider
	address   string
	tlsCert   string // file path to server certificate
	tlsKey    string // file path to server private key
	tlsCA     string // file path to CA certificate (enables mTLS)
	log       *slog.Logger
}

// WithAddress configures the listen address (e.g. ":8300").
func WithAddress(address string) ServerOption {
	return func(s *Server) { s.address = address }
}

// WithTLSCert configures the file path to the server TLS certificate.
func WithTLSCert(path string) ServerOption {
	return func(s *Server) { s.tlsCert = path }
}

// WithTLSKey configures the file path to the server TLS private key.
func WithTLSKey(path string) ServerOption {
	return func(s *Server) { s.tlsKey = path }
}

// WithTLSCA configures the file path to the CA certificate used to
// verify client certificates. When set, the server requires and
// validates client certificates (mTLS).
func WithTLSCA(path string) ServerOption {
	return func(s *Server) { s.tlsCA = path }
}

// WithServer injects a shared atomic server reference. The reference
// is typically owned by a TunnelProvider; init will store the fully
// initialized server into it so that both sides share the same
// running instance.
func WithServer(ref *atomic.Pointer[chserver.Server]) ServerOption {
	return func(s *Server) { s.serverRef = ref }
}

// WithServerLogger configures a structured logger. Defaults to
// slog.Default with a "component" attribute.
func WithServerLogger(log *slog.Logger) ServerOption {
	return func(s *Server) { s.log = log }
}

// NewServer creates a tunnel server. The underlying chisel server is
// fully initialized so that AddUser (via TunnelProvider) works
// immediately, even before Start is called.
func NewServer(opts ...ServerOption) (*Server, error) {
	s := &Server{
		serverRef: &atomic.Pointer[chserver.Server]{},
		address:   ":8300",
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.log == nil {
		s.log = slog.Default().With("component", "tunnel-server")
	}
	if err := s.init(); err != nil {
		return nil, fmt.Errorf("tunnel server init: %w", err)
	}
	return s, nil
}

// Start begins accepting connections and blocks until ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	host, port, err := net.SplitHostPort(s.address)
	if err != nil {
		return fmt.Errorf("parse address %q: %w", s.address, err)
	}

	s.log.Info("starting", "address", s.address)

	srv := s.serverRef.Load()
	if err := srv.StartContext(ctx, host, port); err != nil {
		return fmt.Errorf("tunnel server start: %w", err)
	}

	return srv.Wait()
}

// Stop gracefully shuts down the tunnel server.
func (s *Server) Stop(_ context.Context) error {
	srv := s.serverRef.Load()
	if srv == nil {
		return nil
	}
	s.log.Info("shutting down")
	return srv.Close()
}

// init creates the real chisel server and stores it into the shared
// atomic reference so that any TunnelProvider holding the same
// reference sees the fully initialized instance.
func (s *Server) init() error {
	cfg := &chserver.Config{
		Reverse: true,
	}

	// Configure TLS for mTLS when certificate paths are provided.
	if s.tlsCert != "" && s.tlsKey != "" {
		cfg.TLS = chserver.TLSConfig{
			Cert: s.tlsCert,
			Key:  s.tlsKey,
			CA:   s.tlsCA,
		}
	}

	ch, err := chserver.NewServer(cfg)
	if err != nil {
		return err
	}

	// Chisel allows anonymous connections when no users exist.
	// Add a disabled sentinel user to enforce authentication.
	if err := ch.AddUser(uuid.NewString(), uuid.NewString(), "127.0.0.1"); err != nil {
		return err
	}

	// Store the pointer into the shared atomic reference so the
	// TunnelProvider sees the initialized server.
	s.serverRef.Store(ch)
	return nil
}
