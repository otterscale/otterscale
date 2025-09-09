package core

import (
	"context"
	"testing"

	"github.com/juju/juju/api/client/action"
	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/core/base"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"
	"github.com/stretchr/testify/assert"
)

// Mock FacilityRepo
type facilitymockFacilityRepo struct{}

func (m *facilitymockFacilityRepo) Create(ctx context.Context, uuid, name, configYAML, charmName, channel string, revision, number int, base *base.Base, placements []instance.Placement, constraint *constraints.Value, trust bool) (*application.DeployInfo, error) {
	return &application.DeployInfo{}, nil
}

func (m *facilitymockFacilityRepo) Update(ctx context.Context, uuid, name, configYAML string) error {
	return nil
}

func (m *facilitymockFacilityRepo) Delete(ctx context.Context, uuid, name string, destroyStorage, force bool) error {
	return nil
}

func (m *facilitymockFacilityRepo) Expose(ctx context.Context, uuid, name string, endpoints map[string]params.ExposedEndpoint) error {
	return nil
}

func (m *facilitymockFacilityRepo) AddUnits(ctx context.Context, uuid, name string, number int, placements []instance.Placement) ([]string, error) {
	return []string{"unit1"}, nil
}

func (m *facilitymockFacilityRepo) ResolveUnitErrors(ctx context.Context, uuid string, units []string) error {
	return nil
}

func (m *facilitymockFacilityRepo) CreateRelation(ctx context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error) {
	return &params.AddRelationResults{}, nil
}

func (m *facilitymockFacilityRepo) DeleteRelation(ctx context.Context, uuid string, id int) error {
	return nil
}

func (m *facilitymockFacilityRepo) GetConfig(ctx context.Context, uuid, name string) (map[string]any, error) {
	return map[string]any{"foo": "bar"}, nil
}

func (m *facilitymockFacilityRepo) GetLeader(ctx context.Context, uuid, name string) (string, error) {
	return "leader", nil
}

func (m *facilitymockFacilityRepo) GetUnitInfo(ctx context.Context, uuid, name string) (*application.UnitInfo, error) {
	return &application.UnitInfo{}, nil
}

// Implement the missing Consume method to satisfy FacilityRepo interface
func (m *facilitymockFacilityRepo) Consume(ctx context.Context, uuid string, args *crossmodel.ConsumeApplicationArgs) error {
	return nil
}

// Mock ServerRepo, ClientRepo, ActionRepo, CharmRepo, MachineRepo
type (
	facilityMockServerRepo struct{}
	facilitymockClientRepo struct{}
	facilityMockActionRepo struct{}
	mockCharmRepo          struct{}
)

func (m *facilityMockServerRepo) Get(ctx context.Context, name string) (string, error) {
	return "val", nil
}
func (m *facilityMockServerRepo) Update(ctx context.Context, name, value string) error { return nil }
func (m *facilitymockClientRepo) Status(ctx context.Context, uuid string, patterns []string) (*params.FullStatus, error) {
	return &params.FullStatus{
		Applications: map[string]params.ApplicationStatus{
			"app1": {Charm: "charm1"},
		},
	}, nil
}

func (m *facilityMockActionRepo) List(ctx context.Context, uuid, appName string) (map[string]ActionSpec, error) {
	return map[string]ActionSpec{"action1": {}}, nil
}

func (m *facilityMockActionRepo) RunCommand(ctx context.Context, uuid, unitName, command string) (string, error) {
	return "id", nil
}

func (m *facilityMockActionRepo) RunAction(ctx context.Context, uuid, unitName, actionName string, parameters map[string]any) (string, error) {
	return "id", nil
}

// Remove mockActionResult and return the correct type for GetResult
func (m *facilityMockActionRepo) GetResult(ctx context.Context, uuid, id string) (*action.ActionResult, error) {
	return &action.ActionResult{}, nil
}

func (m *mockCharmRepo) List(ctx context.Context) ([]Charm, error) {
	return []Charm{{ID: "c1", Type: "charm", Name: "testcharm"}}, nil
}

func (m *mockCharmRepo) Get(ctx context.Context, name string) (*Charm, error) {
	return &Charm{ID: "c1", Type: "charm", Name: name}, nil
}

func (m *mockCharmRepo) ListArtifacts(ctx context.Context, name string) ([]CharmArtifact, error) {
	return []CharmArtifact{{Channel: CharmChannel{Name: "stable"}}}, nil
}

func TestFacilityUseCase_ListFacilities(t *testing.T) {
	uc := NewFacilityUseCase(&facilitymockFacilityRepo{}, &facilityMockServerRepo{}, &facilitymockClientRepo{}, &facilityMockActionRepo{}, &mockCharmRepo{}, &mockMachineRepo{})
	facs, err := uc.ListFacilities(context.Background(), "uuid")
	assert.NoError(t, err)
	assert.Len(t, facs, 1)
	assert.Equal(t, "app1", facs[0].Name)
}

func TestFacilityUseCase_GetFacility(t *testing.T) {
	uc := NewFacilityUseCase(&facilitymockFacilityRepo{}, &facilityMockServerRepo{}, &facilitymockClientRepo{}, &facilityMockActionRepo{}, &mockCharmRepo{}, &mockMachineRepo{})
	fac, err := uc.GetFacility(context.Background(), "uuid", "app1")
	assert.NoError(t, err)
	assert.Equal(t, "app1", fac.Name)
	assert.Contains(t, fac.Metadata.ConfigYAML, "foo")
}

/*
func TestFacilityUseCase_CreateFacility(t *testing.T) {
	uc := NewFacilityUseCase(&facilitymockFacilityRepo{}, &facilityMockServerRepo{}, &facilitymockClientRepo{}, &facilityMockActionRepo{}, &mockCharmRepo{}, &mockMachineRepo{})
	fac, err := uc.CreateFacility(context.Background(), "uuid", "name", "yaml", "charm", "stable", 1, 1, nil, nil, true)
	assert.NoError(t, err)
	assert.NotNil(t, fac)
}
*/

func TestFacilityUseCase_UpdateFacility(t *testing.T) {
	uc := NewFacilityUseCase(&facilitymockFacilityRepo{}, &facilityMockServerRepo{}, &facilitymockClientRepo{}, &facilityMockActionRepo{}, &mockCharmRepo{}, &mockMachineRepo{})
	fac, err := uc.UpdateFacility(context.Background(), "uuid", "name", "yaml")
	assert.NoError(t, err)
	assert.NotNil(t, fac)
}

func TestFacilityUseCase_DeleteFacility(t *testing.T) {
	uc := NewFacilityUseCase(&facilitymockFacilityRepo{}, &facilityMockServerRepo{}, &facilitymockClientRepo{}, &facilityMockActionRepo{}, &mockCharmRepo{}, &mockMachineRepo{})
	err := uc.DeleteFacility(context.Background(), "uuid", "name", false, false)
	assert.NoError(t, err)
}

func TestFacilityUseCase_ListCharms(t *testing.T) {
	uc := NewFacilityUseCase(&facilitymockFacilityRepo{}, &facilityMockServerRepo{}, &facilitymockClientRepo{}, &facilityMockActionRepo{}, &mockCharmRepo{}, &mockMachineRepo{})
	charms, err := uc.ListCharms(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, charms)
}

func TestFacilityUseCase_GetCharm(t *testing.T) {
	uc := NewFacilityUseCase(&facilitymockFacilityRepo{}, &facilityMockServerRepo{}, &facilitymockClientRepo{}, &facilityMockActionRepo{}, &mockCharmRepo{}, &mockMachineRepo{})
	charm, err := uc.GetCharm(context.Background(), "testcharm")
	assert.NoError(t, err)
	assert.Equal(t, "testcharm", charm.Name)
}

func TestFacilityUseCase_ListArtifacts(t *testing.T) {
	uc := NewFacilityUseCase(&facilitymockFacilityRepo{}, &facilityMockServerRepo{}, &facilitymockClientRepo{}, &facilityMockActionRepo{}, &mockCharmRepo{}, &mockMachineRepo{})
	arts, err := uc.ListArtifacts(context.Background(), "testcharm")
	assert.NoError(t, err)
	assert.NotEmpty(t, arts)
}
