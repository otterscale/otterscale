package core

import (
	"context"
	"testing"
	"time"

	"github.com/juju/juju/api/client/action"
	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/core/base"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"
	"github.com/stretchr/testify/assert"
)

// Mock CephFSRepo
type mockCephFSRepo struct{}

func (m *mockCephFSRepo) ListVolumes(ctx context.Context, config *StorageConfig) ([]Volume, error) {
	return []Volume{{ID: 1, Name: "vol1", CreatedAt: time.Now()}}, nil
}

func (m *mockCephFSRepo) ListSubvolumes(ctx context.Context, config *StorageConfig, volume, group string) ([]Subvolume, error) {
	return []Subvolume{{Name: "sub1", Path: "/vol1/sub1"}}, nil
}

func (m *mockCephFSRepo) GetSubvolume(ctx context.Context, config *StorageConfig, volume, subvolume, group string) (*Subvolume, error) {
	return &Subvolume{Name: subvolume, Path: "/vol1/" + subvolume}, nil
}

func (m *mockCephFSRepo) ListSubvolumeSnapshots(ctx context.Context, config *StorageConfig, volume, subvolume, group string) ([]SubvolumeSnapshot, error) {
	return []SubvolumeSnapshot{{Name: "snap1"}}, nil
}

func (m *mockCephFSRepo) GetSubvolumeSnapshot(ctx context.Context, config *StorageConfig, volume, subvolume, group, snapshot string) (*SubvolumeSnapshot, error) {
	return &SubvolumeSnapshot{Name: snapshot}, nil
}

func (m *mockCephFSRepo) CreateSubvolume(ctx context.Context, config *StorageConfig, volume, subvolume, group string, size uint64) error {
	return nil
}

func (m *mockCephFSRepo) ResizeSubvolume(ctx context.Context, config *StorageConfig, volume, subvolume, group string, size uint64) error {
	return nil
}

func (m *mockCephFSRepo) DeleteSubvolume(ctx context.Context, config *StorageConfig, volume, subvolume, group string) error {
	return nil
}

func (m *mockCephFSRepo) CreateSubvolumeSnapshot(ctx context.Context, config *StorageConfig, volume, subvolume, group, snapshot string) error {
	return nil
}

func (m *mockCephFSRepo) DeleteSubvolumeSnapshot(ctx context.Context, config *StorageConfig, volume, subvolume, group, snapshot string) error {
	return nil
}

func (m *mockCephFSRepo) ListSubvolumeGroups(ctx context.Context, config *StorageConfig, volume string) ([]SubvolumeGroup, error) {
	return []SubvolumeGroup{{Name: "group1"}}, nil
}

func (m *mockCephFSRepo) GetSubvolumeGroup(ctx context.Context, config *StorageConfig, volume, group string) (*SubvolumeGroup, error) {
	return &SubvolumeGroup{Name: group}, nil
}

func (m *mockCephFSRepo) CreateSubvolumeGroup(ctx context.Context, config *StorageConfig, volume, group string, size uint64) error {
	return nil
}

func (m *mockCephFSRepo) ResizeSubvolumeGroup(ctx context.Context, config *StorageConfig, volume, group string, size uint64) error {
	return nil
}

func (m *mockCephFSRepo) DeleteSubvolumeGroup(ctx context.Context, config *StorageConfig, volume, group string) error {
	return nil
}

func (m *mockCephFSRepo) ListPathToExportClients(ctx context.Context, config *StorageConfig, pool string) (map[string][]string, error) {
	return map[string][]string{"/vol1/sub1": {"10.0.0.1"}}, nil
}

// Mock FacilityRepo
type fsMockFacilityRepo struct{}

func (m *fsMockFacilityRepo) GetLeader(ctx context.Context, uuid, name string) (string, error) {
	return "leader", nil
}

func (m *fsMockFacilityRepo) GetConfig(ctx context.Context, uuid, name string) (map[string]any, error) {
	return map[string]any{"vip": map[string]any{"value": "10.0.0.1"}}, nil
}

// AddUnits is a mock implementation to satisfy the FacilityRepo interface.
func (m *fsMockFacilityRepo) AddUnits(ctx context.Context, uuid, name string, count int, placements []instance.Placement) ([]string, error) {
	return []string{}, nil
}

// Consume is a mock implementation to satisfy the FacilityRepo interface.
func (m *fsMockFacilityRepo) Consume(ctx context.Context, uuid string, args *crossmodel.ConsumeApplicationArgs) error {
	return nil
}

// Create is a mock implementation to satisfy the FacilityRepo interface.
func (m *fsMockFacilityRepo) Create(ctx context.Context, uuid, name, app, series, channel string, numUnits, expose int, base *base.Base, placements []instance.Placement, cons *constraints.Value, force bool) (*application.DeployInfo, error) {
	return &application.DeployInfo{}, nil
}

// CreateRelation is a mock implementation to satisfy the FacilityRepo interface.
func (m *fsMockFacilityRepo) CreateRelation(ctx context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error) {
	return &params.AddRelationResults{}, nil
}

// Delete is a mock implementation to satisfy the FacilityRepo interface.
func (m *fsMockFacilityRepo) Delete(ctx context.Context, uuid, name string, force, destroyStorage bool) error {
	return nil
}

// DeleteRelation is a mock implementation to satisfy the FacilityRepo interface.
func (m *fsMockFacilityRepo) DeleteRelation(ctx context.Context, uuid string, relationID int) error {
	return nil
}

// Expose is a mock implementation to satisfy the FacilityRepo interface.
func (m *fsMockFacilityRepo) Expose(ctx context.Context, uuid, name string, endpoints map[string]params.ExposedEndpoint) error {
	return nil
}

// GetUnitInfo is a mock implementation to satisfy the FacilityRepo interface.
func (m *fsMockFacilityRepo) GetUnitInfo(ctx context.Context, uuid, name string) (*application.UnitInfo, error) {
	return &application.UnitInfo{}, nil
}

// ResolveUnitErrors is a mock implementation to satisfy the FacilityRepo interface.
func (m *fsMockFacilityRepo) ResolveUnitErrors(ctx context.Context, uuid string, unitNames []string) error {
	return nil
}

// Update is a mock implementation to satisfy the FacilityRepo interface.
func (m *fsMockFacilityRepo) Update(ctx context.Context, uuid, name string, config string) error {
	return nil
}

// Mock ActionRepo
type fsMockActionRepo struct{}

func (m *fsMockActionRepo) RunCommand(ctx context.Context, uuid, leader, command string) (string, error) {
	switch command {
	case "ceph config generate-minimal-conf && ceph auth get client.admin":
		return "action-id-ceph-config", nil
	case "radosgw-admin user list":
		return "action-id-user-list", nil
	case "radosgw-admin user info --uid=otterscale --format=json":
		return "action-id-user-info", nil
	default:
		return "action-id-ceph-config", nil
	}
}

func (m *fsMockActionRepo) RunAction(ctx context.Context, uuid, leader, action string, params map[string]any) (string, error) {
	return "action-id", nil
}

func (m *fsMockActionRepo) GetResult(ctx context.Context, uuid, id string) (*action.ActionResult, error) {
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
		return &action.ActionResult{Status: "completed"}, nil
	}
}

// List is a mock implementation to satisfy the ActionRepo interface.
func (m *fsMockActionRepo) List(ctx context.Context, uuid, leader string) (map[string]ActionSpec, error) {
	return map[string]ActionSpec{
		"mock-action": {
			Description: "A mock action",
		},
	}, nil
}

func TestStorageUseCase_ListVolumes(t *testing.T) {
	uc := NewStorageUseCase(&fsMockActionRepo{}, &fsMockFacilityRepo{}, nil, nil, &mockCephFSRepo{}, nil, nil)
	vols, err := uc.ListVolumes(context.Background(), "uuid", "facility")
	assert.NoError(t, err)
	assert.NotEmpty(t, vols)
}

func TestStorageUseCase_ListSubvolumes(t *testing.T) {
	uc := NewStorageUseCase(&fsMockActionRepo{}, &fsMockFacilityRepo{}, nil, nil, &mockCephFSRepo{}, nil, nil)
	subs, err := uc.ListSubvolumes(context.Background(), "uuid", "facility", "vol1", "")
	assert.NoError(t, err)
	assert.NotEmpty(t, subs)
	assert.NotNil(t, subs[0].Export)
}

func TestStorageUseCase_CreateSubvolume(t *testing.T) {
	uc := NewStorageUseCase(&fsMockActionRepo{}, &fsMockFacilityRepo{}, nil, nil, &mockCephFSRepo{}, nil, nil)
	sub, err := uc.CreateSubvolume(context.Background(), "uuid", "facility", "vol1", "sub1", "", 1024, false)
	assert.NoError(t, err)
	assert.NotNil(t, sub)
}

func TestStorageUseCase_UpdateSubvolume(t *testing.T) {
	uc := NewStorageUseCase(&fsMockActionRepo{}, &fsMockFacilityRepo{}, nil, nil, &mockCephFSRepo{}, nil, nil)
	sub, err := uc.UpdateSubvolume(context.Background(), "uuid", "facility", "vol1", "sub1", "", 2048)
	assert.NoError(t, err)
	assert.NotNil(t, sub)
}

func TestStorageUseCase_DeleteSubvolume(t *testing.T) {
	uc := NewStorageUseCase(&fsMockActionRepo{}, &fsMockFacilityRepo{}, nil, nil, &mockCephFSRepo{}, nil, nil)
	err := uc.DeleteSubvolume(context.Background(), "uuid", "facility", "vol1", "sub1", "")
	assert.NoError(t, err)
}

func TestStorageUseCase_GrantSubvolumeClient(t *testing.T) {
	uc := NewStorageUseCase(&fsMockActionRepo{}, &fsMockFacilityRepo{}, nil, nil, &mockCephFSRepo{}, nil, nil)
	err := uc.GrantSubvolumeClient(context.Background(), "uuid", "facility", "sub1", "10.0.0.2")
	assert.NoError(t, err)
}

func TestStorageUseCase_RevokeSubvolumeClient(t *testing.T) {
	uc := NewStorageUseCase(&fsMockActionRepo{}, &fsMockFacilityRepo{}, nil, nil, &mockCephFSRepo{}, nil, nil)
	err := uc.RevokeSubvolumeClient(context.Background(), "uuid", "facility", "sub1", "10.0.0.2")
	assert.NoError(t, err)
}

func TestStorageUseCase_CreateSubvolumeSnapshot(t *testing.T) {
	uc := NewStorageUseCase(&fsMockActionRepo{}, &fsMockFacilityRepo{}, nil, nil, &mockCephFSRepo{}, nil, nil)
	snap, err := uc.CreateSubvolumeSnapshot(context.Background(), "uuid", "facility", "vol1", "sub1", "", "snap1")
	assert.NoError(t, err)
	assert.NotNil(t, snap)
}

func TestStorageUseCase_DeleteSubvolumeSnapshot(t *testing.T) {
	uc := NewStorageUseCase(&fsMockActionRepo{}, &fsMockFacilityRepo{}, nil, nil, &mockCephFSRepo{}, nil, nil)
	err := uc.DeleteSubvolumeSnapshot(context.Background(), "uuid", "facility", "vol1", "sub1", "", "snap1")
	assert.NoError(t, err)
}

func TestStorageUseCase_ListSubvolumeGroups(t *testing.T) {
	uc := NewStorageUseCase(&fsMockActionRepo{}, &fsMockFacilityRepo{}, nil, nil, &mockCephFSRepo{}, nil, nil)
	grps, err := uc.ListSubvolumeGroups(context.Background(), "uuid", "facility", "vol1")
	assert.NoError(t, err)
	assert.NotEmpty(t, grps)
}

func TestStorageUseCase_CreateSubvolumeGroup(t *testing.T) {
	uc := NewStorageUseCase(&fsMockActionRepo{}, &fsMockFacilityRepo{}, nil, nil, &mockCephFSRepo{}, nil, nil)
	grp, err := uc.CreateSubvolumeGroup(context.Background(), "uuid", "facility", "vol1", "group1", 1024)
	assert.NoError(t, err)
	assert.NotNil(t, grp)
}

func TestStorageUseCase_UpdateSubvolumeGroup(t *testing.T) {
	uc := NewStorageUseCase(&fsMockActionRepo{}, &fsMockFacilityRepo{}, nil, nil, &mockCephFSRepo{}, nil, nil)
	grp, err := uc.UpdateSubvolumeGroup(context.Background(), "uuid", "facility", "vol1", "group1", 2048)
	assert.NoError(t, err)
	assert.NotNil(t, grp)
}

func TestStorageUseCase_DeleteSubvolumeGroup(t *testing.T) {
	uc := NewStorageUseCase(&fsMockActionRepo{}, &fsMockFacilityRepo{}, nil, nil, &mockCephFSRepo{}, nil, nil)
	err := uc.DeleteSubvolumeGroup(context.Background(), "uuid", "facility", "vol1", "group1")
	assert.NoError(t, err)
}
