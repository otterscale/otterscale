// Package kubernetes provides Kubernetes API access through the reverse tunnel
// established by the agent. It implements core.DiscoveryClient and
// core.ResourceRepo.
//
// All requests are impersonated: the authenticated user's subject and groups
// are forwarded to the target cluster's API server via impersonation headers,
// so RBAC is enforced at the cluster rather than at this proxy.
package kubernetes

import (
	"context"
	"net/http"
	"sync"
	"time"

	"k8s.io/client-go/rest"

	"github.com/otterscale/otterscale/internal/core"
)

// clientTimeout bounds Kubernetes API calls that take no context.Context (the
// discovery client, for one), so they cannot block indefinitely.
const clientTimeout = 30 * time.Second

// clusterTransport caches one cluster's HTTP transport. Only the RoundTripper
// is cached; per-request clients (discovery, dynamic, clientset) are built on
// the fly from the impersonation config.
type clusterTransport struct {
	address string
	rt      http.RoundTripper
}

// Kubernetes is the shared foundation for discoveryClient and resourceRepo. It
// resolves cluster names to tunnel addresses and builds impersonated
// rest.Configs, caching transports per cluster and invalidating them when the
// tunnel address changes.
type Kubernetes struct {
	mu         sync.Mutex
	tunnel     core.TunnelProvider
	transports map[string]*clusterTransport // keyed by cluster name
}

func New(tunnel core.TunnelProvider) *Kubernetes {
	return &Kubernetes{
		tunnel:     tunnel,
		transports: make(map[string]*clusterTransport),
	}
}

// impersonation rejects an empty subject rather than merely checking presence.
// Every authorization decision in this system is made by the target cluster
// against the impersonated identity, and client-go attaches the impersonation
// headers only when at least one of UserName, UID, Groups or Extra is set. A
// UserInfo with no subject and no groups would therefore produce a request
// carrying no identity at all — which the API server answers as the agent's own
// ServiceAccount, an account with far broader rights than any user. Nothing
// constructs such a value today; this makes that a checked invariant rather
// than a property of the current auth middleware.
func impersonation(ctx context.Context) (rest.ImpersonationConfig, error) {
	userInfo, ok := core.UserInfoFromContext(ctx)
	if !ok || userInfo.Subject == "" {
		return rest.ImpersonationConfig{}, &core.DomainError{
			Code:    core.ErrorCodeUnauthenticated,
			Message: "no authenticated subject in context",
		}
	}

	return rest.ImpersonationConfig{
		UserName: userInfo.Subject,
		Groups:   userInfo.Groups,
	}, nil
}

// impersonationConfig targets the cluster through its tunnel address and
// impersonates the caller in ctx.
func (k *Kubernetes) impersonationConfig(ctx context.Context, cluster string) (*rest.Config, error) {
	impersonate, err := impersonation(ctx)
	if err != nil {
		return nil, err
	}

	address, err := k.tunnel.ResolveAddress(ctx, cluster)
	if err != nil {
		// No longer registered; drop stale clients and their TCP connections.
		k.evictClients(cluster)
		return nil, err // ResolveAddress already returns *core.ErrClusterNotFound
	}

	rt, err := k.roundTripper(cluster, address)
	if err != nil {
		return nil, err
	}

	cfg := &rest.Config{
		Host:        address,
		Impersonate: impersonate,
		Transport:   rt,
		Timeout:     clientTimeout,
	}

	return cfg, nil
}

// streamConfig serves exec and port-forward. Unlike impersonationConfig it sets
// no pre-built Transport, because streaming executors and dialers negotiate
// their own connection upgrade.
func (k *Kubernetes) streamConfig(ctx context.Context, cluster string) (*rest.Config, error) {
	impersonate, err := impersonation(ctx)
	if err != nil {
		return nil, err
	}

	address, err := k.tunnel.ResolveAddress(ctx, cluster)
	if err != nil {
		// No longer registered; drop stale clients and their TCP connections.
		k.evictClients(cluster)
		return nil, err // ResolveAddress already returns *core.ErrClusterNotFound
	}

	return &rest.Config{
		Host:        address,
		Impersonate: impersonate,
	}, nil
}

// roundTripper caches one transport per cluster, replacing it when the tunnel
// address changes (as after re-registration). Transports are shared across
// users because impersonation happens in HTTP headers, not at the transport
// level, which avoids a new TCP connection per request.
func (k *Kubernetes) roundTripper(cluster, address string) (http.RoundTripper, error) {
	k.mu.Lock()
	defer k.mu.Unlock()

	if entry, ok := k.transports[cluster]; ok && entry.address == address {
		return entry.rt, nil
	}

	// Close idle connections on the old transport, or they leak to a stale
	// tunnel address.
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

// evictClients runs when a cluster is no longer registered, so its transport
// and idle TCP connections do not leak.
func (k *Kubernetes) evictClients(cluster string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if old, ok := k.transports[cluster]; ok {
		closeTransport(old.rt)
		delete(k.transports, cluster)
	}
}

func closeTransport(rt http.RoundTripper) {
	type idleCloser interface {
		CloseIdleConnections()
	}
	if ic, ok := rt.(idleCloser); ok {
		ic.CloseIdleConnections()
	}
}
