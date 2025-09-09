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

// Mock CephRBDRepo
type mockCephRBDRepo struct{}

func (m *mockCephRBDRepo) ListImages(ctx context.Context, config *StorageConfig, pool string) ([]RBDImage, error) {
	return []RBDImage{{Name: "img1", PoolName: pool}}, nil
}

func (m *mockCephRBDRepo) GetImage(ctx context.Context, config *StorageConfig, pool, image string) (*RBDImage, error) {
	return &RBDImage{Name: image, PoolName: pool}, nil
}

func (m *mockCephRBDRepo) CreateImage(ctx context.Context, config *StorageConfig, pool, image string, order int, stripeUnit, stripeCount, size, features uint64) (*RBDImage, error) {
	return &RBDImage{Name: image, PoolName: pool}, nil
}

func (m *mockCephRBDRepo) UpdateImageSize(ctx context.Context, config *StorageConfig, pool, image string, size uint64) error {
	return nil
}

func (m *mockCephRBDRepo) DeleteImage(ctx context.Context, config *StorageConfig, pool, image string) error {
	return nil
}

func (m *mockCephRBDRepo) CreateImageSnapshot(ctx context.Context, config *StorageConfig, pool, image, snapshot string) error {
	return nil
}

func (m *mockCephRBDRepo) DeleteImageSnapshot(ctx context.Context, config *StorageConfig, pool, image, snapshot string) error {
	return nil
}

func (m *mockCephRBDRepo) RollbackImageSnapshot(ctx context.Context, config *StorageConfig, pool, image, snapshot string) error {
	return nil
}

func (m *mockCephRBDRepo) ProtectImageSnapshot(ctx context.Context, config *StorageConfig, pool, image, snapshot string) error {
	return nil
}

func (m *mockCephRBDRepo) UnprotectImageSnapshot(ctx context.Context, config *StorageConfig, pool, image, snapshot string) error {
	return nil
}

// Mock CephClusterRepo
type rbdMockCephClusterRepo struct{}

func (m *rbdMockCephClusterRepo) ListPools(ctx context.Context, config *StorageConfig) ([]Pool, error) {
	return []Pool{{Name: "pool1"}}, nil
}

func (m *rbdMockCephClusterRepo) ListPoolsByApplication(ctx context.Context, config *StorageConfig, application string) ([]Pool, error) {
	return []Pool{{Name: "pool1"}}, nil
}

func (m *rbdMockCephClusterRepo) CreatePool(ctx context.Context, config *StorageConfig, poolName string, application string) error {
	return nil
}

// DeletePool is a mock implementation to satisfy the CephClusterRepo interface.
func (m *rbdMockCephClusterRepo) DeletePool(ctx context.Context, config *StorageConfig, poolName string) error {
	return nil
}

// DoSMART is a mock implementation to satisfy the CephClusterRepo interface.
func (m *rbdMockCephClusterRepo) DoSMART(ctx context.Context, config *StorageConfig, poolName string) (map[string][]string, error) {
	return map[string][]string{
		"mock-device": {"mock-smart-output"},
	}, nil
}

// EnableApplication is a mock implementation to satisfy the CephClusterRepo interface.
func (m *rbdMockCephClusterRepo) EnableApplication(ctx context.Context, config *StorageConfig, poolName, application string) error {
	return nil
}

// GetECProfile is a mock implementation to satisfy the CephClusterRepo interface.
func (m *rbdMockCephClusterRepo) GetECProfile(ctx context.Context, config *StorageConfig, profileName string) (string, string, error) {
	return "profile-name", "profile-value", nil
}

// GetParameter is a mock implementation to satisfy the CephClusterRepo interface.
func (m *rbdMockCephClusterRepo) GetParameter(ctx context.Context, config *StorageConfig, key string, pool string) (string, error) {
	return "mock-value", nil
}

// GetQuota is a mock implementation to satisfy the CephClusterRepo interface.
func (m *rbdMockCephClusterRepo) GetQuota(ctx context.Context, config *StorageConfig, poolName string) (uint64, uint64, error) {
	return 0, 0, nil
}

// SetParameter is a mock implementation to satisfy the CephClusterRepo interface.
func (m *rbdMockCephClusterRepo) SetParameter(ctx context.Context, config *StorageConfig, key, value, pool string) error {
	return nil
}

// SetQuota is a mock implementation to satisfy the CephClusterRepo interface.
func (m *rbdMockCephClusterRepo) SetQuota(ctx context.Context, config *StorageConfig, poolName string, maxBytes, maxObjects uint64) error {
	return nil
}

// ListMONs is a mock implementation to satisfy the CephClusterRepo interface.
func (m *rbdMockCephClusterRepo) ListMONs(ctx context.Context, config *StorageConfig) ([]MON, error) {
	return []MON{{PublicAddress: "1.2.3.4"}}, nil
}

// ListOSDs is a mock implementation to satisfy the CephClusterRepo interface.
func (m *rbdMockCephClusterRepo) ListOSDs(ctx context.Context, config *StorageConfig) ([]OSD, error) {
	return []OSD{
		{ID: 0, Name: "osd.0"},
		{ID: 1, Name: "osd.1"},
	}, nil
}

// Mock ActionRepo for storageConfig
type rbdMockActionRepo struct{}

func (m *rbdMockActionRepo) RunCommand(ctx context.Context, uuid, leader, command string) (string, error) {
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

func (m *rbdMockActionRepo) RunAction(ctx context.Context, uuid, leader, action string, params map[string]any) (string, error) {
	return "action-id", nil
}

func (m *rbdMockActionRepo) GetResult(ctx context.Context, uuid, id string) (*action.ActionResult, error) {
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
func (m *rbdMockActionRepo) List(ctx context.Context, uuid, leader string) (map[string]ActionSpec, error) {
	return map[string]ActionSpec{}, nil
}

// Mock FacilityRepo for storageConfig
type rbdMockFacilityRepo struct{}

func (m *rbdMockFacilityRepo) Update(ctx context.Context, uuid, name string, params string) error {
	return nil
}

func (m *rbdMockFacilityRepo) GetConfig(ctx context.Context, uuid, facility string) (map[string]any, error) {
	return map[string]any{}, nil
}

func (m *rbdMockFacilityRepo) GetLeader(ctx context.Context, uuid, name string) (string, error) {
	return "leader", nil
}

func (m *rbdMockFacilityRepo) GetUnitInfo(ctx context.Context, uuid, name string) (*application.UnitInfo, error) {
	return &application.UnitInfo{PublicAddress: "10.0.0.1"}, nil
}

// AddUnits is a mock implementation to satisfy the FacilityRepo interface.
func (m *rbdMockFacilityRepo) AddUnits(ctx context.Context, uuid, app string, count int, placements []instance.Placement) ([]string, error) {
	return []string{}, nil
}

// Consume is a mock implementation to satisfy the FacilityRepo interface.
func (m *rbdMockFacilityRepo) Consume(ctx context.Context, uuid string, args *crossmodel.ConsumeApplicationArgs) error {
	return nil
}

func (m *rbdMockFacilityRepo) Create(ctx context.Context, uuid, name, series, charm, channel string, numUnits, expose int, base *base.Base, placements []instance.Placement, cons *constraints.Value, force bool) (*application.DeployInfo, error) {
	return &application.DeployInfo{}, nil
}

// Delete is a mock implementation to satisfy the FacilityRepo interface.
func (m *rbdMockFacilityRepo) Delete(ctx context.Context, uuid, name string, force, destroyStorage bool) error {
	return nil
}

// CreateRelation is a mock implementation to satisfy the FacilityRepo interface.
func (m *rbdMockFacilityRepo) CreateRelation(ctx context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error) {
	return &params.AddRelationResults{}, nil
}

// DeleteRelation is a mock implementation to satisfy the FacilityRepo interface.
func (m *rbdMockFacilityRepo) DeleteRelation(ctx context.Context, uuid string, relationID int) error {
	return nil
}

// Expose is a mock implementation to satisfy the FacilityRepo interface.
func (m *rbdMockFacilityRepo) Expose(ctx context.Context, uuid, name string, expose map[string]params.ExposedEndpoint) error {
	return nil
}

// ResolveUnitErrors is a mock implementation to satisfy the FacilityRepo interface.
func (m *rbdMockFacilityRepo) ResolveUnitErrors(ctx context.Context, uuid string, units []string) error {
	return nil
}

func TestStorageUseCase_ListImages(t *testing.T) {
	uc := NewStorageUseCase(&rbdMockActionRepo{}, &rbdMockFacilityRepo{}, &rbdMockCephClusterRepo{}, &mockCephRBDRepo{}, nil, nil, nil)
	images, err := uc.ListImages(context.Background(), "uuid", "facility")
	assert.NoError(t, err)
	assert.NotEmpty(t, images)
}

func TestStorageUseCase_CreateImage(t *testing.T) {
	uc := NewStorageUseCase(&rbdMockActionRepo{}, &rbdMockFacilityRepo{}, &rbdMockCephClusterRepo{}, &mockCephRBDRepo{}, nil, nil, nil)
	img, err := uc.CreateImage(context.Background(), "uuid", "facility", "pool1", "img1", 4096, 4096, 1, 1024, true, false, false, false, false)
	assert.NoError(t, err)
	assert.NotNil(t, img)
}

func TestStorageUseCase_UpdateImage(t *testing.T) {
	uc := NewStorageUseCase(&rbdMockActionRepo{}, &rbdMockFacilityRepo{}, &rbdMockCephClusterRepo{}, &mockCephRBDRepo{}, nil, nil, nil)
	img, err := uc.UpdateImage(context.Background(), "uuid", "facility", "pool1", "img1", 2048)
	assert.NoError(t, err)
	assert.NotNil(t, img)
}

func TestStorageUseCase_DeleteImage(t *testing.T) {
	uc := NewStorageUseCase(&rbdMockActionRepo{}, &rbdMockFacilityRepo{}, &rbdMockCephClusterRepo{}, &mockCephRBDRepo{}, nil, nil, nil)
	err := uc.DeleteImage(context.Background(), "uuid", "facility", "pool1", "img1")
	assert.NoError(t, err)
}

func TestStorageUseCase_CreateImageSnapshot(t *testing.T) {
	uc := NewStorageUseCase(&rbdMockActionRepo{}, &rbdMockFacilityRepo{}, &rbdMockCephClusterRepo{}, &mockCephRBDRepo{}, nil, nil, nil)
	snap, err := uc.CreateImageSnapshot(context.Background(), "uuid", "facility", "pool1", "img1", "snap1")
	assert.NoError(t, err)
	assert.NotNil(t, snap)
	assert.Equal(t, "snap1", snap.Name)
}

func TestStorageUseCase_DeleteImageSnapshot(t *testing.T) {
	uc := NewStorageUseCase(&rbdMockActionRepo{}, &rbdMockFacilityRepo{}, &rbdMockCephClusterRepo{}, &mockCephRBDRepo{}, nil, nil, nil)
	err := uc.DeleteImageSnapshot(context.Background(), "uuid", "facility", "pool1", "img1", "snap1")
	assert.NoError(t, err)
}

func TestStorageUseCase_RollbackImageSnapshot(t *testing.T) {
	uc := NewStorageUseCase(&rbdMockActionRepo{}, &rbdMockFacilityRepo{}, &rbdMockCephClusterRepo{}, &mockCephRBDRepo{}, nil, nil, nil)
	err := uc.RollbackImageSnapshot(context.Background(), "uuid", "facility", "pool1", "img1", "snap1")
	assert.NoError(t, err)
}

func TestStorageUseCase_ProtectImageSnapshot(t *testing.T) {
	uc := NewStorageUseCase(&rbdMockActionRepo{}, &rbdMockFacilityRepo{}, &rbdMockCephClusterRepo{}, &mockCephRBDRepo{}, nil, nil, nil)
	err := uc.ProtectImageSnapshot(context.Background(), "uuid", "facility", "pool1", "img1", "snap1")
	assert.NoError(t, err)
}

func TestStorageUseCase_UnprotectImageSnapshot(t *testing.T) {
	uc := NewStorageUseCase(&rbdMockActionRepo{}, &rbdMockFacilityRepo{}, &rbdMockCephClusterRepo{}, &mockCephRBDRepo{}, nil, nil, nil)
	err := uc.UnprotectImageSnapshot(context.Background(), "uuid", "facility", "pool1", "img1", "snap1")
	assert.NoError(t, err)
}
