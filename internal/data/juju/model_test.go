package juju

import (
	"context"
	"testing"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/core/status"
	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

/* -------------------------------------------------------------------------- *
 * Helper – creates a minimal UserModelSummary for tests
 * -------------------------------------------------------------------------- */
func newSummary(name string, st status.Status) base.UserModelSummary {
	return base.UserModelSummary{
		UUID:         "uuid-" + name,
		Name:         name,
		Type:         "iaas",
		ProviderType: "aws",
		Life:         "alive",
		Status:       base.Status{Status: st},
		IsController: false,
	}
}

/* -------------------------------------------------------------------------- *
 * Test_List_Error – when the Juju connection fails, List should return an error
 * -------------------------------------------------------------------------- */
func Test_List_Error(t *testing.T) {
	// Empty config makes Juju.connection() fail
	cfg := &config.Config{} // no Juju parameters
	j := New(cfg)

	repo := NewModel(j)

	_, err := repo.List(context.Background())
	if err == nil {
		t.Fatalf("expected error from List when connection fails, got nil")
	}
}

/* -------------------------------------------------------------------------- *
 * Test_Create_Error – when the Juju connection fails, Create should return an error
 * -------------------------------------------------------------------------- */
func Test_Create_Error(t *testing.T) {
	// Empty config makes Juju.connection() fail
	cfg := &config.Config{}
	j := New(cfg)

	repo := NewModel(j)

	_, err := repo.Create(context.Background(), "some-model")
	if err == nil {
		t.Fatalf("expected error from Create when connection fails, got nil")
	}
}

/* -------------------------------------------------------------------------- *
 * TestModel_filterValidModels
 *   - removes models with name == "controller"
 *   - removes models whose status is not considered valid by status.ValidModelStatus
 *   - converts the remaining models to core.Scope
 * -------------------------------------------------------------------------- */
func TestModel_filterValidModels(t *testing.T) {
	m := &model{} // receiver is not used inside filterValidModels

	// Build test data:
	// 1. a normal valid model (kept)
	// 2. a controller model (removed)
	// 3. a model with an invalid status (removed)
	// 4. another normal valid model (kept)
	summaries := []base.UserModelSummary{
		newSummary("app1", status.Available), // keep
		{
			UUID:         "uuid-controller",
			Name:         "controller", // name == "controller" → filtered
			Type:         "iaas",
			ProviderType: "aws",
			Life:         "alive",
			Status:       base.Status{Status: status.Active},
			IsController: true,
		},
		newSummary("broken", status.Blocked), // invalid status → filtered
		newSummary("app2", status.Suspended), // keep
	}

	got := m.filterValidModels(summaries)

	// Expect only two valid models
	if want := 2; len(got) != want {
		t.Fatalf("expected %d scopes, got %d", want, len(got))
	}

	// Verify the UUIDs (DeleteFunc preserves order)
	expectedUUIDs := []string{"uuid-app1", "uuid-app2"}
	for i, sc := range got {
		if sc.UUID != expectedUUIDs[i] {
			t.Errorf("scope %d: expected UUID %q, got %q", i, expectedUUIDs[i], sc.UUID)
		}
	}
}

/* -------------------------------------------------------------------------- *
 * Compile‑time interface check – ensure model implements core.ScopeRepo
 * -------------------------------------------------------------------------- */
var _ core.ScopeRepo = (*model)(nil)
