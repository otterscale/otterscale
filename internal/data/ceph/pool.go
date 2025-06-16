package ceph

import (
	"context"

	"github.com/openhdc/otterscale/internal/core"
)

type pool struct {
	ceph *Ceph
}

func NewPool(ceph *Ceph) core.CephPoolRepo {
	return &pool{
		ceph: ceph,
	}
}

var _ core.CephPoolRepo = (*pool)(nil)

func (r *pool) List(ctx context.Context, config *core.StorageConfig) ([]core.Pool, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	names, err := conn.ListPools()
	if err != nil {
		return nil, err
	}
	pools := []core.Pool{}
	for _, name := range names {
		pools = append(pools, core.Pool{
			Name: name,
		})
	}
	return pools, nil
}

func (r *pool) Create(ctx context.Context, config *core.StorageConfig, name string) (*core.Pool, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	if err := conn.MakePool(name); err != nil {
		return nil, err
	}
	return &core.Pool{
		Name: name,
	}, nil
}

func (r *pool) Delete(ctx context.Context, config *core.StorageConfig, name string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return conn.DeletePool(name)
}
