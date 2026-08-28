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

type ServerOption func(*Server)

// Server manages a chisel reverse-tunnel listener with mTLS authentication and
// automatic user provisioning.
type Server struct {
	serverRef *atomic.Pointer[chserver.Server] // shared with TunnelProvider
	address   string
	tlsCert   string // file path to server certificate
	tlsKey    string // file path to server private key
	tlsCA     string // file path to CA certificate (enables mTLS)
	log       *slog.Logger
}

// WithAddress sets the listen address (e.g. ":8300").
func WithAddress(address string) ServerOption {
	return func(s *Server) { s.address = address }
}

func WithTLSCert(path string) ServerOption {
	return func(s *Server) { s.tlsCert = path }
}

func WithTLSKey(path string) ServerOption {
	return func(s *Server) { s.tlsKey = path }
}

// WithTLSCA turns on mTLS: the server then requires and validates client
// certificates against this CA.
func WithTLSCA(path string) ServerOption {
	return func(s *Server) { s.tlsCA = path }
}

// WithServer injects the shared reference a TunnelProvider owns; init stores
// the initialized server into it so both sides share one running instance.
func WithServer(ref *atomic.Pointer[chserver.Server]) ServerOption {
	return func(s *Server) { s.serverRef = ref }
}

// WithServerLogger defaults to slog.Default with a "component" attribute.
func WithServerLogger(log *slog.Logger) ServerOption {
	return func(s *Server) { s.log = log }
}

// NewServer fully initializes the underlying chisel server, so AddUser works
// through the TunnelProvider even before Start is called.
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

// Start begins accepting connections and blocks until ctx is canceled.
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

func (s *Server) Stop(_ context.Context) error {
	srv := s.serverRef.Load()
	if srv == nil {
		return nil
	}
	s.log.Info("shutting down")
	return srv.Close()
}

// init stores the real chisel server into the shared atomic reference, so any
// TunnelProvider holding it sees the initialized instance.
func (s *Server) init() error {
	cfg := &chserver.Config{
		Reverse: true,
	}

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

	// Chisel allows anonymous connections while no users exist; this disabled
	// sentinel user is what forces authentication.
	if err := ch.AddUser(uuid.NewString(), uuid.NewString(), "127.0.0.1"); err != nil {
		return err
	}

	s.serverRef.Store(ch)
	return nil
}
