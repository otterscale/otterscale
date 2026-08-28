// Package chisel implements core.TunnelProvider using jpillora/chisel.
//
// Each registered cluster is assigned a unique loopback address in the
// 127.x.x.x range, so chisel can route reverse-tunnel traffic to the correct
// agent without port conflicts.
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

// tunnelPort is shared by all cluster tunnels; clusters are differentiated by
// loopback host, not by port.
const tunnelPort = 16598

// maxHosts counts the unique loopback addresses in 127.1.1.1 – 127.254.254.254
// (octets 0 and 255 are avoided).
const maxHosts = 254 * 254 * 254

// Service maps cluster names to unique loopback addresses and provisions chisel
// users for each agent. It implements core.TunnelProvider and
// transport.TunnelService.
type Service struct {
	server atomic.Pointer[chserver.Server]
	ca     *pki.CA
	log    *slog.Logger
	addrs  *addressAllocator

	mu    sync.RWMutex
	links map[string]core.Link // cluster name -> tunnel state
}

// NewService takes the CA that signs agent CSRs. The underlying chisel server
// is lazily initialized by the tunnel transport layer; see tunnel.NewServer.
func NewService(ca *pki.CA) *Service {
	return &Service{
		ca:    ca,
		log:   slog.Default().With("component", "tunnel-provider"),
		addrs: newAddressAllocator(),
		links: make(map[string]core.Link),
	}
}

var _ core.TunnelProvider = (*Service)(nil)

// ServerRef exposes the atomic chisel server reference that the tunnel
// transport fills in at startup, so both sides share one instance. It is
// deliberately not part of core.TunnelProvider, to keep the domain layer free
// of chisel.
func (s *Service) ServerRef() *atomic.Pointer[chserver.Server] {
	return &s.server
}

func (s *Service) CA() *pki.CA {
	return s.ca
}

func (s *Service) CACertPEM() []byte {
	return s.ca.CertPEM()
}

func (s *Service) ListLinks() map[string]core.Link {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return maps.Clone(s.links)
}

// RegisterLink signs the agent's CSR, associates the cluster with a unique
// loopback host, creates a chisel user with a fresh password, and returns the
// grant the agent needs to connect.
//
// A previous registration's chisel user is retired only after the replacement
// is live, and the cluster keeps the address it already had.
func (s *Service) RegisterLink(_ context.Context, cluster, agentID, agentVersion string, csrPEM []byte) (core.TunnelGrant, error) {
	// The CSR is request data, so a rejection is the caller's problem; a CA that
	// cannot sign is ours, and the two must not report the same code.
	certPEM, err := s.ca.SignCSR(csrPEM)
	if err != nil {
		code := core.ErrorCodeInternal
		if errors.Is(err, pki.ErrInvalidCSR) {
			code = core.ErrorCodeInvalidArgument
		}
		return core.TunnelGrant{}, &core.DomainError{Code: code, Message: "sign CSR", Cause: err}
	}

	// The tunnel user is the cluster, never the agent. Agents identify
	// themselves by hostname, and two clusters running the same deployment
	// report the same one — chisel's user index is keyed by name, so an
	// agent-derived name lets one cluster's registration overwrite another's
	// credentials and strand a working tunnel. Keying by cluster also keeps the
	// index aligned with this registry, so deregistering a cluster cannot delete
	// a user that now belongs to a different one.
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

	// Nothing is torn down until the replacement is in place. An earlier version
	// released the previous host and user first, so a failure below left the
	// cluster deregistered — a re-registration that could not complete took the
	// agent already serving that cluster offline with it.
	//
	// A cluster also keeps its address: re-registration replaces the agent
	// behind a cluster, not the cluster's identity on the tunnel, and churning
	// the address would invalidate the per-cluster transport cached against it.
	prev, registered := s.links[cluster]

	host := prev.Host
	if !registered {
		host, err = s.addrs.allocate(cluster)
		if err != nil {
			return core.TunnelGrant{}, err
		}
	}

	// Anchoring the regex is what stops the agent from reverse-tunneling
	// anything but the allocated host:port.
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

	// The new credential is live. AddUser replaces an entry of the same name, so
	// only a predecessor registered under a different one needs retiring.
	if registered && prev.User != user {
		srv.DeleteUser(prev.User)
	}

	s.links[cluster] = core.Link{
		Host:         host,
		User:         user,
		AgentVersion: agentVersion,
	}

	// agent_id no longer decides anything, but it is still the only record of
	// which agent claimed the cluster.
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

// DeregisterCluster deletes the chisel user and releases the loopback host. It
// is a no-op for a cluster that is not registered.
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

func (s *Service) ResolveAddress(_ context.Context, cluster string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.links[cluster]
	if !ok {
		return "", &core.ErrClusterNotFound{Cluster: cluster}
	}

	return fmt.Sprintf("http://%s:%d", entry.Host, tunnelPort), nil
}
