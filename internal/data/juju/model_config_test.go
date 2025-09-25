package juju

import (
	"context"
	"testing"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

// newBrokenJuju returns a *Juju whose connection method will always fail.
// An empty Config causes NewConnection to return an error, which is
// sufficient for our error‑path tests.
func newBrokenJuju(t *testing.T) *Juju {
	t.Helper()
	// Empty config → newConnection will return an error.
	cfg := &config.Config{}
	return New(cfg)
}

// --------------------------------------------------------------------------
// TestModelConfig_List_Error
// --------------------------------------------------------------------------
func TestModelConfig_List_Error(t *testing.T) {
	j := newBrokenJuju(t) // connection will fail
	repo := NewModelConfig(j)

	_, err := repo.List(context.Background(), "some-uuid")
	if err == nil {
		t.Fatalf("expected error from List when connection fails, got nil")
	}
}

// --------------------------------------------------------------------------
// TestModelConfig_Set_Error
// --------------------------------------------------------------------------
func TestModelConfig_Set_Error(t *testing.T) {
	j := newBrokenJuju(t)
	repo := NewModelConfig(j)

	err := repo.Set(context.Background(), "some-uuid", map[string]any{"key": "value"})
	if err == nil {
		t.Fatalf("expected error from Set when connection fails, got nil")
	}
}

// --------------------------------------------------------------------------
// TestModelConfig_Unset_Error
// --------------------------------------------------------------------------
func TestModelConfig_Unset_Error(t *testing.T) {
	j := newBrokenJuju(t)
	repo := NewModelConfig(j)

	err := repo.Unset(context.Background(), "some-uuid", "key")
	if err == nil {
		t.Fatalf("expected error from Unset when connection fails, got nil")
	}
}

// --------------------------------------------------------------------------
// Compile‑time interface check (optional, ensures our repo implements the
// expected interface even if the implementation changes later).
// --------------------------------------------------------------------------
var _ core.ScopeConfigRepo = (*modelConfig)(nil)
