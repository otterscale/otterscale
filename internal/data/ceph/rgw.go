package ceph

import (
	"context"

	"github.com/openhdc/otterscale/internal/core"
)

type rgw struct {
	ceph *Ceph
}

func NewRGW(ceph *Ceph) core.CephRGWRepo {
	return &rgw{
		ceph: ceph,
	}
}

var _ core.CephRGWRepo = (*rgw)(nil)

func (r *rgw) ListBuckets(ctx context.Context, config *core.StorageConfig) ([]core.RGWBucket, error) {
	client, err := r.ceph.client(config)
	if err != nil {
		return nil, err
	}
	return client.ListBucketsWithStat(ctx)
}

func (r *rgw) ListRoles(ctx context.Context, config *core.StorageConfig) ([]core.RGWRole, error) {
	return nil, nil
}

func (r *rgw) ListUsers(ctx context.Context, config *core.StorageConfig) ([]core.RGWUser, error) {
	client, err := r.ceph.client(config)
	if err != nil {
		return nil, err
	}
	users := []core.RGWUser{}
	userNames, err := client.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	if userNames != nil {
		for _, userName := range *userNames {
			users = append(users, core.RGWUser{
				Name: userName,
			})
		}
	}
	return users, nil
}

func (r *rgw) ListAccessKeys(ctx context.Context, config *core.StorageConfig) ([]core.RGWAccessKey, error) {
	return nil, nil
}
