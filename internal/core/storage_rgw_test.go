package core

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/juju/juju/api/client/action"
	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/core/base"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"
	"github.com/stretchr/testify/assert"
)

// Mock CephRGWRepo
type mockCephRGWRepo struct{}

func (m *mockCephRGWRepo) ListBuckets(ctx context.Context, config *StorageConfig) ([]RGWBucket, error) {
	return []RGWBucket{{Bucket: &admin.Bucket{Bucket: "bucket1"}}}, nil
}

func (m *mockCephRGWRepo) GetBucket(ctx context.Context, config *StorageConfig, bucket string) (*RGWBucket, error) {
	return &RGWBucket{Bucket: &admin.Bucket{Bucket: bucket}}, nil
}

func (m *mockCephRGWRepo) CreateBucket(ctx context.Context, config *StorageConfig, bucket string, acl types.BucketCannedACL) error {
	return nil
}

func (m *mockCephRGWRepo) UpdateBucketOwner(ctx context.Context, config *StorageConfig, bucket, owner string) error {
	return nil
}

func (m *mockCephRGWRepo) UpdateBucketACL(ctx context.Context, config *StorageConfig, bucket string, acl types.BucketCannedACL) error {
	return nil
}

func (m *mockCephRGWRepo) UpdateBucketPolicy(ctx context.Context, config *StorageConfig, bucket, policy string) error {
	return nil
}

func (m *mockCephRGWRepo) DeleteBucket(ctx context.Context, config *StorageConfig, bucket string) error {
	return nil
}

func (m *mockCephRGWRepo) ListUsers(ctx context.Context, config *StorageConfig) ([]RGWUser, error) {
	return []RGWUser{{ID: "user1"}}, nil
}

func (m *mockCephRGWRepo) CreateUser(ctx context.Context, config *StorageConfig, id, name string, suspended bool) (*RGWUser, error) {
	return &RGWUser{ID: id, DisplayName: name}, nil
}

func (m *mockCephRGWRepo) UpdateUser(ctx context.Context, config *StorageConfig, id, name string, suspended bool) (*RGWUser, error) {
	return &RGWUser{ID: id, DisplayName: name}, nil
}

func (m *mockCephRGWRepo) DeleteUser(ctx context.Context, config *StorageConfig, id string) error {
	return nil
}

func (m *mockCephRGWRepo) CreateUserKey(ctx context.Context, config *StorageConfig, id string) (*RGWUserKey, error) {
	return &admin.UserKeySpec{AccessKey: "ak", SecretKey: "sk"}, nil
}

func (m *mockCephRGWRepo) DeleteUserKey(ctx context.Context, config *StorageConfig, id, accessKey string) error {
	return nil
}

// Mock ActionRepo for storageConfig
type rgwMockActionRepo struct{}

func (m *rgwMockActionRepo) List(ctx context.Context, uuid, leader string) (map[string]ActionSpec, error) {
	// Return a mock map of actions as needed for tests
	return map[string]ActionSpec{
		"mock-action-1": {},
		"mock-action-2": {},
	}, nil
}

func (m *rgwMockActionRepo) RunCommand(ctx context.Context, uuid, leader, command string) (string, error) {
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

func (m *rgwMockActionRepo) RunAction(ctx context.Context, uuid, leader, action string, params map[string]any) (string, error) {
	return "action-id", nil
}

func (m *rgwMockActionRepo) GetResult(ctx context.Context, uuid, id string) (*action.ActionResult, error) {
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

// Mock FacilityRepo for storageConfig
type rgwMockFacilityRepo struct{}

func (m *rgwMockFacilityRepo) GetLeader(ctx context.Context, uuid, name string) (string, error) {
	return "leader", nil
}

func (m *rgwMockFacilityRepo) GetUnitInfo(ctx context.Context, uuid, name string) (*application.UnitInfo, error) {
	return &application.UnitInfo{PublicAddress: "10.0.0.1"}, nil
}

// AddUnits is a mock implementation to satisfy the FacilityRepo interface.
func (m *rgwMockFacilityRepo) AddUnits(ctx context.Context, uuid, name string, count int, placements []instance.Placement) ([]string, error) {
	return []string{}, nil
}

// ...

func (m *rgwMockFacilityRepo) Create(ctx context.Context, uuid, name, channel, series, config string, numUnits, expose int, base *base.Base, placements []instance.Placement, cons *constraints.Value, trust bool) (*application.DeployInfo, error) {
	return &application.DeployInfo{}, nil
}

// Consume is a mock implementation to satisfy the FacilityRepo interface.
func (m *rgwMockFacilityRepo) Consume(ctx context.Context, uuid string, args *crossmodel.ConsumeApplicationArgs) error {
	return nil
}

// CreateRelation is a mock implementation to satisfy the FacilityRepo interface.
func (m *rgwMockFacilityRepo) CreateRelation(ctx context.Context, uuid string, apps []string) (*params.AddRelationResults, error) {
	return &params.AddRelationResults{}, nil
}

// Delete is a mock implementation to satisfy the FacilityRepo interface.
func (m *rgwMockFacilityRepo) Delete(ctx context.Context, uuid, name string, force, destroyStorage bool) error {
	return nil
}

// DeleteRelation is a mock implementation to satisfy the FacilityRepo interface.
func (m *rgwMockFacilityRepo) DeleteRelation(ctx context.Context, uuid string, relationID int) error {
	return nil
}

// Expose is a mock implementation to satisfy the FacilityRepo interface.
func (m *rgwMockFacilityRepo) Expose(ctx context.Context, uuid, name string, expose map[string]params.ExposedEndpoint) error {
	return nil
}

// GetConfig is a mock implementation to satisfy the FacilityRepo interface.
func (m *rgwMockFacilityRepo) GetConfig(ctx context.Context, uuid, name string) (map[string]any, error) {
	return map[string]any{"mock": "config"}, nil
}

// ResolveUnitErrors is a mock implementation to satisfy the FacilityRepo interface.
func (m *rgwMockFacilityRepo) ResolveUnitErrors(ctx context.Context, uuid string, names []string) error {
	return nil
}

// Update is a mock implementation to satisfy the FacilityRepo interface.
func (m *rgwMockFacilityRepo) Update(ctx context.Context, uuid, name string, config string) error {
	return nil
}

func TestStorageUseCase_ListBuckets(t *testing.T) {
	uc := NewStorageUseCase(&rgwMockActionRepo{}, &rgwMockFacilityRepo{}, nil, nil, nil, &mockCephRGWRepo{}, nil)
	buckets, err := uc.ListBuckets(context.Background(), "uuid", "facility")
	assert.NoError(t, err)
	assert.NotEmpty(t, buckets)
}

func TestStorageUseCase_CreateBucket(t *testing.T) {
	uc := NewStorageUseCase(&rgwMockActionRepo{}, &rgwMockFacilityRepo{}, nil, nil, nil, &mockCephRGWRepo{}, nil)
	bucket, err := uc.CreateBucket(context.Background(), "uuid", "facility", "bucket1", "owner1", "policy", types.BucketCannedACLPrivate)
	assert.NoError(t, err)
	assert.NotNil(t, bucket)
}

func TestStorageUseCase_UpdateBucket(t *testing.T) {
	uc := NewStorageUseCase(&rgwMockActionRepo{}, &rgwMockFacilityRepo{}, nil, nil, nil, &mockCephRGWRepo{}, nil)
	bucket, err := uc.UpdateBucket(context.Background(), "uuid", "facility", "bucket1", "owner1", "policy", types.BucketCannedACLPrivate)
	assert.NoError(t, err)
	assert.NotNil(t, bucket)
}

func TestStorageUseCase_DeleteBucket(t *testing.T) {
	uc := NewStorageUseCase(&rgwMockActionRepo{}, &rgwMockFacilityRepo{}, nil, nil, nil, &mockCephRGWRepo{}, nil)
	err := uc.DeleteBucket(context.Background(), "uuid", "facility", "bucket1")
	assert.NoError(t, err)
}

func TestStorageUseCase_ListUsers(t *testing.T) {
	uc := NewStorageUseCase(&rgwMockActionRepo{}, &rgwMockFacilityRepo{}, nil, nil, nil, &mockCephRGWRepo{}, nil)
	users, err := uc.ListUsers(context.Background(), "uuid", "facility")
	assert.NoError(t, err)
	assert.NotEmpty(t, users)
}

func TestStorageUseCase_CreateUser(t *testing.T) {
	uc := NewStorageUseCase(&rgwMockActionRepo{}, &rgwMockFacilityRepo{}, nil, nil, nil, &mockCephRGWRepo{}, nil)
	user, err := uc.CreateUser(context.Background(), "uuid", "facility", "id1", "name1", false)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestStorageUseCase_UpdateUser(t *testing.T) {
	uc := NewStorageUseCase(&rgwMockActionRepo{}, &rgwMockFacilityRepo{}, nil, nil, nil, &mockCephRGWRepo{}, nil)
	user, err := uc.UpdateUser(context.Background(), "uuid", "facility", "id1", "name1", false)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestStorageUseCase_DeleteUser(t *testing.T) {
	uc := NewStorageUseCase(&rgwMockActionRepo{}, &rgwMockFacilityRepo{}, nil, nil, nil, &mockCephRGWRepo{}, nil)
	err := uc.DeleteUser(context.Background(), "uuid", "facility", "id1")
	assert.NoError(t, err)
}

func TestStorageUseCase_CreateUserKey(t *testing.T) {
	uc := NewStorageUseCase(&rgwMockActionRepo{}, &rgwMockFacilityRepo{}, nil, nil, nil, &mockCephRGWRepo{}, nil)
	key, err := uc.CreateUserKey(context.Background(), "uuid", "facility", "id1")
	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.Equal(t, "ak", key.AccessKey)
	assert.Equal(t, "sk", key.SecretKey)
}

func TestStorageUseCase_DeleteUserKey(t *testing.T) {
	uc := NewStorageUseCase(&rgwMockActionRepo{}, &rgwMockFacilityRepo{}, nil, nil, nil, &mockCephRGWRepo{}, nil)
	err := uc.DeleteUserKey(context.Background(), "uuid", "facility", "id1", "ak")
	assert.NoError(t, err)
}
