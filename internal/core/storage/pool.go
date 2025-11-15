package storage

import (
	"context"
	"strconv"
	"time"
)

type Pool struct {
	ID                  int64
	Name                string
	Updating            bool
	Type                string
	ECOverwrites        bool
	DataChunks          uint64
	CodingChunks        uint64
	ReplicatedSize      uint64
	QuotaBytes          uint64
	QuotaObjects        uint64
	UsedBytes           uint64
	UsedObjects         uint64
	MaxBytes            uint64
	PlacementGroupCount uint64
	PlacementGroupState map[string]int64
	CreatedAt           time.Time
	Applications        []string
}

// Note: Ceph create and update operations only return error status.
type PoolRepo interface {
	List(ctx context.Context, scope, application string) ([]Pool, error)
	Create(ctx context.Context, scope string, pool, poolType string) error
	Delete(ctx context.Context, scope string, pool string) error
	Enable(ctx context.Context, scope string, pool, application string) error
	GetParameter(ctx context.Context, scope string, pool, key string) (string, error)
	SetParameter(ctx context.Context, scope string, pool, key, value string) error
	GetQuota(ctx context.Context, scope string, pool string) (maxBytes, maxObjects uint64, err error)
	SetQuota(ctx context.Context, scope string, pool string, maxBytes, maxObjects uint64) error
	GetECProfile(ctx context.Context, scope string, pool string) (k, m string, err error)
}

func (uc *UseCase) ListPools(ctx context.Context, scope, application string) ([]Pool, error) {
	pools, err := uc.pool.List(ctx, scope, application)
	if err != nil {
		return nil, err
	}

	uc.setPoolParameters(ctx, scope, pools)

	return pools, nil
}

func (uc *UseCase) CreatePool(ctx context.Context, scope, pool, poolType string, ecOverwrites bool, replicatedSize, quotaMaxBytes, quotaMaxObjects uint64, applications []string) (*Pool, error) {
	if err := uc.pool.Create(ctx, scope, pool, poolType); err != nil {
		return nil, err
	}

	if poolType == "erasure" && ecOverwrites {
		if err := uc.pool.SetParameter(ctx, scope, pool, "allow_ec_overwrites", "true"); err != nil {
			return nil, err
		}
	}

	if poolType == "replicated" && replicatedSize > 1 {
		if err := uc.pool.SetParameter(ctx, scope, pool, "size", strconv.FormatUint(replicatedSize, 10)); err != nil {
			return nil, err
		}
	}

	for _, app := range applications {
		if err := uc.pool.Enable(ctx, scope, pool, app); err != nil {
			return nil, err
		}
	}

	if err := uc.pool.SetQuota(ctx, scope, pool, quotaMaxBytes, quotaMaxObjects); err != nil {
		return nil, err
	}

	return &Pool{
		Name: pool,
	}, nil
}

func (uc *UseCase) UpdatePool(ctx context.Context, scope, pool string, quotaMaxBytes, quotaMaxObjects uint64) (*Pool, error) {
	if err := uc.pool.SetQuota(ctx, scope, pool, quotaMaxBytes, quotaMaxObjects); err != nil {
		return nil, err
	}

	return &Pool{
		Name: pool,
	}, nil
}

func (uc *UseCase) DeletePool(ctx context.Context, scope, pool string) error {
	return uc.pool.Delete(ctx, scope, pool)
}

func (uc *UseCase) setPoolParameters(ctx context.Context, scope string, pools []Pool) {
	for i := range pools {
		if pools[i].Type == "erasure" {
			ecOverwrites, _ := uc.pool.GetParameter(ctx, scope, pools[i].Name, "allow_ec_overwrites")
			if ecOverwrites == "true" {
				pools[i].ECOverwrites = true
			}

			ecProfile, _ := uc.pool.GetParameter(ctx, scope, pools[i].Name, "erasure_code_profile")
			if ecProfile != "" {
				k, m, _ := uc.pool.GetECProfile(ctx, scope, ecProfile)

				pools[i].DataChunks, _ = strconv.ParseUint(k, 10, 64)
				pools[i].CodingChunks, _ = strconv.ParseUint(m, 10, 64)
			}
		}

		maxBytes, maxObjects, _ := uc.pool.GetQuota(ctx, scope, pools[i].Name)
		pools[i].QuotaBytes = maxBytes
		pools[i].QuotaObjects = maxObjects
	}
}
