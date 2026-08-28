// Package chisel implements core.TunnelProvider using jpillora/chisel.
//
// Each registered cluster is assigned a unique loopback address in
// the 127.x.x.x range so that chisel can route reverse-tunnel traffic
// to the correct agent without port conflicts.
package chisel

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"regexp"
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

	mu    sync.RWMutex
	links map[string]core.Link // cluster name -> tunnel state
}

// NewService returns a new Service backed by chisel. The CA is
// required for signing agent CSRs and must be provided at
// construction time (dependency injection).
// The underlying chisel server is lazily initialized by the tunnel
// transport layer; see tunnel.NewServer.
func NewService(ca *pki.CA) *Service {
	return &Service{
		ca:    ca,
		log:   slog.Default().With("component", "tunnel-provider"),
		addrs: newAddressAllocator(),
		links: make(map[string]core.Link),
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

// ListLinks returns the names of all currently registered links.
func (s *Service) ListLinks() map[string]core.Link {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return maps.Clone(s.links)
}

// RegisterLink validates and signs the agent's CSR, associates a
// cluster with a unique loopback host, creates a chisel user with a
// freshly generated password, and returns the grant the agent needs to
// connect: endpoint, signed certificate, and tunnel credential.
//
// If the cluster was previously registered, the old host allocation
// and chisel user are released first so that stale credentials do not
// accumulate. The replacement address is normally the same one: the
// allocator derives its starting probe from a hash of the cluster
// name, so a released host is handed straight back to it.
func (s *Service) RegisterLink(_ context.Context, cluster, agentID, agentVersion string, csrPEM []byte) (core.TunnelGrant, error) {
	// Sign the agent's CSR with the internal CA. The CSR is request
	// data, so a rejection is the caller's problem; a CA that cannot
	// sign is ours, and the two must not report the same code.
	certPEM, err := s.ca.SignCSR(csrPEM)
	if err != nil {
		code := core.ErrorCodeInternal
		if errors.Is(err, pki.ErrInvalidCSR) {
			code = core.ErrorCodeInvalidArgument
		}
		return core.TunnelGrant{}, &core.DomainError{Code: code, Message: "sign CSR", Cause: err}
	}

	// The tunnel user is the cluster, never the agent. Agents identify
	// themselves by hostname, and two clusters running the same
	// deployment report the same one — chisel's user index is keyed by
	// name, so an agent-derived name lets one cluster's registration
	// overwrite another's credentials and strand a working tunnel.
	// Keying by cluster also keeps the index aligned with this
	// registry, so deregistering a cluster cannot delete a user that
	// now belongs to a different one.
	user := cluster
	pass, err := pki.NewTunnelPassword()
	if err != nil {
		return core.TunnelGrant{}, err
	}

	srv := s.server.Load()
	if srv == nil {
		return core.TunnelGrant{}, &core.ErrNotReady{Subsystem: "chisel server"}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Nothing is torn down until the replacement is in place. An
	// earlier version released the previous host and user first, so a
	// failure below left the cluster deregistered — a re-registration
	// that could not complete took the agent already serving that
	// cluster offline with it.
	//
	// A cluster keeps the address it was given. Re-registration
	// replaces the agent behind a cluster, not the cluster's identity
	// on the tunnel, and churning the address would invalidate the
	// per-cluster transport cached against it for no reason.
	prev, registered := s.links[cluster]

	host := prev.Host
	if !registered {
		host, err = s.addrs.allocate(cluster)
		if err != nil {
			return core.TunnelGrant{}, err
		}
	}

	// Restrict the user to reverse-tunneling only the allocated
	// host:port combination. The regex anchors prevent the agent
	// from binding arbitrary endpoints.
	allowed := fmt.Sprintf("^R:%s:%d(:.*)?$", regexp.QuoteMeta(host), tunnelPort)
	if err := srv.AddUser(user, pass, allowed); err != nil {
		if !registered {
			s.addrs.release(host)
		}
		return core.TunnelGrant{}, &core.DomainError{
			Code:    core.ErrorCodeInternal,
			Message: "provision tunnel user",
			Cause:   err,
		}
	}

	// The new credential is live. AddUser replaces an entry of the same
	// name, so only a predecessor registered under a different one has
	// to be retired explicitly.
	if registered && prev.User != user {
		srv.DeleteUser(prev.User)
	}

	s.links[cluster] = core.Link{
		Host:         host,
		User:         user,
		AgentVersion: agentVersion,
	}

	// agent_id no longer decides anything, but it is still the only
	// record of which agent claimed the cluster.
	s.log.Info("registered link",
		"cluster", cluster,
		"agent_id", agentID,
		"agent_version", agentVersion,
		"host", host,
	)

	return core.TunnelGrant{
		Endpoint:    fmt.Sprintf("%s:%d", host, tunnelPort),
		Certificate: certPEM,
		User:        user,
		Password:    pass,
	}, nil
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

	entry, ok := s.links[cluster]
	if !ok {
		return
	}
	srv.DeleteUser(entry.User)
	s.addrs.release(entry.Host)
	delete(s.links, cluster)
}

// ResolveAddress returns the HTTP base URL for the given cluster's
// tunnel endpoint. Returns an error if the cluster is not registered.
func (s *Service) ResolveAddress(_ context.Context, cluster string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.links[cluster]
	if !ok {
		return "", &core.ErrClusterNotFound{Cluster: cluster}
	}

	return fmt.Sprintf("http://%s:%d", entry.Host, tunnelPort), nil
}
