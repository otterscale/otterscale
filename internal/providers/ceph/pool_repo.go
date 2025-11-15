package ceph

import (
	"context"
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

func (r *poolRepo) List(_ context.Context, scope, application string) ([]storage.Pool, error) {
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

	return slices.DeleteFunc(pools, func(p storage.Pool) bool {
		return !slices.Contains(p.Applications, application)
	}), nil
}

func (r *poolRepo) Create(_ context.Context, scope, pool, poolType string) error {
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

func (r *poolRepo) Enable(_ context.Context, scope, pool, application string) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	return enableOSDPoolApplication(conn, pool, application)
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
			pool.Type = "replicated"

		case 3:
			pool.Type = "erasure"
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
			pool.Applications = append(pool.Applications, app)
		}

		ret = append(ret, pool)
	}

	return ret
}
