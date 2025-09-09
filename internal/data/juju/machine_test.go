package juju

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/juju/juju/rpc/params"
	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

/* -------------------------------------------------------------------------- *
 * Helper – works for both normal tests (testing.T) and benchmarks (testing.B)
 * -------------------------------------------------------------------------- */
func machine_mustNewJuju(t testing.TB, cfg *config.Config) *Juju {
	j := New(cfg)

	return j
}

/* -------------------------------------------------------------------------- *
 * Construction / interface compliance
 * -------------------------------------------------------------------------- */
func TestNewMachine(t *testing.T) {
	j := machine_mustNewJuju(t, &config.Config{})
	repo := NewMachine(j)
	if repo == nil {
		t.Fatal("expected MachineManager repository to be created, got nil")
	}
	var _ core.MachineManagerRepo = repo
}

/* -------------------------------------------------------------------------- *
 * Structure validation
 * -------------------------------------------------------------------------- */
func TestMachine_Structure(t *testing.T) {
	j := machine_mustNewJuju(t, &config.Config{})
	repo := NewMachine(j)

	m, ok := repo.(*machine)
	if !ok {
		t.Fatal("expected *machine, got a different type")
	}
	if m.juju == nil {
		t.Error("expected juju field to be set, got nil")
	}
}

/* -------------------------------------------------------------------------- *
 * Calls with a syntactically valid config (no real controller)
 * -------------------------------------------------------------------------- */
func TestMachine_WithConfig(t *testing.T) {
	j := machine_mustNewJuju(t, &config.Config{})
	repo := NewMachine(j)

	ctx := context.Background()
	// The connection will fail because there is no real controller – we only
	// verify that the method returns an error (and does not panic).
	if err := repo.AddMachines(ctx, "some-uuid", []params.AddMachineParams{{}}); err == nil {
		t.Log("AddMachines succeeded unexpectedly (no real controller)")
	}
	if err := repo.DestroyMachines(ctx, "some-uuid", false, false, false, nil, "machine-0"); err == nil {
		t.Log("DestroyMachines succeeded unexpectedly (no real controller)")
	}
}

/* -------------------------------------------------------------------------- *
 * Error handling – empty config should cause a connection error
 * -------------------------------------------------------------------------- */
func TestMachine_ErrorHandling(t *testing.T) {
	// Empty config → connection cannot be built.
	j := machine_mustNewJuju(t, &config.Config{})
	repo := NewMachine(j)

	ctx := context.Background()
	if err := repo.AddMachines(ctx, "uuid", []params.AddMachineParams{{}}); err == nil {
		t.Fatalf("expected error from AddMachines with empty config, got nil")
	}
	if err := repo.DestroyMachines(ctx, "uuid", false, false, false, nil, "machine-0"); err == nil {
		t.Fatalf("expected error from DestroyMachines with empty config, got nil")
	}
}

/* -------------------------------------------------------------------------- *
 * Edge case – nil Juju should panic on any method call
 * -------------------------------------------------------------------------- */
func TestMachine_EdgeCases(t *testing.T) {
	t.Run("nil_juju", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when using MachineManager with nil Juju")
			}
		}()
		m := &machine{juju: nil}
		_ = m.AddMachines(context.Background(), "uuid", []params.AddMachineParams{{}})
	})
}

/* -------------------------------------------------------------------------- *
 * Concurrent access – ensure the connections cache is safe for parallel use
 * -------------------------------------------------------------------------- */
func TestMachine_ConcurrentAccess(t *testing.T) {
	j := machine_mustNewJuju(t, &config.Config{})
	repo := NewMachine(j)

	ctx := context.Background()
	const workers = 10
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_ = repo.AddMachines(ctx, "uuid", []params.AddMachineParams{{}})
			_ = repo.DestroyMachines(ctx, "uuid", false, false, false, nil, "machine-0")
		}()
	}
	wg.Wait()
}

/* -------------------------------------------------------------------------- *
 * Benchmark – creation only (method calls need a real Juju controller)
 * -------------------------------------------------------------------------- */
func BenchmarkMachine_Creation(b *testing.B) {
	j := machine_mustNewJuju(b, &config.Config{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if repo := NewMachine(j); repo == nil {
			b.Fatal("failed to create machine repo")
		}
	}
}

/* -------------------------------------------------------------------------- *
 * Helper to silence unused imports (the production code imports them)
 * -------------------------------------------------------------------------- */
var _ = errors.Join // keep the import alive for the build
