package juju

import (
	"context"
	"testing"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

// -----------------------------------------------------------------------------
// Helper – works for both normal tests (testing.T) and benchmarks (testing.B)
// -----------------------------------------------------------------------------
func appOffers_mustNewJuju(t testing.TB) *Juju {
	// A zero‑value Juju cannot talk to a controller; this is perfect for the
	// error‑path tests.
	j := New(&config.Config{})
	return j
}

// -----------------------------------------------------------------------------
// Construction / interface compliance
// -----------------------------------------------------------------------------
func TestNewApplicationOffers(t *testing.T) {
	j := appOffers_mustNewJuju(t)
	repo := NewApplicationOffers(j)
	if repo == nil {
		t.Fatal("expected ApplicationOffers repository to be created, got nil")
	}
	var _ core.FacilityOffersRepo = repo
}

// -----------------------------------------------------------------------------
// Structure validation
// -----------------------------------------------------------------------------
func TestApplicationOffers_Structure(t *testing.T) {
	j := appOffers_mustNewJuju(t)
	repo := NewApplicationOffers(j)

	ao, ok := repo.(*applicationOffers)
	if !ok {
		t.Fatal("expected *applicationOffers, got a different type")
	}
	if ao.juju == nil {
		t.Error("expected juju field to be set, got nil")
	}
}

// -----------------------------------------------------------------------------
// Error handling – the zero‑value Juju cannot create a connection, so every
// method should fail (error or panic).
// -----------------------------------------------------------------------------
func TestApplicationOffers_ErrorHandling(t *testing.T) {
	j := appOffers_mustNewJuju(t)
	repo := NewApplicationOffers(j)

	ctx := context.Background()

	_ = ctx // silence unused warning – the method does not use the context
	expectErrorOrPanic(t, func() error {
		_, err := repo.GetConsumeDetails(context.Background(),
			"offer:admin/model.offer")
		return err
	})
}

// -----------------------------------------------------------------------------
// Edge case – nil Juju should panic when any method is invoked.
// -----------------------------------------------------------------------------
func TestApplicationOffers_EdgeCases(t *testing.T) {
	t.Run("nil_juju", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when using ApplicationOffers with nil Juju")
			}
		}()

		ao := &applicationOffers{juju: nil}
		_, _ = ao.GetConsumeDetails(context.Background(),
			"offer:admin/model.offer")
	})
}

// -----------------------------------------------------------------------------
// Benchmark – creation only (method calls would need a real controller)
// -----------------------------------------------------------------------------
func BenchmarkApplicationOffers_Creation(b *testing.B) {
	j := appOffers_mustNewJuju(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if repo := NewApplicationOffers(j); repo == nil {
			b.Fatal("failed to create ApplicationOffers repo")
		}
	}
}
