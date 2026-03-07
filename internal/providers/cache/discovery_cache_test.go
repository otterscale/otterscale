package cache

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/kube-openapi/pkg/validation/spec"

	"github.com/otterscale/otterscale/internal/core"
)

// stubDiscovery implements core.DiscoveryClient with configurable
// columns results and call counting.
type stubDiscovery struct {
	columns    []core.ColumnDefinition
	columnsErr error
	callCount  atomic.Int64
}

func (s *stubDiscovery) LookupResource(context.Context, string, string, string, string) (schema.GroupVersionResource, error) {
	return schema.GroupVersionResource{}, nil
}

func (s *stubDiscovery) ServerResources(context.Context, string) ([]*metav1.APIResourceList, error) {
	return nil, nil
}

func (s *stubDiscovery) ResolveSchema(context.Context, string, string, string, string) (*spec.Schema, error) {
	return &spec.Schema{}, nil
}

func (s *stubDiscovery) ServerVersion(context.Context, string) (*version.Info, error) {
	return &version.Info{}, nil
}

func (s *stubDiscovery) SupportsWatchList(context.Context, string) (bool, error) {
	return false, nil
}

func (s *stubDiscovery) Columns(_ context.Context, _ string, _ schema.GroupVersionResource, _ string) ([]core.ColumnDefinition, error) {
	s.callCount.Add(1)
	return s.columns, s.columnsErr
}

func TestDiscoveryCache_Columns_CachesResult(t *testing.T) {
	cols := []core.ColumnDefinition{
		{Name: "Name", Type: "string", Priority: 0},
		{Name: "Ready", Type: "string", Priority: 0},
	}
	stub := &stubDiscovery{columns: cols}

	now := time.Now()
	cache := NewDiscoveryCache(stub, 10*time.Minute, WithClock(func() time.Time { return now }))
	ctx := context.Background()
	gvr := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	// First call: hits upstream.
	result, err := cache.Columns(ctx, "prod", gvr, "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(result))
	}
	if stub.callCount.Load() != 1 {
		t.Fatalf("expected 1 upstream call, got %d", stub.callCount.Load())
	}

	// Second call: served from cache.
	result, err = cache.Columns(ctx, "prod", gvr, "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(result))
	}
	if stub.callCount.Load() != 1 {
		t.Fatalf("expected 1 upstream call (cached), got %d", stub.callCount.Load())
	}
}

func TestDiscoveryCache_Columns_ExpiredEntry(t *testing.T) {
	cols := []core.ColumnDefinition{
		{Name: "Name", Type: "string"},
	}
	stub := &stubDiscovery{columns: cols}

	now := time.Now()
	cache := NewDiscoveryCache(stub, 1*time.Minute, WithClock(func() time.Time { return now }))
	ctx := context.Background()
	gvr := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}

	// Populate cache.
	_, _ = cache.Columns(ctx, "prod", gvr, "")
	if stub.callCount.Load() != 1 {
		t.Fatalf("expected 1 call, got %d", stub.callCount.Load())
	}

	// Advance time past TTL.
	now = now.Add(2 * time.Minute)

	_, _ = cache.Columns(ctx, "prod", gvr, "")
	if stub.callCount.Load() != 2 {
		t.Fatalf("expected 2 calls after TTL expiry, got %d", stub.callCount.Load())
	}
}

func TestDiscoveryCache_Columns_DifferentKeys(t *testing.T) {
	cols := []core.ColumnDefinition{{Name: "Name"}}
	stub := &stubDiscovery{columns: cols}

	cache := NewDiscoveryCache(stub, 10*time.Minute)
	ctx := context.Background()

	gvr1 := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	gvr2 := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}

	_, _ = cache.Columns(ctx, "prod", gvr1, "")
	_, _ = cache.Columns(ctx, "prod", gvr2, "")

	if stub.callCount.Load() != 2 {
		t.Fatalf("expected 2 calls for different GVRs, got %d", stub.callCount.Load())
	}
}

func TestDiscoveryCache_EvictExpiredColumns(t *testing.T) {
	cols := []core.ColumnDefinition{{Name: "Col"}}
	stub := &stubDiscovery{columns: cols}

	now := time.Now()
	cache := NewDiscoveryCache(stub, 1*time.Minute, WithClock(func() time.Time { return now }))
	ctx := context.Background()
	gvr := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}

	_, _ = cache.Columns(ctx, "prod", gvr, "")

	cache.mu.RLock()
	if len(cache.columnsCache) != 1 {
		cache.mu.RUnlock()
		t.Fatalf("expected 1 cache entry, got %d", len(cache.columnsCache))
	}
	cache.mu.RUnlock()

	// Advance time past TTL and run eviction.
	now = now.Add(2 * time.Minute)
	cache.mu.Lock()
	cache.evictExpiredColumns()
	count := len(cache.columnsCache)
	cache.mu.Unlock()

	if count != 0 {
		t.Fatalf("expected 0 cache entries after eviction, got %d", count)
	}
}
