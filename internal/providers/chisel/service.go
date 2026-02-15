// Package chisel implements core.TunnelProvider using jpillora/chisel.
//
// Each registered cluster is assigned a unique loopback address in
// the 127.x.x.x range so that chisel can route reverse-tunnel traffic
// to the correct agent without port conflicts.
package chisel

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	chserver "github.com/jpillora/chisel/server"

	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/pki"
)

// tunnelPort is the fixed port shared by all cluster tunnels.
// Each cluster is differentiated by its loopback host, not its port.
const tunnelPort = 16598

// maxHosts is the total number of unique loopback addresses available
// in the range 127.1.1.1 – 127.254.254.254 (octets 0 and 255 are
// avoided).
const maxHosts = 254 * 254 * 254

// Service manages the mapping between cluster names and unique
// loopback addresses, and provisions chisel users for each agent.
// It implements core.TunnelProvider and transport.TunnelService.
type Service struct {
	server atomic.Pointer[chserver.Server]
	ca     *pki.CA
	log    *slog.Logger
	addrs  *addressAllocator

	mu       sync.RWMutex
	clusters map[string]core.Cluster // cluster name -> tunnel state
}

// NewService returns a new Service backed by chisel. The CA is
// required for signing agent CSRs and must be provided at
// construction time (dependency injection).
// The underlying chisel server is lazily initialized by the tunnel
// transport layer; see tunnel.NewServer.
func NewService(ca *pki.CA) *Service {
	return &Service{
		ca:       ca,
		log:      slog.Default().With("component", "tunnel-provider"),
		addrs:    newAddressAllocator(),
		clusters: make(map[string]core.Cluster),
	}
}

var _ core.TunnelProvider = (*Service)(nil)

// ServerRef returns a pointer to the atomic chisel server reference.
// The tunnel transport stores the fully initialized server into this
// reference at startup so that both sides share the same instance.
// This method is intentionally NOT part of core.TunnelProvider to keep
// the domain layer free of chisel dependencies.
func (s *Service) ServerRef() *atomic.Pointer[chserver.Server] {
	return &s.server
}

// CA returns the CA used to sign agent CSRs and generate server
// certificates. This is provided at construction time via DI.
func (s *Service) CA() *pki.CA {
	return s.ca
}

// CACertPEM returns the PEM-encoded CA certificate so that agents
// can verify the tunnel server's identity via mTLS.
func (s *Service) CACertPEM() []byte {
	return s.ca.CertPEM()
}

// ListClusters returns the names of all currently registered clusters.
func (s *Service) ListClusters() map[string]core.Cluster {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return maps.Clone(s.clusters)
}

// RegisterCluster validates and signs the agent's CSR, associates a
// cluster with a unique loopback host, creates a chisel user with a
// password derived from the signed certificate, and returns the
// tunnel endpoint and the PEM-encoded signed certificate.
//
// If the cluster was previously registered, the old host allocation
// is released first so that re-registration always moves the cluster
// to a fresh address.
func (s *Service) RegisterCluster(ctx context.Context, cluster, agentID, agentVersion string, csrPEM []byte) (string, []byte, error) {
	// Sign the agent's CSR with the internal CA.
	certPEM, err := s.ca.SignCSR(csrPEM)
	if err != nil {
		return "", nil, fmt.Errorf("sign CSR: %w", err)
	}

	// Derive the chisel password from the signed certificate so
	// that both server and agent can compute it independently.
	auth, err := pki.DeriveAuth(agentID, certPEM)
	if err != nil {
		return "", nil, fmt.Errorf("derive auth: %w", err)
	}
	_, pass, ok := parseAuth(auth)
	if !ok {
		return "", nil, fmt.Errorf("invalid auth format: expected user:pass, got %q", auth)
	}

	srv := s.server.Load()
	if srv == nil {
		return "", nil, &core.ErrNotReady{Subsystem: "chisel server"}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Release the previous host and user for this cluster, if any,
	// so that stale credentials do not accumulate in chisel.
	if prev, ok := s.clusters[cluster]; ok {
		srv.DeleteUser(prev.User)
		s.addrs.release(prev.Host)
		delete(s.clusters, cluster)
	}

	host, err := s.addrs.allocate(cluster)
	if err != nil {
		return "", nil, err
	}

	// Restrict the user to reverse-tunneling only the allocated
	// host:port combination. The regex anchors prevent the agent
	// from binding arbitrary endpoints.
	allowed := fmt.Sprintf("^R:%s:%d(:.*)?$", regexp.QuoteMeta(host), tunnelPort)
	if err := srv.AddUser(agentID, pass, allowed); err != nil {
		s.addrs.release(host)
		return "", nil, err
	}

	s.clusters[cluster] = core.Cluster{
		Host:         host,
		User:         agentID,
		AgentVersion: agentVersion,
	}

	return fmt.Sprintf("%s:%d", host, tunnelPort), certPEM, nil
}

// DeregisterCluster removes a cluster's tunnel allocation, deleting
// the chisel user and releasing the loopback host. It is a no-op if
// the cluster is not currently registered.
func (s *Service) DeregisterCluster(cluster string) {
	srv := s.server.Load()
	if srv == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.clusters[cluster]
	if !ok {
		return
	}
	srv.DeleteUser(entry.User)
	s.addrs.release(entry.Host)
	delete(s.clusters, cluster)
}

// ResolveAddress returns the HTTP base URL for the given cluster's
// tunnel endpoint. Returns an error if the cluster is not registered.
func (s *Service) ResolveAddress(ctx context.Context, cluster string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.clusters[cluster]
	if !ok {
		return "", &core.ErrClusterNotFound{Cluster: cluster}
	}

	return fmt.Sprintf("http://%s:%d", entry.Host, tunnelPort), nil
}

// parseAuth splits a "user:pass" string into its components.
func parseAuth(auth string) (user, pass string, ok bool) {
	return strings.Cut(auth, ":")
}
