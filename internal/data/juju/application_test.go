package juju

import (
	"context"
	"testing"

	"github.com/juju/juju/core/base"
	"github.com/juju/juju/core/crossmodel"
	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

// -----------------------------------------------------------------------------
// Helper – works for both normal tests (testing.T) and benchmarks (testing.B)
// -----------------------------------------------------------------------------
func app_mustNewJuju(t testing.TB) *Juju {
	// The zero‑value Juju cannot contact a controller – this is exactly what we
	// need for the error‑path tests.
	j := New(&config.Config{})

	return j
}

// -----------------------------------------------------------------------------
// Construction / interface compliance
// -----------------------------------------------------------------------------
func TestNewApplication(t *testing.T) {
	j := app_mustNewJuju(t)
	repo := NewApplication(j)
	if repo == nil {
		t.Fatal("expected Application repository to be created, got nil")
	}
	var _ core.FacilityRepo = repo
}

// -----------------------------------------------------------------------------
// Structure validation
// -----------------------------------------------------------------------------
func TestApplication_Structure(t *testing.T) {
	j := app_mustNewJuju(t)
	repo := NewApplication(j)

	app, ok := repo.(*application)
	if !ok {
		t.Fatal("expected *application, got a different type")
	}
	if app.juju == nil {
		t.Error("expected juju field to be set, got nil")
	}
}

// -----------------------------------------------------------------------------
// Error handling – every public method should return an error (or panic) when
// the underlying Juju connection cannot be established.
// -----------------------------------------------------------------------------
func TestApplication_ErrorHandling(t *testing.T) {
	j := app_mustNewJuju(t)
	repo := NewApplication(j)

	ctx := context.Background()

	// Helper that runs a call and expects either an error or a panic.
	expectErrorOrPanic := func(t *testing.T, fn func() error) {
		defer func() {
			if r := recover(); r != nil {
				// Panic is acceptable – treat it as the expected failure.
				t.Logf("recovered panic as expected: %v", r)
			}
		}()

		if err := fn(); err != nil {
			// Non‑nil error is also acceptable.
			return
		}
		// If we get here we received neither an error nor a panic → failure.
		t.Errorf("expected error or panic, got nil")
	}

	// --------------------------------------------------------------------- //
	// Each method is exercised.  The arguments are dummy values – the call
	// will fail long before they are used (at the connection step).
	// --------------------------------------------------------------------- //
	expectErrorOrPanic(t, func() error {
		_, err := repo.Create(ctx, "uuid", "myapp", "", "charm", "", 0, 0,
			(*base.Base)(nil), nil, nil, false)
		return err
	})

	expectErrorOrPanic(t, func() error {
		return repo.Update(ctx, "uuid", "myapp", "")
	})

	expectErrorOrPanic(t, func() error {
		return repo.Delete(ctx, "uuid", "myapp", false, false)
	})

	expectErrorOrPanic(t, func() error {
		return repo.Expose(ctx, "uuid", "myapp", nil)
	})

	expectErrorOrPanic(t, func() error {
		_, err := repo.AddUnits(ctx, "uuid", "myapp", 1, nil)
		return err
	})

	expectErrorOrPanic(t, func() error {
		return repo.ResolveUnitErrors(ctx, "uuid", []string{"unit/0"})
	})

	expectErrorOrPanic(t, func() error {
		_, err := repo.CreateRelation(ctx, "uuid", []string{"app:db", "mysql:db"})
		return err
	})

	expectErrorOrPanic(t, func() error {
		return repo.DeleteRelation(ctx, "uuid", 42)
	})

	expectErrorOrPanic(t, func() error {
		_, err := repo.GetConfig(ctx, "uuid", "myapp")
		return err
	})

	expectErrorOrPanic(t, func() error {
		_, err := repo.GetLeader(ctx, "uuid", "myapp")
		return err
	})

	expectErrorOrPanic(t, func() error {
		_, err := repo.GetUnitInfo(ctx, "uuid", "myapp/0")
		return err
	})
}

// -----------------------------------------------------------------------------
// Edge cases – nil Juju should panic when any method is called.
// -----------------------------------------------------------------------------
func TestApplication_EdgeCases(t *testing.T) {
	// Nil Juju – any method call must panic.
	t.Run("nil_juju", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when using application with nil Juju")
			}
		}()

		app := &application{juju: nil}
		ctx := context.Background()

		// Call a handful of methods – all should panic because `app.juju`
		// is nil and the underlying `connection` method dereferences it.
		_, _ = app.Create(ctx, "uuid", "myapp", "", "", "", 0, 0, nil, nil, nil, false)
		_ = app.Update(ctx, "uuid", "myapp", "")
		_ = app.Delete(ctx, "uuid", "myapp", false, false)
		_ = app.Expose(ctx, "uuid", "myapp", nil)
		_, _ = app.AddUnits(ctx, "uuid", "myapp", 1, nil)
		_ = app.ResolveUnitErrors(ctx, "uuid", []string{"unit/0"})
		_, _ = app.CreateRelation(ctx, "uuid", []string{"a:b", "c:d"})
		_ = app.DeleteRelation(ctx, "uuid", 1)
		_, _ = app.GetConfig(ctx, "uuid", "myapp")
		_, _ = app.GetLeader(ctx, "uuid", "myapp")
		_, _ = app.GetUnitInfo(ctx, "uuid", "myapp/0")
		_ = app.Consume(ctx, "uuid", &crossmodel.ConsumeApplicationArgs{})
	})
}

// -----------------------------------------------------------------------------
// Benchmark – creation only (method calls would require a real controller)
// -----------------------------------------------------------------------------
func BenchmarkApplication_Creation(b *testing.B) {
	j := app_mustNewJuju(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if repo := NewApplication(j); repo == nil {
			b.Fatal("failed to create application repo")
		}
	}
}
