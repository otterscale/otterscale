// Package kubernetes provides Kubernetes API access through the
// reverse-tunnel established by the agent. It implements
// core.DiscoveryClient and core.ResourceRepo.
//
// All requests are impersonated: the authenticated user's identity
// (subject + groups) is forwarded to the target cluster's API server
// via Kubernetes impersonation headers, so RBAC is enforced at the
// cluster level rather than at this proxy.
package kubernetes

import (
	"context"
	"net/http"
	"os"
	"sync"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/otterscale/otterscale/internal/core"
)

// clientTimeout is the default HTTP timeout applied to per-request
// rest.Configs. This ensures that Kubernetes API calls that do not
// accept a context.Context (e.g. the discovery client) are still
// bounded and cannot block indefinitely.
const clientTimeout = 30 * time.Second

// clusterTransport holds a cached HTTP transport for a single cluster.
// The transport is shared across users because impersonation is
// handled via HTTP headers (WrapTransport), not at the transport
// level. Only the RoundTripper is cached — per-request clients
// (discovery, dynamic, clientset) are created on the fly from the
// impersonation config.
type clusterTransport struct {
	address string
	rt      http.RoundTripper
}

// Kubernetes is the shared foundation for discoveryClient and
// resourceRepo. It resolves cluster names to tunnel addresses and
// builds impersonated rest.Configs. Transports are cached per-cluster
// and invalidated when the tunnel address changes.
type Kubernetes struct {
	mu         sync.Mutex
	tunnel     core.TunnelProvider
	transports map[string]*clusterTransport // keyed by cluster name
}

// New creates a Kubernetes helper bound to the given TunnelProvider.
func New(tunnel core.TunnelProvider) *Kubernetes {
	return &Kubernetes{
		tunnel:     tunnel,
		transports: make(map[string]*clusterTransport),
	}
}

// impersonationConfig builds a rest.Config that targets the given
// cluster through its tunnel address and impersonates the calling
// user extracted from the request context.
func (k *Kubernetes) impersonationConfig(ctx context.Context, cluster string) (*rest.Config, error) {
	userInfo, ok := core.UserInfoFromContext(ctx)
	if !ok {
		return nil, &core.DomainError{
			Code:    core.ErrorCodeUnauthenticated,
			Message: "user info not found in context",
		}
	}

	// use local kubeconfig for debugging if OTTERSCALE_DEBUG is set.
	if os.Getenv("OTTERSCALE_DEBUG") != "" {
		return clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	}

	address, err := k.tunnel.ResolveAddress(ctx, cluster)
	if err != nil {
		// Cluster is no longer registered; evict stale cached
		// clients and their TCP connections.
		k.evictClients(cluster)
		return nil, err // ResolveAddress already returns *core.ErrClusterNotFound
	}

	rt, err := k.roundTripper(cluster, address)
	if err != nil {
		return nil, err
	}

	cfg := &rest.Config{
		Host: address,
		Impersonate: rest.ImpersonationConfig{
			UserName: userInfo.Subject,
			Groups:   userInfo.Groups,
		},
		Transport: rt,
		Timeout:   clientTimeout,
	}

	return cfg, nil
}

// spdyConfig builds a rest.Config suitable for SPDY connections
// (exec, port-forward). Unlike impersonationConfig, it does NOT
// set a pre-built Transport because SPDY executors and dialers need
// to negotiate their own connection upgrade.
func (k *Kubernetes) spdyConfig(ctx context.Context, cluster string) (*rest.Config, error) {
	userInfo, ok := core.UserInfoFromContext(ctx)
	if !ok {
		return nil, &core.DomainError{
			Code:    core.ErrorCodeUnauthenticated,
			Message: "user info not found in context",
		}
	}

	address, err := k.tunnel.ResolveAddress(ctx, cluster)
	if err != nil {
		// Cluster is no longer registered; evict stale cached
		// clients and their TCP connections.
		k.evictClients(cluster)
		return nil, err // ResolveAddress already returns *core.ErrClusterNotFound
	}

	return &rest.Config{
		Host: address,
		Impersonate: rest.ImpersonationConfig{
			UserName: userInfo.Subject,
			Groups:   userInfo.Groups,
		},
	}, nil
}

// roundTripper returns a cached HTTP transport for the given cluster.
// If the cached transport's address does not match the current tunnel
// address (e.g. after cluster re-registration), the stale entry is
// evicted and a fresh transport is created.
//
// Transports are shared across users because impersonation is handled
// via HTTP headers, not at the transport level. This avoids creating
// new TCP connections on every request.
func (k *Kubernetes) roundTripper(cluster, address string) (http.RoundTripper, error) {
	k.mu.Lock()
	defer k.mu.Unlock()

	if entry, ok := k.transports[cluster]; ok && entry.address == address {
		return entry.rt, nil
	}

	// Address changed or first access — create a fresh transport.
	// Close idle connections on the old transport to avoid leaking
	// TCP connections to a stale tunnel address.
	if old, ok := k.transports[cluster]; ok {
		closeTransport(old.rt)
	}

	cfg := &rest.Config{Host: address}
	rt, err := rest.TransportFor(cfg)
	if err != nil {
		return nil, &core.DomainError{
			Code:    core.ErrorCodeInternal,
			Message: "create HTTP transport",
			Cause:   err,
		}
	}

	k.transports[cluster] = &clusterTransport{
		address: address,
		rt:      rt,
	}
	return rt, nil
}

// evictClients removes the cached transport for the given cluster and
// closes idle TCP connections. This is called when a cluster is no
// longer registered (e.g. after deregistration) to prevent connection
// and memory leaks.
func (k *Kubernetes) evictClients(cluster string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if old, ok := k.transports[cluster]; ok {
		closeTransport(old.rt)
		delete(k.transports, cluster)
	}
}

// closeTransport closes idle connections on the transport if it
// supports the CloseIdleConnections method (e.g. *http.Transport).
func closeTransport(rt http.RoundTripper) {
	type idleCloser interface {
		CloseIdleConnections()
	}
	if ic, ok := rt.(idleCloser); ok {
		ic.CloseIdleConnections()
	}
}
