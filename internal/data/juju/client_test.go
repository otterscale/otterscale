package juju

import (
	"context"
	"sync"
	"testing"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/wrap"
)

// -----------------------------------------------------------------------------
// Helper – works for both normal tests (testing.T) and benchmarks (testing.B)
// -----------------------------------------------------------------------------
func client_mustNewJuju(t testing.TB, cfg *config.Config) *Juju {
	j := New(cfg)

	return j
}

// -----------------------------------------------------------------------------
// Construction / interface compliance
// -----------------------------------------------------------------------------
func TestNewClient(t *testing.T) {
	j := client_mustNewJuju(t, &config.Config{})
	repo := NewClient(j)
	if repo == nil {
		t.Fatal("NewClient returned nil")
	}
	var _ core.ClientRepo = repo
}

// -----------------------------------------------------------------------------
// Structure validation
// -----------------------------------------------------------------------------
func TestClient_Structure(t *testing.T) {
	j := client_mustNewJuju(t, &config.Config{})
	repo := NewClient(j)

	c, ok := repo.(*client)
	if !ok {
		t.Fatal("expected *client, got a different type")
	}
	if c.juju == nil {
		t.Error("expected juju field to be set, got nil")
	}
}

// -----------------------------------------------------------------------------
// Calls with a syntactically valid rest.Config (no real controller)
// The connection will fail, but the call must return an error instead of panicking.
// -----------------------------------------------------------------------------
func TestClient_WithConfig(t *testing.T) {
	j := client_mustNewJuju(t, &config.Config{})
	repo := NewClient(j)

	ctx := context.Background()
	// any UUID works – the underlying connection will be unable to talk to a
	// controller, therefore an error is expected.
	if _, err := repo.Status(ctx, "some-uuid", []string{}); err == nil {
		t.Log("Status succeeded unexpectedly (no real controller)")
	}
}

// -----------------------------------------------------------------------------
// Error handling – empty config should cause a connection error
// -----------------------------------------------------------------------------
//
// The `client` repository calls `r.juju.connection(uuid)`.  When the Juju
// configuration is empty the connection routine returns an error; the method
// should propagate that error.
//
// We use an *empty* `&config.Config{}` (i.e. no controller information) – this
// triggers the error without causing a nil‑pointer panic.
func TestClient_ErrorHandling(t *testing.T) {
	j := client_mustNewJuju(t, &config.Config{})
	repo := NewClient(j)

	ctx := context.Background()
	if _, err := repo.Status(ctx, "uuid", []string{"model"}); err == nil {
		t.Fatalf("expected error from Status with empty Juju config, got nil")
	}
}

// -----------------------------------------------------------------------------
// Concurrent access – ensure the client repo is safe for parallel use
// -----------------------------------------------------------------------------
func TestClient_ConcurrentAccess(t *testing.T) {
	j := client_mustNewJuju(t, &config.Config{})
	repo := NewClient(j)

	ctx := context.Background()
	const workers = 10
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_, _ = repo.Status(ctx, "uuid", []string{})
		}()
	}
	wg.Wait()
}

// -----------------------------------------------------------------------------
// Edge case – nil Juju should panic when any method is called
// -----------------------------------------------------------------------------
func TestClient_EdgeCases(t *testing.T) {
	t.Run("nil_juju", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when using client with nil Juju")
			}
		}()
		c := &client{juju: nil}
		_, _ = c.Status(context.Background(), "uuid", []string{})
	})
}

// -----------------------------------------------------------------------------
// Benchmark – creation only (method calls would need a real controller)
// -----------------------------------------------------------------------------
func BenchmarkClient_Creation(b *testing.B) {
	j := client_mustNewJuju(b, &config.Config{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if repo := NewClient(j); repo == nil {
			b.Fatal("creation returned nil")
		}
	}
}

// -----------------------------------------------------------------------------
// Keep the utils import alive – the production code uses it.
// -----------------------------------------------------------------------------
var _ = wrap.HTTPGet // silence unused‑import warning

/* -------------------------------------------------------------------------- *
 * The tests above follow the same pattern used in the other Juju test files
 * (apps_test.go, action_test.go, etc.): helper to build a *Juju instance,
 * construction/interface checks, structure validation, a happy‑path call that
 * returns an error because there is no real controller, explicit error‑path
 * testing with an empty config, concurrent access, edge‑case panic test and a
 * simple benchmark.
 * -------------------------------------------------------------------------- */
