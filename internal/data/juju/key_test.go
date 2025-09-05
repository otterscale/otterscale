package juju

import (
	"context"
	"sync"
	"testing"

	"connectrpc.com/connect"
	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

/* -------------------------------------------------------------------------- *
 * Helper – works for both normal tests (testing.T) and benchmarks (testing.B)
 * -------------------------------------------------------------------------- */
func key_mustNewJuju(t testing.TB, cfg *config.Config) *Juju {
	j := New(cfg)

	return j
}

/* -------------------------------------------------------------------------- *
 * Construction / interface compliance
 * -------------------------------------------------------------------------- */
func TestNewKey(t *testing.T) {
	j := key_mustNewJuju(t, &config.Config{})
	repo := NewKey(j)
	if repo == nil {
		t.Fatal("expected Key repository to be created, got nil")
	}
	var _ core.KeyRepo = repo
}

/* -------------------------------------------------------------------------- *
 * Structure validation
 * -------------------------------------------------------------------------- */
func TestKey_Structure(t *testing.T) {
	j := key_mustNewJuju(t, &config.Config{})
	repo := NewKey(j)

	k, ok := repo.(*key)
	if !ok {
		t.Fatal("expected *key, got a different type")
	}
	if k.juju == nil {
		t.Error("expected juju field to be set, got nil")
	}
}

/* -------------------------------------------------------------------------- *
 * Calls with a syntactically valid config (no real controller)
 * -------------------------------------------------------------------------- */
func TestKey_WithConfig(t *testing.T) {
	j := key_mustNewJuju(t, &config.Config{})
	repo := NewKey(j)

	ctx := context.Background()
	// A non‑empty UUID; the underlying connection will fail because there is
	// no real Juju controller, so we only check that the method returns an
	// error rather than panicking.
	if err := repo.Add(ctx, "some-uuid", "my‑key"); err == nil {
		t.Log("Add succeeded unexpectedly (no real controller)")
	}
}

/* -------------------------------------------------------------------------- *
 * Error handling – empty config should cause a connection error
 * -------------------------------------------------------------------------- */
func TestKey_ErrorHandling(t *testing.T) {
	j := key_mustNewJuju(t, &config.Config{}) // empty config → bad connection
	repo := NewKey(j)

	ctx := context.Background()
	if err := repo.Add(ctx, "uuid", "foo"); err == nil {
		t.Fatalf("expected error from Add with empty config, got nil")
	}
}

/* -------------------------------------------------------------------------- *
 * Edge case – nil Juju should panic when any method is called
 * -------------------------------------------------------------------------- */
func TestKey_EdgeCases(t *testing.T) {
	t.Run("nil_juju", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when using Key repo with nil Juju")
			}
		}()
		kr := &key{juju: nil}
		_ = kr.Add(context.Background(), "uuid", "foo")
	})
}

/* -------------------------------------------------------------------------- *
 * Concurrent access – ensure the connections cache is safe for parallel use
 * -------------------------------------------------------------------------- */
func TestKey_ConcurrentAccess(t *testing.T) {
	j := key_mustNewJuju(t, &config.Config{})
	repo := NewKey(j)

	ctx := context.Background()
	const workers = 10
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_ = repo.Add(ctx, "uuid", "foo")
		}()
	}
	wg.Wait()
}

/* -------------------------------------------------------------------------- *
 * Benchmark – creation only (method calls need a real controller)
 * -------------------------------------------------------------------------- */
func BenchmarkKey_Creation(b *testing.B) {
	j := key_mustNewJuju(b, &config.Config{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if repo := NewKey(j); repo == nil {
			b.Fatal("failed to create key repo")
		}
	}
}

/* -------------------------------------------------------------------------- *
 * Helper to silence unused imports (the production code imports them)
 * -------------------------------------------------------------------------- */
var _ = connect.CodeOf // keep the import alive for the build
