package ceph

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/otterscale/otterscale/internal/core/storage"
)

// Note: Ceph API do not support context.
type poolRepo struct {
	ceph *Ceph
}

func NewPoolRepo(ceph *Ceph) storage.PoolRepo {
	return &poolRepo{
		ceph: ceph,
	}
}

var _ storage.PoolRepo = (*poolRepo)(nil)

func (r *poolRepo) List(_ context.Context, scope string, application storage.PoolApplication) ([]storage.Pool, error) {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return nil, err
	}

	osdDump, err := dumpOSD(conn)
	if err != nil {
		return nil, err
	}

	pgDump, err := dumpPG(conn)
	if err != nil {
		return nil, err
	}

	df, err := dfAll(conn)
	if err != nil {
		return nil, err
	}

	pools := r.toPools(osdDump, pgDump, df)

	if application == storage.PoolApplicationUnspecified {
		return pools, nil
	}

	return slices.DeleteFunc(pools, func(p storage.Pool) bool {
		return !slices.Contains(p.Applications, application)
	}), nil
}

func (r *poolRepo) Get(ctx context.Context, scope, pool string) (*storage.Pool, error) {
	pools, err := r.List(ctx, scope, storage.PoolApplicationUnspecified)
	if err != nil {
		return nil, err
	}

	for i := range pools {
		if pools[i].Name == pool {
			return &pools[i], nil
		}
	}

	return nil, fmt.Errorf("storage pool %q not found", pool)
}

func (r *poolRepo) Create(_ context.Context, scope, pool string, poolType storage.PoolType) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	return createOSDPool(conn, pool, poolType)
}

func (r *poolRepo) Delete(_ context.Context, scope, pool string) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	return deleteOSDPool(conn, pool)
}

func (r *poolRepo) Enable(_ context.Context, scope, pool string, application storage.PoolApplication) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	return enableOSDPoolApplication(conn, pool, application.String())
}

func (r *poolRepo) GetParameter(_ context.Context, scope, pool, key string) (string, error) {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return "", err
	}

	return getOSDPool(conn, pool, key)
}

func (r *poolRepo) SetParameter(_ context.Context, scope, pool, key, value string) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	return setOSDPool(conn, pool, key, value)
}

func (r *poolRepo) GetQuota(_ context.Context, scope, pool string) (maxBytes, maxObjects uint64, err error) {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return 0, 0, err
	}

	quota, err := getOSDPoolQuota(conn, pool)
	if err != nil {
		return 0, 0, err
	}

	return quota.QuotaMaxBytes, quota.QuotaMaxObjects, nil
}

func (r *poolRepo) SetQuota(_ context.Context, scope, pool string, maxBytes, maxObjects uint64) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	if err := setOSDPoolQuota(conn, pool, "max_bytes", maxBytes); err != nil {
		return err
	}

	return setOSDPoolQuota(conn, pool, "max_objects", maxObjects)
}

func (r *poolRepo) GetECProfile(_ context.Context, scope, name string) (k, m string, err error) {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return "", "", err
	}

	profile, err := getECProfile(conn, name)
	if err != nil {
		return "", "", err
	}

	return profile.K, profile.M, nil
}

func (r *poolRepo) toPools(d *osdDump, pd *pgDump, df *df) []storage.Pool {
	ret := []storage.Pool{}

	for i := range d.Pools {
		pool := storage.Pool{
			ID:                  d.Pools[i].ID,
			Name:                d.Pools[i].Name,
			Updating:            d.Pools[i].PgNum+d.Pools[i].PgPlacementNum != d.Pools[i].PgNumTarget+d.Pools[i].PgPlacementNumTarget,
			ReplicatedSize:      d.Pools[i].Size,
			PlacementGroupCount: d.Pools[i].PgNum,
			PlacementGroupState: map[string]int64{},
			CreatedAt:           d.Pools[i].CreateTime.Time,
		}

		switch d.Pools[i].Type {
		case 1:
			pool.Type = storage.PoolTypeReplicated

		case 3:
			pool.Type = storage.PoolTypeErasure
		}

		for j := range df.Pools {
			if d.Pools[i].ID != df.Pools[j].ID {
				continue
			}

			pool.UsedBytes = df.Pools[j].Stats.UsedBytes
			pool.UsedObjects = df.Pools[j].Stats.UsedObjects
			pool.MaxBytes = df.Pools[j].Stats.MaxBytes
		}

		for j := range pd.PGMap.PGStats {
			id := strings.Split(pd.PGMap.PGStats[j].ID, ".")[0]

			if strconv.FormatInt(d.Pools[i].ID, 10) != id {
				continue
			}

			state := pd.PGMap.PGStats[j].State
			pool.PlacementGroupState[state]++
		}

		for app := range d.Pools[i].ApplicationMetadata {
			switch app {
			case storage.PoolApplicationBlock.String():
				pool.Applications = append(pool.Applications, storage.PoolApplicationBlock)

			case storage.PoolApplicationFile.String():
				pool.Applications = append(pool.Applications, storage.PoolApplicationFile)

			case storage.PoolApplicationObject.String():
				pool.Applications = append(pool.Applications, storage.PoolApplicationObject)
			}
		}

		ret = append(ret, pool)
	}

	return ret
}
