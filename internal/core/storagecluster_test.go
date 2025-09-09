package core

import (
	"context"
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

// Mock ActionRepo for storageConfig
type clusterMockActionRepo struct {
	call int
}

// List is a stub to satisfy the ActionRepo interface.
func (m *clusterMockActionRepo) List(ctx context.Context, uuid, leader string) (map[string]ActionSpec, error) {
	return map[string]ActionSpec{}, nil
}

func (m *clusterMockActionRepo) RunCommand(ctx context.Context, uuid, leader, command string) (string, error) {
	switch {
	case command == "ceph config generate-minimal-conf && ceph auth get client.admin":
		return "action-id-ceph-config", nil
	case command == "radosgw-admin user list":
		return "action-id-user-list", nil
	case command == "radosgw-admin user info --uid=otterscale --format=json":
		return "action-id-user-info", nil
	default:
		return "action-id-ceph-config", nil
	}
}

func (m *clusterMockActionRepo) RunAction(ctx context.Context, uuid, leader, action string, params map[string]any) (string, error) {
	return "action-id", nil
}

func (m *clusterMockActionRepo) GetResult(ctx context.Context, uuid, id string) (*action.ActionResult, error) {
	switch id {
	case "action-id-ceph-config":
		return &action.ActionResult{
			Status: "completed",
			Output: map[string]any{
				"stdout": `
[global]
fsid = test-fsid
mon_host = 1.2.3.4
[client.admin]
key = test-key
`,
			},
		}, nil
	case "action-id-user-list":
		return &action.ActionResult{
			Status: "completed",
			Output: map[string]any{
				"stdout": `["otterscale"]`,
			},
		}, nil
	case "action-id-user-info":
		return &action.ActionResult{
			Status: "completed",
			Output: map[string]any{
				"stdout": `{"keys":[{"access_key":"ak","secret_key":"sk"}]}`,
			},
		}, nil
	default:
		// fallback: return ceph config for first call
		return &action.ActionResult{
			Status: "completed",
			Output: map[string]any{
				"stdout": `
[global]
fsid = test-fsid
mon_host = 1.2.3.4
[client.admin]
key = test-key
`,
			},
		}, nil
	}
}

// Mock CephClusterRepo
type mockCephClusterRepo struct{}

func (m *mockCephClusterRepo) ListMONs(ctx context.Context, config *StorageConfig) ([]MON, error) {
	return []MON{{Name: "mon1"}}, nil
}

func (m *mockCephClusterRepo) ListOSDs(ctx context.Context, config *StorageConfig) ([]OSD, error) {
	return []OSD{{ID: 1, Name: "osd.1", Hostname: "host1"}}, nil
}

func (m *mockCephClusterRepo) DoSMART(ctx context.Context, config *StorageConfig, who string) (map[string][]string, error) {
	return map[string][]string{"osd.1": {"OK"}}, nil
}

func (m *mockCephClusterRepo) ListPools(ctx context.Context, config *StorageConfig) ([]Pool, error) {
	return []Pool{{ID: 1, Name: "pool1", Type: "replicated"}}, nil
}

func (m *mockCephClusterRepo) ListPoolsByApplication(ctx context.Context, config *StorageConfig, application string) ([]Pool, error) {
	return []Pool{{ID: 2, Name: "pool2", Type: "replicated"}}, nil
}

func (m *mockCephClusterRepo) CreatePool(ctx context.Context, config *StorageConfig, pool, poolType string) error {
	return nil
}

func (m *mockCephClusterRepo) DeletePool(ctx context.Context, config *StorageConfig, pool string) error {
	return nil
}

func (m *mockCephClusterRepo) EnableApplication(ctx context.Context, config *StorageConfig, pool, application string) error {
	return nil
}

func (m *mockCephClusterRepo) GetParameter(ctx context.Context, config *StorageConfig, pool, key string) (string, error) {
	return "true", nil
}

func (m *mockCephClusterRepo) SetParameter(ctx context.Context, config *StorageConfig, pool, key, value string) error {
	return nil
}

func (m *mockCephClusterRepo) SetQuota(ctx context.Context, config *StorageConfig, pool string, maxBytes, maxObjects uint64) error {
	return nil
}

func (m *mockCephClusterRepo) GetQuota(ctx context.Context, config *StorageConfig, pool string) (uint64, uint64, error) {
	return 100, 10, nil
}

func (m *mockCephClusterRepo) GetECProfile(ctx context.Context, config *StorageConfig, name string) (string, string, error) {
	return "2", "1", nil
}

// Mock FacilityRepo
type storageMockFacilityRepo struct{}

// GetConfig is a mock implementation to satisfy the FacilityRepo interface.
func (m *storageMockFacilityRepo) GetConfig(ctx context.Context, facility string, modelUUID string) (map[string]any, error) {
	return map[string]any{}, nil
}

func (m *storageMockFacilityRepo) Get(ctx context.Context, facility string) (*Facility, error) {
	return &Facility{Name: facility}, nil
}

// AddUnits is a mock implementation to satisfy the FacilityRepo interface.
func (m *storageMockFacilityRepo) AddUnits(ctx context.Context, facility string, app string, units int, placements []instance.Placement) ([]string, error) {
	return []string{}, nil
}

// Consume is a mock implementation to satisfy the FacilityRepo interface.
func (m *storageMockFacilityRepo) Consume(ctx context.Context, facility string, args *crossmodel.ConsumeApplicationArgs) error {
	return nil
}

// Create is a mock implementation to satisfy the FacilityRepo interface.
func (m *storageMockFacilityRepo) Create(
	ctx context.Context,
	modelUUID, facility, app, channel, series string,
	units, minUnits int,
	base *base.Base,
	placements []instance.Placement,
	constraints *constraints.Value,
	expose bool,
) (*application.DeployInfo, error) {
	return &application.DeployInfo{
		Name: app,
	}, nil
}

// CreateRelation is a mock implementation to satisfy the FacilityRepo interface.
func (m *storageMockFacilityRepo) CreateRelation(ctx context.Context, facility string, endpoints []string) (*params.AddRelationResults, error) {
	return &params.AddRelationResults{}, nil
}

// Delete is a mock implementation to satisfy the FacilityRepo interface.
func (m *storageMockFacilityRepo) Delete(ctx context.Context, facility string, app string, force bool, destroyStorage bool) error {
	return nil
}

// DeleteRelation is a mock implementation to satisfy the FacilityRepo interface.
func (m *storageMockFacilityRepo) DeleteRelation(ctx context.Context, facility string, relationID int) error {
	return nil
}

// Expose is a mock implementation to satisfy the FacilityRepo interface.
func (m *storageMockFacilityRepo) Expose(ctx context.Context, facility string, app string, endpoints map[string]params.ExposedEndpoint) error {
	return nil
}

// GetLeader is a mock implementation to satisfy the FacilityRepo interface.
func (m *storageMockFacilityRepo) GetLeader(ctx context.Context, facility string, app string) (string, error) {
	return "leader-unit/0", nil
}

// GetUnitInfo is a mock implementation to satisfy the FacilityRepo interface.
func (m *storageMockFacilityRepo) GetUnitInfo(ctx context.Context, facility string, unit string) (*application.UnitInfo, error) {
	return &application.UnitInfo{}, nil
}

// ResolveUnitErrors is a mock implementation to satisfy the FacilityRepo interface.
func (m *storageMockFacilityRepo) ResolveUnitErrors(ctx context.Context, facility string, units []string) error {
	return nil
}

// Update is a mock implementation to satisfy the FacilityRepo interface.
func (m *storageMockFacilityRepo) Update(ctx context.Context, facility string, app string, config string) error {
	return nil
}

// Mock MachineRepo
type storageMockMachineRepo struct{}

func (m *storageMockMachineRepo) List(ctx context.Context) ([]Machine, error) {
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

// PowerOff is a mock implementation to satisfy the MachineRepo interface.
func (m *storageMockMachineRepo) PowerOff(ctx context.Context, systemID string, params *entity.MachinePowerOffParams) (*Machine, error) {
	return &Machine{
		Machine: &entity.Machine{
			SystemID: systemID,
		},
		LastCommissioned: time.Now(),
	}, nil
}

// Get is a mock implementation to satisfy the MachineRepo interface.
func (m *storageMockMachineRepo) Get(ctx context.Context, systemID string) (*Machine, error) {
	return &Machine{
		Machine: &entity.Machine{
			SystemID: systemID,
		},
		LastCommissioned: time.Now(),
	}, nil
}

// Commission is a mock implementation to satisfy the MachineRepo interface.
func (m *storageMockMachineRepo) Commission(ctx context.Context, systemID string, params *entity.MachineCommissionParams) (*Machine, error) {
	return &Machine{
		Machine: &entity.Machine{
			SystemID: systemID,
		},
		LastCommissioned: time.Now(),
	}, nil
}

// Release is a mock implementation to satisfy the MachineRepo interface.
func (m *storageMockMachineRepo) Release(ctx context.Context, systemID string, params *entity.MachineReleaseParams) (*Machine, error) {
	return &Machine{
		Machine: &entity.Machine{
			SystemID: systemID,
		},
		LastCommissioned: time.Now(),
	}, nil
}

func TestStorageUseCase_ListMONs(t *testing.T) {
	uc := NewStorageUseCase(&clusterMockActionRepo{}, &storageMockFacilityRepo{}, &mockCephClusterRepo{}, nil, nil, nil, &storageMockMachineRepo{})
	mons, err := uc.ListMONs(context.Background(), "uuid", "facility")
	assert.NoError(t, err)
	assert.NotEmpty(t, mons)
}

func TestStorageUseCase_ListOSDs(t *testing.T) {
	uc := NewStorageUseCase(&clusterMockActionRepo{}, &storageMockFacilityRepo{}, &mockCephClusterRepo{}, nil, nil, nil, &storageMockMachineRepo{})
	osds, err := uc.ListOSDs(context.Background(), "uuid", "facility")
	assert.NoError(t, err)
	assert.NotEmpty(t, osds)
}

func TestStorageUseCase_DoSMART(t *testing.T) {
	uc := NewStorageUseCase(&clusterMockActionRepo{}, &storageMockFacilityRepo{}, &mockCephClusterRepo{}, nil, nil, nil, nil)
	smart, err := uc.DoSMART(context.Background(), "uuid", "facility", "osd.1")
	assert.NoError(t, err)
	assert.Contains(t, smart, "osd.1")
}

func TestStorageUseCase_ListPools(t *testing.T) {
	uc := NewStorageUseCase(&clusterMockActionRepo{}, &storageMockFacilityRepo{}, &mockCephClusterRepo{}, nil, nil, nil, nil)
	pools, err := uc.ListPools(context.Background(), "uuid", "facility", "")
	assert.NoError(t, err)
	assert.NotEmpty(t, pools)
}

func TestStorageUseCase_CreatePool(t *testing.T) {
	uc := NewStorageUseCase(&clusterMockActionRepo{}, &storageMockFacilityRepo{}, &mockCephClusterRepo{}, nil, nil, nil, nil)
	pool, err := uc.CreatePool(context.Background(), "uuid", "facility", "pool3", "replicated", false, 3, 100, 10, []string{"rbd"})
	assert.NoError(t, err)
	assert.Equal(t, "pool3", pool.Name)
}

func TestStorageUseCase_UpdatePool(t *testing.T) {
	uc := NewStorageUseCase(&clusterMockActionRepo{}, &storageMockFacilityRepo{}, &mockCephClusterRepo{}, nil, nil, nil, nil)
	pool, err := uc.UpdatePool(context.Background(), "uuid", "facility", "pool3", 200, 20)
	assert.NoError(t, err)
	assert.Equal(t, "pool3", pool.Name)
}

func TestStorageUseCase_DeletePool(t *testing.T) {
	uc := NewStorageUseCase(&clusterMockActionRepo{}, &storageMockFacilityRepo{}, &mockCephClusterRepo{}, nil, nil, nil, nil)
	err := uc.DeletePool(context.Background(), "uuid", "facility", "pool3")
	assert.NoError(t, err)
}
