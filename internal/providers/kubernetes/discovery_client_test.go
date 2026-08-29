package kubernetes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core"
)

// TestWatchListCache covers the TTL bookkeeping that keeps a /version
// round-trip off the critical path of every new watch.
func TestWatchListCache(t *testing.T) {
	now := time.Now()
	d := &discoveryClient{
		now:       func() time.Time { return now },
		watchList: make(map[string]watchListSupport),
	}

	if _, ok := d.cachedWatchList("prod"); ok {
		t.Fatal("expected a miss on an empty cache")
	}

	d.cacheWatchList("prod", true)

	supported, ok := d.cachedWatchList("prod")
	if !ok {
		t.Fatal("expected a hit right after caching")
	}
	if !supported {
		t.Error("cached capability = false, want true")
	}

	if _, ok := d.cachedWatchList("staging"); ok {
		t.Error("expected a miss for a cluster that was never cached")
	}

	// A negative answer is cached too, so an older cluster is not
	// re-probed on every watch either.
	d.cacheWatchList("staging", false)
	supported, ok = d.cachedWatchList("staging")
	if !ok {
		t.Fatal("expected a hit for a cached negative answer")
	}
	if supported {
		t.Error("cached capability = true, want false")
	}

	now = now.Add(watchListTTL + time.Second)

	if _, ok := d.cachedWatchList("prod"); ok {
		t.Error("expected a miss once the entry has expired")
	}
}

// fakeTunnel points every cluster at a test server.
type fakeTunnel struct{ address string }

func (fakeTunnel) CACertPEM() []byte               { return nil }
func (fakeTunnel) ListLinks() map[string]core.Link { return nil }
func (fakeTunnel) RegisterLink(context.Context, string, string, string, []byte) (core.TunnelGrant, error) {
	return core.TunnelGrant{}, nil
}

func (f fakeTunnel) ResolveAddress(context.Context, string) (string, error) { return f.address, nil }

// discoveryServer is a stand-in API server that counts how often its
// discovery endpoint is hit and lets a test change what it advertises.
type discoveryServer struct {
	mu        sync.Mutex
	resources []string
	hits      int
}

func (s *discoveryServer) setResources(names ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resources = names
}

func (s *discoveryServer) hitCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.hits
}

func (s *discoveryServer) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/apis/apps/v1" {
			http.NotFound(w, r)
			return
		}

		s.mu.Lock()
		s.hits++
		list := metav1.APIResourceList{GroupVersion: "apps/v1"}
		for _, name := range s.resources {
			list.APIResources = append(list.APIResources, metav1.APIResource{Name: name, Kind: "X"})
		}
		s.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(list)
	})
}

// newDiscoveryTest wires a discoveryClient to a counting test server.
func newDiscoveryTest(t *testing.T) (*discoveryClient, *discoveryServer, context.Context) {
	t.Helper()

	srv := &discoveryServer{}
	srv.setResources("deployments", "deployments/scale")

	httpSrv := httptest.NewServer(srv.handler())
	t.Cleanup(httpSrv.Close)

	d, ok := NewDiscoveryClient(New(fakeTunnel{address: httpSrv.URL})).(*discoveryClient)
	if !ok {
		t.Fatal("NewDiscoveryClient did not return *discoveryClient")
	}

	ctx := core.WithUserInfo(t.Context(), core.UserInfo{
		Subject: "alice",
		Groups:  []string{"system:authenticated"},
	})
	return d, srv, ctx
}

// TestLookupResourceCachesGroupVersion is the regression test for the
// round-trip every resource operation used to pay. Each List, Get,
// Create, Watch or Scale validates its GVR first, so an uncached lookup
// doubled the tunnel traffic of the entire API.
func TestLookupResourceCachesGroupVersion(t *testing.T) {
	d, srv, ctx := newDiscoveryTest(t)

	for range 10 {
		if _, err := d.LookupResource(ctx, "prod", "apps", "v1", "deployments", ""); err != nil {
			t.Fatalf("LookupResource: %v", err)
		}
	}

	// A subresource in the same group/version shares the fetch, which
	// is the point of caching per group/version rather than per lookup.
	if _, err := d.LookupResource(ctx, "prod", "apps", "v1", "deployments", "scale"); err != nil {
		t.Fatalf("LookupResource subresource: %v", err)
	}

	if got := srv.hitCount(); got != 1 {
		t.Errorf("discovery endpoint hit %d times, want 1", got)
	}
}

// TestLookupResourceRefreshesOnMiss covers the freshness rule that lets
// the TTL stay long: a type installed after the list was cached must be
// usable immediately, not once the entry expires.
func TestLookupResourceRefreshesOnMiss(t *testing.T) {
	d, srv, ctx := newDiscoveryTest(t)

	if _, err := d.LookupResource(ctx, "prod", "apps", "v1", "deployments", ""); err != nil {
		t.Fatalf("LookupResource: %v", err)
	}

	// A CRD shows up after the list was cached.
	srv.setResources("deployments", "widgets")

	if _, err := d.LookupResource(ctx, "prod", "apps", "v1", "widgets", ""); err != nil {
		t.Fatalf("newly installed resource not found without waiting for the TTL: %v", err)
	}
	if got := srv.hitCount(); got != 2 {
		t.Errorf("discovery endpoint hit %d times, want 2 (initial + refresh on miss)", got)
	}

	// The refreshed list is cached in turn.
	if _, err := d.LookupResource(ctx, "prod", "apps", "v1", "widgets", ""); err != nil {
		t.Fatalf("LookupResource after refresh: %v", err)
	}
	if got := srv.hitCount(); got != 2 {
		t.Errorf("discovery endpoint hit %d times, want the refreshed list to be cached", got)
	}
}

// TestLookupResourceUnknownResource checks that a genuinely absent
// resource is still reported, after the one refresh.
func TestLookupResourceUnknownResource(t *testing.T) {
	d, srv, ctx := newDiscoveryTest(t)

	if _, err := d.LookupResource(ctx, "prod", "apps", "v1", "deployments", ""); err != nil {
		t.Fatalf("LookupResource: %v", err)
	}

	_, err := d.LookupResource(ctx, "prod", "apps", "v1", "nonexistent", "")
	if err == nil {
		t.Fatal("expected an error for a resource the cluster does not serve")
	}
	code, ok := core.DomainErrorCode(err)
	if !ok || code != core.ErrorCodeInvalidArgument {
		t.Errorf("error = %v (code %v), want ErrorCodeInvalidArgument", err, code)
	}
	if got := srv.hitCount(); got != 2 {
		t.Errorf("discovery endpoint hit %d times, want 2 (initial + one refresh)", got)
	}
}

// TestLookupResourceExpiry covers the TTL bookkeeping.
func TestLookupResourceExpiry(t *testing.T) {
	d, srv, ctx := newDiscoveryTest(t)

	now := time.Now()
	d.now = func() time.Time { return now }

	if _, err := d.LookupResource(ctx, "prod", "apps", "v1", "deployments", ""); err != nil {
		t.Fatalf("LookupResource: %v", err)
	}
	if _, err := d.LookupResource(ctx, "prod", "apps", "v1", "deployments", ""); err != nil {
		t.Fatalf("LookupResource: %v", err)
	}
	if got := srv.hitCount(); got != 1 {
		t.Fatalf("discovery endpoint hit %d times before expiry, want 1", got)
	}

	now = now.Add(resourceListTTL + time.Second)

	if _, err := d.LookupResource(ctx, "prod", "apps", "v1", "deployments", ""); err != nil {
		t.Fatalf("LookupResource after expiry: %v", err)
	}
	if got := srv.hitCount(); got != 2 {
		t.Errorf("discovery endpoint hit %d times, want 2 once the entry expired", got)
	}
}

// TestLookupResourceDeduplicatesConcurrentMisses checks that a burst of
// requests on a cold cache collapses into one fetch, so that a newly
// opened dashboard does not fan out one discovery call per panel.
func TestLookupResourceDeduplicatesConcurrentMisses(t *testing.T) {
	d, srv, ctx := newDiscoveryTest(t)

	const callers = 32

	var wg sync.WaitGroup
	wg.Add(callers)
	for range callers {
		go func() {
			defer wg.Done()
			if _, err := d.LookupResource(ctx, "prod", "apps", "v1", "deployments", ""); err != nil {
				t.Errorf("LookupResource: %v", err)
			}
		}()
	}
	wg.Wait()

	if got := srv.hitCount(); got >= callers {
		t.Errorf("discovery endpoint hit %d times for %d concurrent callers, want them collapsed", got, callers)
	}
}
