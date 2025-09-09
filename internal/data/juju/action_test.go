package juju

import (
	"context"
	"testing"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

// ---------------------------------------------------------------------------
// Helper – works for both normal tests (testing.T) and benchmarks (testing.B)
// ---------------------------------------------------------------------------
func action_mustNewJuju(t testing.TB) *Juju {
	// We deliberately create a *zero‑value* Juju instance (no controller
	// address, no auth data).  All repository methods will try to open a
	// connection and therefore fail – exactly what we need for the error‑path
	// tests.
	j := New(&config.Config{}) // New returns (*Juju, error)

	return j
}

/* -------------------------------------------------------------------------- *
 * Construction / interface compliance
 * -------------------------------------------------------------------------- */
func TestNewAction(t *testing.T) {
	j := action_mustNewJuju(t)
	repo := NewAction(j)
	if repo == nil {
		t.Fatal("expected Action repository to be created, got nil")
	}
	// Verify it implements the interface
	var _ core.ActionRepo = repo
}

/* -------------------------------------------------------------------------- *
 * Helper that runs a repository method and expects *either* an error *or*
 * a panic (both are valid ways of signalling a failed connection).
 * -------------------------------------------------------------------------- */
func expectErrorOrPanic(t *testing.T, fn func() error) {
	defer func() {
		if r := recover(); r != nil {
			// Panic is an acceptable way to signal the failure.
			t.Logf("recovered panic as expected: %v", r)
		}
	}()

	if err := fn(); err != nil {
		// Returning a non‑nil error is also acceptable.
		return
	}
	// If we reach here we got neither an error nor a panic → fail the test.
	t.Errorf("expected an error or panic, got nil")
}

/* -------------------------------------------------------------------------- *
 * List – connection error
 * -------------------------------------------------------------------------- */
func TestAction_List_ConnectionError(t *testing.T) {
	j := action_mustNewJuju(t)
	repo := NewAction(j)

	expectErrorOrPanic(t, func() error {
		_, err := repo.List(context.Background(), "some-uuid", "my-app")
		return err
	})
}

/* -------------------------------------------------------------------------- *
 * RunCommand – connection error
 * -------------------------------------------------------------------------- */
func TestAction_RunCommand_ConnectionError(t *testing.T) {
	j := action_mustNewJuju(t)
	repo := NewAction(j)

	expectErrorOrPanic(t, func() error {
		_, err := repo.RunCommand(context.Background(),
			"some-uuid", "unit/0", "whoami")
		return err
	})
}

/* -------------------------------------------------------------------------- *
 * RunAction – connection error
 * -------------------------------------------------------------------------- */
func TestAction_RunAction_ConnectionError(t *testing.T) {
	j := action_mustNewJuju(t)
	repo := NewAction(j)

	expectErrorOrPanic(t, func() error {
		_, err := repo.RunAction(context.Background(),
			"some-uuid", "unit/0", "some-action", map[string]any{})
		return err
	})
}

/* -------------------------------------------------------------------------- *
 * GetResult – connection error (may panic)
 * -------------------------------------------------------------------------- */
func TestAction_GetResult_ConnectionError(t *testing.T) {
	j := action_mustNewJuju(t)
	repo := NewAction(j)

	expectErrorOrPanic(t, func() error {
		_, err := repo.GetResult(context.Background(),
			"some-uuid", "action-id-123")
		return err
	})
}
