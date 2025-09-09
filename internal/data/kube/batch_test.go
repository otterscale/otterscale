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
// Helper – works for both normal tests (testing.T) and benchmarks (testing.B)
// -----------------------------------------------------------------------------
func batch_mustNewKube(t testing.TB, cfg *config.Config) *Kube {
	k, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create Kube instance: %v", err)
	}
	return k
}

// -----------------------------------------------------------------------------
// Construction / interface compliance
// -----------------------------------------------------------------------------
func TestNewBatch(t *testing.T) {
	k := batch_mustNewKube(t, &config.Config{})
	repo := NewBatch(k)
	if repo == nil {
		t.Fatal("expected Batch repository to be created, got nil")
	}
	// Verify it implements the interface
	var _ oscore.KubeBatchRepo = repo
}

func TestBatch_InterfaceCompliance(t *testing.T) {
	k := batch_mustNewKube(t, &config.Config{})
	repo := NewBatch(k)
	var _ oscore.KubeBatchRepo = repo
}

// -----------------------------------------------------------------------------
// Structure validation
// -----------------------------------------------------------------------------
func TestBatch_Structure(t *testing.T) {
	k := batch_mustNewKube(t, &config.Config{})
	repo := NewBatch(k)

	b, ok := repo.(*batch)
	if !ok {
		t.Fatal("expected *batch, got a different type")
	}
	if b.kube == nil {
		t.Error("expected kube field to be set, got nil")
	}
}

// -----------------------------------------------------------------------------
// Calls with a syntactically valid rest.Config (no real cluster)
// -----------------------------------------------------------------------------
func TestBatch_WithConfig(t *testing.T) {
	k := batch_mustNewKube(t, &config.Config{})
	repo := NewBatch(k)

	ctx := context.Background()
	restCfg := &rest.Config{Host: "https://example.invalid"} // any non‑empty host

	// List all jobs
	if _, err := repo.ListJobs(ctx, restCfg, "default"); err == nil {
		t.Log("ListJobs succeeded unexpectedly (no real cluster)")
	}
	// List by label
	if _, err := repo.ListJobsByLabel(ctx, restCfg, "default", "app=demo"); err == nil {
		t.Log("ListJobsByLabel succeeded unexpectedly (no real cluster)")
	}
	// Create a job (nil spec)
	if _, err := repo.CreateJob(ctx, restCfg, "default", "demo", nil, nil, nil); err == nil {
		t.Log("CreateJob succeeded unexpectedly (no real cluster)")
	}
	// Delete a job
	if err := repo.DeleteJob(ctx, restCfg, "default", "demo"); err == nil {
		t.Log("DeleteJob succeeded unexpectedly (no real cluster)")
	}
}

// -----------------------------------------------------------------------------
// Error handling – nil *rest.Config should cause clientset creation failure
// -----------------------------------------------------------------------------
func TestBatch_ErrorHandling(t *testing.T) {
	k := batch_mustNewKube(t, &config.Config{})
	repo := NewBatch(k)
	ctx := context.Background()

	// an empty config – Host is missing, clientset will return an error
	emptyCfg := &rest.Config{}

	cases := []struct {
		name string
		call func() error
	}{
		{
			name: "ListJobs",
			call: func() error {
				_, err := repo.ListJobs(ctx, emptyCfg, "default")
				return err
			},
		},
		{
			name: "ListJobsByLabel",
			call: func() error {
				_, err := repo.ListJobsByLabel(ctx, emptyCfg, "default", "app=demo")
				return err
			},
		},
		{
			name: "CreateJob",
			call: func() error {
				_, err := repo.CreateJob(ctx, emptyCfg, "default", "demo", nil, nil, nil)
				return err
			},
		},
		{
			name: "DeleteJob",
			call: func() error {
				return repo.DeleteJob(ctx, emptyCfg, "default", "demo")
			},
		},
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
func TestBatch_MethodBehavior(t *testing.T) {
	k := batch_mustNewKube(t, &config.Config{})
	repo := NewBatch(k)

	ctx := context.Background()
	restCfg := &rest.Config{Host: "https://example.invalid"}

	t.Run("ListJobs", func(t *testing.T) {
		_, err := repo.ListJobs(ctx, restCfg, "default")
		if err != nil {
			t.Logf("ListJobs returned expected error: %v", err)
		}
	})

	t.Run("ListJobsByLabel", func(t *testing.T) {
		_, err := repo.ListJobsByLabel(ctx, restCfg, "default", "app=demo")
		if err != nil {
			t.Logf("ListJobsByLabel returned expected error: %v", err)
		}
	})

	t.Run("CreateJob", func(t *testing.T) {
		_, err := repo.CreateJob(ctx, restCfg, "default", "demo", map[string]string{"app": "demo"}, nil, nil)
		if err != nil {
			t.Logf("CreateJob returned expected error: %v", err)
		}
	})

	t.Run("DeleteJob", func(t *testing.T) {
		err := repo.DeleteJob(ctx, restCfg, "default", "demo")
		if err != nil {
			t.Logf("DeleteJob returned expected error: %v", err)
		}
	})
}

// -----------------------------------------------------------------------------
// Integration‑style sequential calls
// -----------------------------------------------------------------------------
func TestBatch_IntegrationPatterns(t *testing.T) {
	k := batch_mustNewKube(t, &config.Config{})
	repo := NewBatch(k)

	ctx := context.Background()
	restCfg := &rest.Config{Host: "https://example.invalid"}

	// List
	if _, err := repo.ListJobs(ctx, restCfg, "default"); err != nil {
		t.Logf("ListJobs error (expected): %v", err)
	}
	// List by label
	if _, err := repo.ListJobsByLabel(ctx, restCfg, "default", "app=demo"); err != nil {
		t.Logf("ListJobsByLabel error (expected): %v", err)
	}
	// Create
	if _, err := repo.CreateJob(ctx, restCfg, "default", "demo", nil, nil, nil); err != nil {
		t.Logf("CreateJob error (expected): %v", err)
	}
	// Delete
	if err := repo.DeleteJob(ctx, restCfg, "default", "demo"); err != nil {
		t.Logf("DeleteJob error (expected): %v", err)
	}
}

// -----------------------------------------------------------------------------
// Benchmarks
// -----------------------------------------------------------------------------
func BenchmarkBatch_Creation(b *testing.B) {
	cfg := &config.Config{}
	k := batch_mustNewKube(b, cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if repo := NewBatch(k); repo == nil {
			b.Fatal("failed to create batch repo")
		}
	}
}

func BenchmarkBatch_MethodCalls(b *testing.B) {
	k := batch_mustNewKube(b, &config.Config{})
	repo := NewBatch(k)
	ctx := context.Background()
	restCfg := &rest.Config{Host: "https://example.invalid"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.ListJobs(ctx, restCfg, "default")
		_, _ = repo.ListJobsByLabel(ctx, restCfg, "default", "app=demo")
		_, _ = repo.CreateJob(ctx, restCfg, "default", "demo", nil, nil, nil)
		_ = repo.DeleteJob(ctx, restCfg, "default", "demo")
	}
}

// -----------------------------------------------------------------------------
// Concurrent access – ensure the clientset cache is safe for parallel use
// -----------------------------------------------------------------------------
func TestBatch_ConcurrentAccess(t *testing.T) {
	k := batch_mustNewKube(t, &config.Config{})
	repo := NewBatch(k)

	ctx := context.Background()
	restCfg := &rest.Config{Host: "https://example.invalid"}

	const workers = 10
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_, _ = repo.ListJobs(ctx, restCfg, "default")
			_, _ = repo.ListJobsByLabel(ctx, restCfg, "default", "app=demo")
			_, _ = repo.CreateJob(ctx, restCfg, "default", "demo", nil, nil, nil)
			_ = repo.DeleteJob(ctx, restCfg, "default", "demo")
		}()
	}
	wg.Wait()
}

// -----------------------------------------------------------------------------
// Type assertions
// -----------------------------------------------------------------------------
func TestBatch_TypeAssertions(t *testing.T) {
	k := batch_mustNewKube(t, &config.Config{})
	repo := NewBatch(k)

	b, ok := repo.(*batch)
	if !ok {
		t.Fatal("could not cast repository to *batch")
	}
	if b.kube != k {
		t.Error("batch.kube field not set correctly")
	}
	var _ oscore.KubeBatchRepo = b
}

// -----------------------------------------------------------------------------
// Edge cases – nil Kube should panic when any method is called
// -----------------------------------------------------------------------------
func TestBatch_EdgeCases(t *testing.T) {
	// Nil Kube
	t.Run("nil_kube", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when using batch with nil Kube")
			}
		}()
		s := &batch{kube: nil}
		ctx := context.Background()
		_, _ = s.ListJobs(ctx, &rest.Config{Host: "https://example.invalid"}, "default")
	})

	// Background context works (method still returns error from clientset)
	t.Run("background_context", func(t *testing.T) {
		k := batch_mustNewKube(t, &config.Config{})
		s := NewBatch(k)
		ctx := context.Background()
		_, err := s.ListJobs(ctx, &rest.Config{Host: "https://example.invalid"}, "default")
		if err == nil {
			t.Log("ListJobs succeeded with background context (no real cluster)")
		}
	})

	// TODO context works
	t.Run("todo_context", func(t *testing.T) {
		k := batch_mustNewKube(t, &config.Config{})
		s := NewBatch(k)
		ctx := context.TODO()
		_, err := s.ListJobs(ctx, &rest.Config{Host: "https://example.invalid"}, "default")
		if err == nil {
			t.Log("ListJobs succeeded with TODO context (no real cluster)")
		}
		_ = ctx // silence unused warning
	})
}
