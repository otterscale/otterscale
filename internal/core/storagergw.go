package core

import (
	"context"

	"github.com/ceph/go-ceph/rgw/admin"
)

type RGWBucket = admin.Bucket

type RGWRole struct {
	Name string
}

type RGWUser struct {
	Name string
}

type RGWAccessKey struct {
	Name string
}

type CephRGWRepo interface {
	ListBuckets(ctx context.Context, config *StorageConfig) ([]RGWBucket, error)
	ListRoles(ctx context.Context, config *StorageConfig) ([]RGWRole, error)
	ListUsers(ctx context.Context, config *StorageConfig) ([]RGWUser, error)
	ListAccessKeys(ctx context.Context, config *StorageConfig) ([]RGWAccessKey, error)
}

func (uc *StorageUseCase) ListBuckets(ctx context.Context, uuid, facility string) ([]RGWBucket, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.rgw.ListBuckets(ctx, config)
}

func (uc *StorageUseCase) ListUsers(ctx context.Context, uuid, facility string) ([]RGWUser, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.rgw.ListUsers(ctx, config)
}
