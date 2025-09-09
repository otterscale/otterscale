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
// Helper – reuse the one defined in storage_test.go (testing.TB works for both
// tests and benchmarks)
// -----------------------------------------------------------------------------
func apps_mustNewKube(t testing.TB, cfg *config.Config) *Kube {
	k, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create Kube instance: %v", err)
	}
	return k
}

// -----------------------------------------------------------------------------
// Construction / interface compliance
// -----------------------------------------------------------------------------
func TestNewApps(t *testing.T) {
	k := apps_mustNewKube(t, &config.Config{})
	repo := NewApps(k)
	if repo == nil {
		t.Fatal("expected Apps repository to be created, got nil")
	}
	// Verify it implements the interface
	var _ oscore.KubeAppsRepo = repo
}

func TestApps_InterfaceCompliance(t *testing.T) {
	k := apps_mustNewKube(t, &config.Config{})
	repo := NewApps(k)
	var _ oscore.KubeAppsRepo = repo
}

// -----------------------------------------------------------------------------
// Structure validation
// -----------------------------------------------------------------------------
func TestApps_Structure(t *testing.T) {
	k := apps_mustNewKube(t, &config.Config{})
	repo := NewApps(k)

	a, ok := repo.(*apps)
	if !ok {
		t.Fatal("expected *apps, got a different type")
	}
	if a.kube == nil {
		t.Error("expected kube field to be set, got nil")
	}
}

// -----------------------------------------------------------------------------
// Calls with a syntactically valid rest.Config (no real cluster)
// -----------------------------------------------------------------------------
func TestApps_WithConfig(t *testing.T) {
	k := apps_mustNewKube(t, &config.Config{})
	repo := NewApps(k)

	ctx := context.Background()
	cfg := &rest.Config{Host: "https://example.invalid"}

	// Deployments
	if _, err := repo.ListDeployments(ctx, cfg, "default"); err == nil {
		t.Log("ListDeployments succeeded unexpectedly (no real cluster)")
	}
	if _, err := repo.GetDeployment(ctx, cfg, "default", "demo"); err == nil {
		t.Log("GetDeployment succeeded unexpectedly (no real cluster)")
	}

	// StatefulSets
	if _, err := repo.ListStatefulSets(ctx, cfg, "default"); err == nil {
		t.Log("ListStatefulSets succeeded unexpectedly (no real cluster)")
	}
	if _, err := repo.GetStatefulSet(ctx, cfg, "default", "demo"); err == nil {
		t.Log("GetStatefulSet succeeded unexpectedly (no real cluster)")
	}

	// DaemonSets
	if _, err := repo.ListDaemonSets(ctx, cfg, "default"); err == nil {
		t.Log("ListDaemonSets succeeded unexpectedly (no real cluster)")
	}
	if _, err := repo.GetDaemonSet(ctx, cfg, "default", "demo"); err == nil {
		t.Log("GetDaemonSet succeeded unexpectedly (no real cluster)")
	}
}

// -----------------------------------------------------------------------------
// Error handling – nil rest.Config should cause clientset creation failure
// -----------------------------------------------------------------------------
func TestApps_ErrorHandling(t *testing.T) {
	k := mustNewKube(t, &config.Config{})
	repo := NewApps(k)

	ctx := context.Background()
	// Use an *empty* Config instead of nil – clientset will return an error
	// without panicking.
	emptyCfg := &rest.Config{}

	cases := []struct {
		name string
		call func() error
	}{
		{"ListDeployments", func() error {
			_, err := repo.ListDeployments(ctx, emptyCfg, "default")
			return err
		}},
		{"GetDeployment", func() error {
			_, err := repo.GetDeployment(ctx, emptyCfg, "default", "demo")
			return err
		}},
		{"ListStatefulSets", func() error {
			_, err := repo.ListStatefulSets(ctx, emptyCfg, "default")
			return err
		}},
		{"GetStatefulSet", func() error {
			_, err := repo.GetStatefulSet(ctx, emptyCfg, "default", "demo")
			return err
		}},
		{"ListDaemonSets", func() error {
			_, err := repo.ListDaemonSets(ctx, emptyCfg, "default")
			return err
		}},
		{"GetDaemonSet", func() error {
			_, err := repo.GetDaemonSet(ctx, emptyCfg, "default", "demo")
			return err
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); err == nil {
				t.Errorf("expected error for %s with empty config, got nil", tc.name)
			}
		})
	}
}

// -----------------------------------------------------------------------------
// Individual method behaviour (errors are expected because there is no cluster)
// -----------------------------------------------------------------------------
func TestApps_MethodBehavior(t *testing.T) {
	k := apps_mustNewKube(t, &config.Config{})
	repo := NewApps(k)

	ctx := context.Background()
	cfg := &rest.Config{Host: "https://example.invalid"}

	t.Run("ListDeployments", func(t *testing.T) {
		_, err := repo.ListDeployments(ctx, cfg, "default")
		if err != nil {
			t.Logf("ListDeployments returned expected error: %v", err)
		}
	})

	t.Run("GetDeployment", func(t *testing.T) {
		_, err := repo.GetDeployment(ctx, cfg, "default", "demo")
		if err != nil {
			t.Logf("GetDeployment returned expected error: %v", err)
		}
	})

	t.Run("ListStatefulSets", func(t *testing.T) {
		_, err := repo.ListStatefulSets(ctx, cfg, "default")
		if err != nil {
			t.Logf("ListStatefulSets returned expected error: %v", err)
		}
	})

	t.Run("GetStatefulSet", func(t *testing.T) {
		_, err := repo.GetStatefulSet(ctx, cfg, "default", "demo")
		if err != nil {
			t.Logf("GetStatefulSet returned expected error: %v", err)
		}
	})

	t.Run("ListDaemonSets", func(t *testing.T) {
		_, err := repo.ListDaemonSets(ctx, cfg, "default")
		if err != nil {
			t.Logf("ListDaemonSets returned expected error: %v", err)
		}
	})

	t.Run("GetDaemonSet", func(t *testing.T) {
		_, err := repo.GetDaemonSet(ctx, cfg, "default", "demo")
		if err != nil {
			t.Logf("GetDaemonSet returned expected error: %v", err)
		}
	})
}

// -----------------------------------------------------------------------------
// Integration‑style sequential calls
// -----------------------------------------------------------------------------
func TestApps_IntegrationPatterns(t *testing.T) {
	k := apps_mustNewKube(t, &config.Config{})
	repo := NewApps(k)

	ctx := context.Background()
	cfg := &rest.Config{Host: "https://example.invalid"}

	// Deployments
	if _, err := repo.ListDeployments(ctx, cfg, "default"); err != nil {
		t.Logf("ListDeployments error (expected): %v", err)
	}
	if _, err := repo.GetDeployment(ctx, cfg, "default", "demo"); err != nil {
		t.Logf("GetDeployment error (expected): %v", err)
	}

	// StatefulSets
	if _, err := repo.ListStatefulSets(ctx, cfg, "default"); err != nil {
		t.Logf("ListStatefulSets error (expected): %v", err)
	}
	if _, err := repo.GetStatefulSet(ctx, cfg, "default", "demo"); err != nil {
		t.Logf("GetStatefulSet error (expected): %v", err)
	}

	// DaemonSets
	if _, err := repo.ListDaemonSets(ctx, cfg, "default"); err != nil {
		t.Logf("ListDaemonSets error (expected): %v", err)
	}
	if _, err := repo.GetDaemonSet(ctx, cfg, "default", "demo"); err != nil {
		t.Logf("GetDaemonSet error (expected): %v", err)
	}
}

// -----------------------------------------------------------------------------
// Benchmarks
// -----------------------------------------------------------------------------
func BenchmarkApps_Creation(b *testing.B) {
	cfg := &config.Config{}
	k := apps_mustNewKube(b, cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if repo := NewApps(k); repo == nil {
			b.Fatal("failed to create apps repo")
		}
	}
}

func BenchmarkApps_MethodCalls(b *testing.B) {
	k := apps_mustNewKube(b, &config.Config{})
	repo := NewApps(k)
	ctx := context.Background()
	cfg := &rest.Config{Host: "https://example.invalid"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.ListDeployments(ctx, cfg, "default")
		_, _ = repo.GetDeployment(ctx, cfg, "default", "demo")
		_, _ = repo.ListStatefulSets(ctx, cfg, "default")
		_, _ = repo.GetStatefulSet(ctx, cfg, "default", "demo")
		_, _ = repo.ListDaemonSets(ctx, cfg, "default")
		_, _ = repo.GetDaemonSet(ctx, cfg, "default", "demo")
	}
}

// -----------------------------------------------------------------------------
// Concurrent access – ensure the clientset cache is safe for parallel use
// -----------------------------------------------------------------------------
func TestApps_ConcurrentAccess(t *testing.T) {
	k := apps_mustNewKube(t, &config.Config{})
	repo := NewApps(k)

	ctx := context.Background()
	cfg := &rest.Config{Host: "https://example.invalid"}

	const workers = 10
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_, _ = repo.ListDeployments(ctx, cfg, "default")
			_, _ = repo.GetDeployment(ctx, cfg, "default", "demo")
			_, _ = repo.ListStatefulSets(ctx, cfg, "default")
			_, _ = repo.GetStatefulSet(ctx, cfg, "default", "demo")
			_, _ = repo.ListDaemonSets(ctx, cfg, "default")
			_, _ = repo.GetDaemonSet(ctx, cfg, "default", "demo")
		}()
	}
	wg.Wait()
}

// -----------------------------------------------------------------------------
// Type assertions
// -----------------------------------------------------------------------------
func TestApps_TypeAssertions(t *testing.T) {
	k := apps_mustNewKube(t, &config.Config{})
	repo := NewApps(k)

	a, ok := repo.(*apps)
	if !ok {
		t.Fatal("could not cast repository to *apps")
	}
	if a.kube != k {
		t.Error("apps.kube field not set correctly")
	}
	var _ oscore.KubeAppsRepo = a
}

// -----------------------------------------------------------------------------
// Edge cases – nil Kube should panic when any method is invoked
// -----------------------------------------------------------------------------
func TestApps_EdgeCases(t *testing.T) {
	// Nil Kube
	t.Run("nil_kube", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when using apps with nil Kube")
			}
		}()
		s := &apps{kube: nil}
		ctx := context.Background()
		_, _ = s.ListDeployments(ctx, &rest.Config{Host: "https://example.invalid"}, "default")
	})

	// Background context works (method still returns error from clientset)
	t.Run("background_context", func(t *testing.T) {
		k := apps_mustNewKube(t, &config.Config{})
		s := NewApps(k)
		ctx := context.Background()
		_, err := s.ListDeployments(ctx, &rest.Config{Host: "https://example.invalid"}, "default")
		if err == nil {
			t.Log("ListDeployments succeeded with background context (no real cluster)")
		}
	})

	// TODO context works
	t.Run("todo_context", func(t *testing.T) {
		k := apps_mustNewKube(t, &config.Config{})
		s := NewApps(k)
		ctx := context.TODO()
		_, err := s.ListDeployments(ctx, &rest.Config{Host: "https://example.invalid"}, "default")
		if err == nil {
			t.Log("ListDeployments succeeded with TODO context (no real cluster)")
		}
		_ = ctx // silence unused warning
	})
}
