// kube/storage_test.go
package kube

import (
	"context"
	"sync"
	"testing"

	"github.com/otterscale/otterscale/internal/config"
	oscore "github.com/otterscale/otterscale/internal/core"
	"k8s.io/client-go/rest"
)

// -----------------------------------------------------------------------------
// Helper – create a minimal Kube instance (ignores errors for the happy path)
// -----------------------------------------------------------------------------
func storage_mustNewKube(t testing.TB, cfg *config.Config) *Kube {
	k, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create Kube instance: %v", err)
	}
	return k
}

// -----------------------------------------------------------------------------
// Basic construction / interface compliance tests
// -----------------------------------------------------------------------------
func TestNewStorage(t *testing.T) {
	k := storage_mustNewKube(t, &config.Config{})
	repo := NewStorage(k)
	if repo == nil {
		t.Fatal("expected Storage repository to be created, got nil")
	}
	// Verify it implements the interface
	var _ oscore.KubeStorageRepo = repo
}

func TestStorage_InterfaceCompliance(t *testing.T) {
	k := storage_mustNewKube(t, &config.Config{})
	repo := NewStorage(k)
	var _ oscore.KubeStorageRepo = repo
}

// -----------------------------------------------------------------------------
// Structure validation
// -----------------------------------------------------------------------------
func TestStorage_Structure(t *testing.T) {
	k := storage_mustNewKube(t, &config.Config{})
	repo := NewStorage(k)

	st, ok := repo.(*storage)
	if !ok {
		t.Fatal("expected *storage, got different type")
	}
	if st.kube == nil {
		t.Error("expected kube field to be set, got nil")
	}
}

// -----------------------------------------------------------------------------
// Tests that run with a valid (but non‑functional) rest.Config
// -----------------------------------------------------------------------------
// The clientset can be created with any host; the actual API calls will fail
// because there is no real Kubernetes API server – that is acceptable for unit
// tests. The purpose is to verify that the code paths are exercised without
// panicking.
func TestStorage_WithConfig(t *testing.T) {
	k := storage_mustNewKube(t, &config.Config{})
	repo := NewStorage(k)

	ctx := context.Background()
	restCfg := &rest.Config{Host: "https://example.invalid"} // any non‑empty host

	// ListStorageClasses
	if _, err := repo.ListStorageClasses(ctx, restCfg); err == nil {
		t.Log("ListStorageClasses succeeded unexpectedly (no real API server)")
	}

	// ListStorageClassesByLabel
	if _, err := repo.ListStorageClassesByLabel(ctx, restCfg, "env=dev"); err == nil {
		t.Log("ListStorageClassesByLabel succeeded unexpectedly (no real API server)")
	}

	// GetStorageClass
	if _, err := repo.GetStorageClass(ctx, restCfg, "standard"); err == nil {
		t.Log("GetStorageClass succeeded unexpectedly (no real API server)")
	}
}

// -----------------------------------------------------------------------------
// Error handling – invalid REST config should cause client creation errors
// -----------------------------------------------------------------------------
func TestStorage_ErrorHandling(t *testing.T) {
	k := storage_mustNewKube(t, &config.Config{})
	repo := NewStorage(k)

	ctx := context.Background()
	invalidCfg := &rest.Config{} // empty Host triggers error in clientset creation

	// ListStorageClasses
	if _, err := repo.ListStorageClasses(ctx, invalidCfg); err == nil {
		t.Error("expected error from ListStorageClasses with empty REST config, got nil")
	}

	// ListStorageClassesByLabel
	if _, err := repo.ListStorageClassesByLabel(ctx, invalidCfg, "foo=bar"); err == nil {
		t.Error("expected error from ListStorageClassesByLabel with empty REST config, got nil")
	}

	// GetStorageClass
	if _, err := repo.GetStorageClass(ctx, invalidCfg, "xyz"); err == nil {
		t.Error("expected error from GetStorageClass with empty REST config, got nil")
	}
}

// -----------------------------------------------------------------------------
// Behaviour of each method when the underlying client works (or fails)
// -----------------------------------------------------------------------------
func TestStorage_MethodBehavior(t *testing.T) {
	k := storage_mustNewKube(t, &config.Config{})
	repo := NewStorage(k)

	ctx := context.Background()
	restCfg := &rest.Config{Host: "https://example.invalid"}

	t.Run("ListStorageClasses", func(t *testing.T) {
		_, err := repo.ListStorageClasses(ctx, restCfg)
		if err != nil {
			t.Logf("ListStorageClasses returned expected error: %v", err)
		}
	})

	t.Run("ListStorageClassesByLabel", func(t *testing.T) {
		_, err := repo.ListStorageClassesByLabel(ctx, restCfg, "app=test")
		if err != nil {
			t.Logf("ListStorageClassesByLabel returned expected error: %v", err)
		}
	})

	t.Run("GetStorageClass", func(t *testing.T) {
		_, err := repo.GetStorageClass(ctx, restCfg, "fast")
		if err != nil {
			t.Logf("GetStorageClass returned expected error: %v", err)
		}
	})
}

// -----------------------------------------------------------------------------
// Simple integration‑style pattern – sequential calls
// -----------------------------------------------------------------------------
func TestStorage_IntegrationPatterns(t *testing.T) {
	k := storage_mustNewKube(t, &config.Config{})
	repo := NewStorage(k)

	ctx := context.Background()
	restCfg := &rest.Config{Host: "https://example.invalid"}

	t.Run("sequential_calls", func(t *testing.T) {
		if _, err := repo.ListStorageClasses(ctx, restCfg); err != nil {
			t.Logf("ListStorageClasses error (expected): %v", err)
		}
		if _, err := repo.ListStorageClassesByLabel(ctx, restCfg, "team=dev"); err != nil {
			t.Logf("ListStorageClassesByLabel error (expected): %v", err)
		}
		if _, err := repo.GetStorageClass(ctx, restCfg, "gold"); err != nil {
			t.Logf("GetStorageClass error (expected): %v", err)
		}
	})
}

// -----------------------------------------------------------------------------
// Benchmarks
// -----------------------------------------------------------------------------
func BenchmarkStorage_Creation(b *testing.B) {
	cfg := &config.Config{}
	k := storage_mustNewKube(b, cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if repo := NewStorage(k); repo == nil {
			b.Fatal("failed to create storage repo")
		}
	}
}

func BenchmarkStorage_MethodCalls(b *testing.B) {
	k := storage_mustNewKube(b, &config.Config{})
	repo := NewStorage(k)
	ctx := context.Background()
	restCfg := &rest.Config{Host: "https://example.invalid"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.ListStorageClasses(ctx, restCfg)
		_, _ = repo.ListStorageClassesByLabel(ctx, restCfg, "env=prod")
		_, _ = repo.GetStorageClass(ctx, restCfg, "standard")
	}
}

// -----------------------------------------------------------------------------
// Concurrency test – ensure no data races on the client‑set cache
// -----------------------------------------------------------------------------
func TestStorage_ConcurrentAccess(t *testing.T) {
	k := storage_mustNewKube(t, &config.Config{})
	repo := NewStorage(k)

	ctx := context.Background()
	restCfg := &rest.Config{Host: "https://example.invalid"}

	const workers = 10
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			if _, err := repo.ListStorageClasses(ctx, restCfg); err != nil {
				// errors are expected; we only care about panics / data races
			}
			if _, err := repo.ListStorageClassesByLabel(ctx, restCfg, "team=ops"); err != nil {
				// ignore
			}
			if _, err := repo.GetStorageClass(ctx, restCfg, "fast"); err != nil {
				// ignore
			}
		}()
	}
	wg.Wait()
}

// -----------------------------------------------------------------------------
// Type assertions
// -----------------------------------------------------------------------------
func TestStorage_TypeAssertions(t *testing.T) {
	k := storage_mustNewKube(t, &config.Config{})
	repo := NewStorage(k)

	st, ok := repo.(*storage)
	if !ok {
		t.Fatal("could not cast repository to *storage")
	}
	if st.kube != k {
		t.Error("storage.kube field not set correctly")
	}
	var _ oscore.KubeStorageRepo = st
}

// -----------------------------------------------------------------------------
// Edge cases – nil Kube pointer should cause panic when methods are invoked
// -----------------------------------------------------------------------------
func TestStorage_EdgeCases(t *testing.T) {
	// Nil Kube
	t.Run("nil_kube", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when using storage with nil Kube")
			}
		}()
		s := &storage{kube: nil}
		ctx := context.Background()
		_, _ = s.ListStorageClasses(ctx, &rest.Config{Host: "https://example.invalid"})
	})

	// Background context works (method should still return an error from the client)
	t.Run("background_context", func(t *testing.T) {
		k := storage_mustNewKube(t, &config.Config{})
		s := NewStorage(k)
		ctx := context.Background()
		_, err := s.ListStorageClasses(ctx, &rest.Config{Host: "https://example.invalid"})
		if err == nil {
			t.Log("ListStorageClasses succeeded with background context (no real cluster)")
		}
	})

	// TODO context works
	t.Run("todo_context", func(t *testing.T) {
		k := storage_mustNewKube(t, &config.Config{})
		s := NewStorage(k)
		ctx := context.TODO()
		_, err := s.ListStorageClasses(ctx, &rest.Config{Host: "https://example.invalid"})
		if err == nil {
			t.Log("ListStorageClasses succeeded with TODO context (no real cluster)")
		}
	})
}
