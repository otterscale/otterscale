package core

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/canonical/gomaasclient/entity"
	"github.com/juju/juju/api/client/action"
	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/core/base"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"
	"github.com/stretchr/testify/assert"
)

// Mock MachineRepo
type mockMachineRepo struct {
	listCalled bool
}

func (m *mockMachineRepo) List(ctx context.Context) ([]Machine, error) {
	m.listCalled = true
	return []Machine{
		{
			Machine: &entity.Machine{
				SystemID:            "id1",
				WorkloadAnnotations: map[string]string{"juju-model-uuid": "scope"},
			},
			LastCommissioned: time.Now(),
		},
	}, nil
}

func (m *mockMachineRepo) Get(ctx context.Context, systemID string) (*Machine, error) {
	return &Machine{
		Machine: &entity.Machine{
			SystemID:            systemID,
			WorkloadAnnotations: map[string]string{"juju-model-uuid": "scope"},
		},
		LastCommissioned: time.Now(),
	}, nil
}

func (m *mockMachineRepo) Release(ctx context.Context, systemID string, params *entity.MachineReleaseParams) (*Machine, error) {
	return &Machine{Machine: &entity.Machine{SystemID: systemID}}, nil
}

func (m *mockMachineRepo) PowerOff(ctx context.Context, systemID string, params *entity.MachinePowerOffParams) (*Machine, error) {
	return &Machine{Machine: &entity.Machine{SystemID: systemID}}, nil
}

func (m *mockMachineRepo) Commission(ctx context.Context, systemID string, params *entity.MachineCommissionParams) (*Machine, error) {
	return &Machine{Machine: &entity.Machine{SystemID: systemID}}, nil
}

// Mock TagRepo
type machineMockTag struct{}

func (m *machineMockTag) AddMachines(ctx context.Context, name string, machineIDs []string) error {
	return nil
}

func (m *machineMockTag) RemoveMachines(ctx context.Context, name string, machineIDs []string) error {
	return nil
}
func (m *machineMockTag) List(ctx context.Context) ([]Tag, error)            { return nil, nil }
func (m *machineMockTag) Get(ctx context.Context, name string) (*Tag, error) { return nil, nil }
func (m *machineMockTag) Create(ctx context.Context, name, comment string) (*Tag, error) {
	return nil, nil
}
func (m *machineMockTag) Delete(ctx context.Context, name string) error { return nil }

// Mock MachineManagerRepo
type mockMachineManagerRepo struct{}

func (m *mockMachineManagerRepo) AddMachines(ctx context.Context, uuid string, params []params.AddMachineParams) error {
	return nil
}

func (m *mockMachineManagerRepo) DestroyMachines(ctx context.Context, uuid string, force, keep, dryRun bool, maxWait *time.Duration, machines ...string) error {
	return nil
}

// Mock ServerRepo
type mockServerRepo struct{}

func (m *mockServerRepo) Get(ctx context.Context, name string) (string, error) { return "value", nil }
func (m *mockServerRepo) Update(ctx context.Context, name, value string) error { return nil }

// Mock ClientRepo
type mockClientRepo struct{}

func (m *mockClientRepo) Status(ctx context.Context, uuid string, patterns []string) (*params.FullStatus, error) {
	return nil, nil
}

// Mock ActionRepo
type mockActionRepo struct{}

func (m *mockActionRepo) RunCommand(ctx context.Context, uuid, uname, cmd string) (string, error) {
	return "action-id", nil
}

func (m *mockActionRepo) GetResult(ctx context.Context, uuid, actionID string) (*action.ActionResult, error) {
	return &action.ActionResult{Status: "test"}, nil
}

// RunAction is a mock implementation to satisfy the ActionRepo interface.
func (m *mockActionRepo) RunAction(ctx context.Context, uuid, uname, actionName string, params map[string]interface{}) (string, error) {
	return "action-id", nil
}

// List is a mock implementation to satisfy the ActionRepo interface.
func (m *mockActionRepo) List(ctx context.Context, uuid, action string) (map[string]ActionSpec, error) {
	return map[string]ActionSpec{}, nil
}

// Mock FacilityRepo
type machineMockFacility struct{}

func (m *machineMockFacility) GetConfig(ctx context.Context, uuid, name string) (map[string]interface{}, error) {
	return map[string]interface{}{"osd-devices": map[string]interface{}{"value": "/dev/sda"}}, nil
}

// AddUnits is a mock implementation to satisfy the FacilityRepo interface.
func (m *machineMockFacility) AddUnits(ctx context.Context, uuid, app string, numUnits int, placements []instance.Placement) ([]string, error) {
	return []string{}, nil
}

// Consume is a mock implementation to satisfy the FacilityRepo interface.
func (m *machineMockFacility) Consume(ctx context.Context, uuid string, args *crossmodel.ConsumeApplicationArgs) error {
	return nil
}

// Create is a mock implementation to satisfy the FacilityRepo interface.
func (m *machineMockFacility) Create(ctx context.Context, uuid, name, app, series, channel string, numUnits, expose int, base *base.Base, placements []instance.Placement, cons *constraints.Value, dryRun bool) (*application.DeployInfo, error) {
	return &application.DeployInfo{}, nil
}

// CreateRelation is a mock implementation to satisfy the FacilityRepo interface.
func (m *machineMockFacility) CreateRelation(ctx context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error) {
	return &params.AddRelationResults{}, nil
}

// Delete is a mock implementation to satisfy the FacilityRepo interface.
func (m *machineMockFacility) Delete(ctx context.Context, uuid, name string, force, dryRun bool) error {
	return nil
}

// DeleteRelation is a mock implementation to satisfy the FacilityRepo interface.
func (m *machineMockFacility) DeleteRelation(ctx context.Context, uuid string, relationID int) error {
	return nil
}

// Expose is a mock implementation to satisfy the FacilityRepo interface.
func (m *machineMockFacility) Expose(ctx context.Context, uuid, name string, endpoints map[string]params.ExposedEndpoint) error {
	return nil
}

// Update is a mock implementation to satisfy the FacilityRepo interface.
func (m *machineMockFacility) Update(ctx context.Context, uuid, name string, config string) error {
	return nil
}

// GetLeader is a mock implementation to satisfy the FacilityRepo interface.
func (m *machineMockFacility) GetLeader(ctx context.Context, uuid, app string) (string, error) {
	return "", nil
}

// GetUnitInfo is a mock implementation to satisfy the FacilityRepo interface.
func (m *machineMockFacility) GetUnitInfo(ctx context.Context, uuid, unit string) (*application.UnitInfo, error) {
	return &application.UnitInfo{}, nil
}

// ResolveUnitErrors is a mock implementation to satisfy the FacilityRepo interface.
func (m *machineMockFacility) ResolveUnitErrors(ctx context.Context, uuid string, errors []string) error {
	return nil
}

// Mock EventRepos
type mockEventRepo struct{}

func (m *mockEventRepo) Get(ctx context.Context, systemID string) ([]Event, error) {
	return []Event{
		{Type: "Commissioning", Created: "Mon, 02 Jan. 2006 15:04:05"},
	}, nil
}

type mockNodeDeviceRepo struct{}

func (m *mockNodeDeviceRepo) List(ctx context.Context, systemID string, deviceType string) ([]NodeDevice, error) {
	return []NodeDevice{}, nil
}

func TestMachineUseCase_ListMachines(t *testing.T) {
	machineRepo := &mockMachineRepo{}
	uc := NewMachineUseCase(
		machineRepo,
		&mockMachineManagerRepo{},
		&mockNodeDeviceRepo{},
		&mockServerRepo{},
		&mockClientRepo{},
		&machineMockTag{},
		&mockActionRepo{},
		&machineMockFacility{},
		&mockEventRepo{},
	)
	machines, err := uc.ListMachines(context.Background(), "scope")
	assert.NoError(t, err)
	assert.Len(t, machines, 1)
	assert.True(t, machineRepo.listCalled)
}

func TestMachineUseCase_GetMachine(t *testing.T) {
	uc := NewMachineUseCase(
		&mockMachineRepo{},
		&mockMachineManagerRepo{},
		&mockNodeDeviceRepo{},
		&mockServerRepo{},
		&mockClientRepo{},
		&machineMockTag{},
		&mockActionRepo{},
		&machineMockFacility{},
		&mockEventRepo{},
	)
	machine, err := uc.GetMachine(context.Background(), "id1")
	assert.NoError(t, err)
	assert.Equal(t, "id1", machine.SystemID)
}

func TestMachineUseCase_AddMachineTags(t *testing.T) {
	uc := NewMachineUseCase(
		&mockMachineRepo{},
		&mockMachineManagerRepo{},
		&mockNodeDeviceRepo{},
		&mockServerRepo{},
		&mockClientRepo{},
		&machineMockTag{},
		&mockActionRepo{},
		&machineMockFacility{},
		&mockEventRepo{},
	)
	err := uc.AddMachineTags(context.Background(), "id1", []string{"tag1", "tag2"})
	assert.NoError(t, err)
}

func TestMachineUseCase_RemoveMachineTags(t *testing.T) {
	uc := NewMachineUseCase(
		&mockMachineRepo{},
		&mockMachineManagerRepo{},
		&mockNodeDeviceRepo{},
		&mockServerRepo{},
		&mockClientRepo{},
		&machineMockTag{},
		&mockActionRepo{},
		&machineMockFacility{},
		&mockEventRepo{},
	)
	err := uc.RemoveMachineTags(context.Background(), "id1", []string{"tag1", "tag2"})
	assert.NoError(t, err)
}

func TestMachineUseCase_DeleteMachine(t *testing.T) {
	uc := NewMachineUseCase(
		&mockMachineRepo{},
		&mockMachineManagerRepo{},
		&mockNodeDeviceRepo{},
		&mockServerRepo{},
		&mockClientRepo{},
		&machineMockTag{},
		&mockActionRepo{},
		&machineMockFacility{},
		&mockEventRepo{},
	)
	err := uc.DeleteMachine(context.Background(), "id1", false, false)
	assert.NoError(t, err)
}

func TestMachineUseCase_PowerOffMachine(t *testing.T) {
	machineRepo := &mockMachineRepo{}
	uc := NewMachineUseCase(
		machineRepo,
		&mockMachineManagerRepo{},
		&mockNodeDeviceRepo{},
		&mockServerRepo{},
		&mockClientRepo{},
		&machineMockTag{},
		&mockActionRepo{},
		&machineMockFacility{},
		&mockEventRepo{},
	)
	machine, err := uc.PowerOffMachine(context.Background(), "id1", "test")
	assert.NoError(t, err)
	assert.Equal(t, "id1", machine.SystemID)
}

func TestMachineUseCase_GetLastCommissioned(t *testing.T) {
	uc := NewMachineUseCase(
		&mockMachineRepo{},
		&mockMachineManagerRepo{},
		&mockNodeDeviceRepo{},
		&mockServerRepo{},
		&mockClientRepo{},
		&machineMockTag{},
		&mockActionRepo{},
		&machineMockFacility{},
		&mockEventRepo{},
	)
	tm, err := uc.getLastCommissioned(context.Background(), "id1")
	assert.NoError(t, err)
	assert.False(t, tm.IsZero())
}

type badTagRepo struct {
	machineMockTag
}

func (b *badTagRepo) AddMachines(ctx context.Context, name string, machineIDs []string) error {
	return errors.New("fail")
}

func TestMachineUseCase_CreateMachine_TagError(t *testing.T) {
	uc := NewMachineUseCase(
		&mockMachineRepo{},
		&mockMachineManagerRepo{},
		&mockNodeDeviceRepo{},
		&mockServerRepo{},
		&mockClientRepo{},
		&badTagRepo{},
		&mockActionRepo{},
		&machineMockFacility{},
		&mockEventRepo{},
	)
	_, err := uc.CreateMachine(context.Background(), "id1", false, false, false, false, "uuid", []string{"name"})
	assert.Error(t, err)
}
