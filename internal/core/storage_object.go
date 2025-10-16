package core

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ceph/go-ceph/rgw/admin"
)

const RGWBucketCannedACLPrivate = types.BucketCannedACLPrivate

type (
	RGWBucketCannedACL = types.BucketCannedACL
	RGWGrant           = types.Grant
	RGWUser            = admin.User
	RGWUserKey         = admin.UserKeySpec
)

type RGWBucket struct {
	*admin.Bucket
	Policy *string
	Grants []types.Grant
}

type CephRGWRepo interface {
	ListBuckets(ctx context.Context, config *CephConfig) ([]RGWBucket, error)
	GetBucket(ctx context.Context, config *CephConfig, bucket string) (*RGWBucket, error)
	CreateBucket(ctx context.Context, config *CephConfig, bucket string, acl types.BucketCannedACL) error
	UpdateBucketOwner(ctx context.Context, config *CephConfig, bucket, owner string) error
	UpdateBucketACL(ctx context.Context, config *CephConfig, bucket string, acl types.BucketCannedACL) error
	UpdateBucketPolicy(ctx context.Context, config *CephConfig, bucket, policy string) error
	DeleteBucket(ctx context.Context, config *CephConfig, bucket string) error
	ListUsers(ctx context.Context, config *CephConfig) ([]RGWUser, error)
	CreateUser(ctx context.Context, config *CephConfig, id, name string, suspended bool) (*RGWUser, error)
	UpdateUser(ctx context.Context, config *CephConfig, id, name string, suspended bool) (*RGWUser, error)
	DeleteUser(ctx context.Context, config *CephConfig, id string) error
	CreateUserKey(ctx context.Context, config *CephConfig, id string) (*RGWUserKey, error)
	DeleteUserKey(ctx context.Context, config *CephConfig, id, accessKey string) error
}

func (uc *StorageUseCase) ListBuckets(ctx context.Context, scope, facility string) ([]RGWBucket, error) {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}
	return uc.cephRGW.ListBuckets(ctx, config)
}

func (uc *StorageUseCase) CreateBucket(ctx context.Context, scope, facility, bucket, owner, policy string, acl types.BucketCannedACL) (*RGWBucket, error) {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}
	if err := uc.cephRGW.CreateBucket(ctx, config, bucket, acl); err != nil {
		return nil, err
	}
	if err := uc.cephRGW.UpdateBucketOwner(ctx, config, bucket, owner); err != nil {
		return nil, err
	}
	if policy != "" {
		if err := uc.cephRGW.UpdateBucketPolicy(ctx, config, bucket, policy); err != nil {
			return nil, err
		}
	}
	return uc.cephRGW.GetBucket(ctx, config, bucket)
}

func (uc *StorageUseCase) UpdateBucket(ctx context.Context, scope, facility, bucket, owner, policy string, acl types.BucketCannedACL) (*RGWBucket, error) {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}
	if err := uc.cephRGW.UpdateBucketACL(ctx, config, bucket, acl); err != nil {
		return nil, err
	}
	if err := uc.cephRGW.UpdateBucketOwner(ctx, config, bucket, owner); err != nil {
		return nil, err
	}
	if policy != "" {
		if err := uc.cephRGW.UpdateBucketPolicy(ctx, config, bucket, policy); err != nil {
			return nil, err
		}
	}
	return uc.cephRGW.GetBucket(ctx, config, bucket)
}

func (uc *StorageUseCase) DeleteBucket(ctx context.Context, scope, facility, bucket string) error {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return err
	}
	return uc.cephRGW.DeleteBucket(ctx, config, bucket)
}

func (uc *StorageUseCase) ListUsers(ctx context.Context, scope, facility string) ([]RGWUser, error) {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}
	return uc.cephRGW.ListUsers(ctx, config)
}

func (uc *StorageUseCase) CreateUser(ctx context.Context, scope, facility, id, name string, suspended bool) (*RGWUser, error) {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}
	return uc.cephRGW.CreateUser(ctx, config, id, name, suspended)
}

func (uc *StorageUseCase) UpdateUser(ctx context.Context, scope, facility, id, name string, suspended bool) (*RGWUser, error) {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}
	return uc.cephRGW.UpdateUser(ctx, config, id, name, suspended)
}

func (uc *StorageUseCase) DeleteUser(ctx context.Context, scope, facility, id string) error {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return err
	}
	return uc.cephRGW.DeleteUser(ctx, config, id)
}

func (uc *StorageUseCase) CreateUserKey(ctx context.Context, scope, facility, id string) (*admin.UserKeySpec, error) {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}
	return uc.cephRGW.CreateUserKey(ctx, config, id)
}

func (uc *StorageUseCase) DeleteUserKey(ctx context.Context, scope, facility, id, accessKey string) error {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return err
	}
	return uc.cephRGW.DeleteUserKey(ctx, config, id, accessKey)
}
