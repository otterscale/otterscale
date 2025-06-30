package core

import (
	"context"
	"strconv"
	"time"
)

type MON struct {
	Leader        bool
	Name          string
	Rank          uint64
	PublicAddress string

	Hostname string
}

type OSD struct {
	ID          int64
	Name        string
	Up          bool
	In          bool
	Exists      bool
	DeviceClass string
	Size        uint64
	Used        uint64
	PGCount     uint64

	Hostname string
}

type Pool struct {
	ID                  int64
	Name                string
	Type                string
	ECOverwrites        bool
	DataChunks          uint64
	CodingChunks        uint64
	ReplicatedSize      uint64
	QuotaBytes          uint64
	QuotaObjects        uint64
	UsedBytes           uint64
	UsedObjects         uint64
	PlacementGroupCount uint64
	PlacementGroupState map[string]int64
	CreatedAt           time.Time
	Applications        []string
}

type CephClusterRepo interface {
	ListMONs(ctx context.Context, config *StorageConfig) ([]MON, error)
	ListOSDs(ctx context.Context, config *StorageConfig) ([]OSD, error)
	DoSMART(ctx context.Context, config *StorageConfig, who string) (map[string][]string, error)
	ListPools(ctx context.Context, config *StorageConfig) ([]Pool, error)
	ListPoolsByApplication(ctx context.Context, config *StorageConfig, application string) ([]Pool, error)
	CreatePool(ctx context.Context, config *StorageConfig, pool, poolType string) error
	DeletePool(ctx context.Context, config *StorageConfig, pool string) error
	EnableApplication(ctx context.Context, config *StorageConfig, pool, application string) error
	GetParameter(ctx context.Context, config *StorageConfig, pool, key string) (string, error)
	SetParameter(ctx context.Context, config *StorageConfig, pool, key, value string) error
	SetQuota(ctx context.Context, config *StorageConfig, pool string, maxBytes, maxObjects uint64) error
	GetQuota(ctx context.Context, config *StorageConfig, pool string) (maxBytes, maxObjects uint64, err error)
	GetECProfile(ctx context.Context, config *StorageConfig, name string) (k, m string, err error)
}

func (uc *StorageUseCase) ListMONs(ctx context.Context, uuid, facility string) ([]MON, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.cluster.ListMONs(ctx, config)
}

func (uc *StorageUseCase) ListOSDs(ctx context.Context, uuid, facility string) ([]OSD, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.cluster.ListOSDs(ctx, config)
}

func (uc *StorageUseCase) DoSMART(ctx context.Context, uuid, facility, osd string) (map[string][]string, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.cluster.DoSMART(ctx, config, osd)
}

func (uc *StorageUseCase) ListPools(ctx context.Context, uuid, facility, application string) ([]Pool, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	if application != "" {
		return uc.cluster.ListPoolsByApplication(ctx, config, application)
	}
	pools, err := uc.cluster.ListPools(ctx, config)
	if err != nil {
		return nil, err
	}
	for i := range pools {
		if pools[i].Type == "erasure" {
			ecOverwrites, _ := uc.cluster.GetParameter(ctx, config, pools[i].Name, "allow_ec_overwrites")
			if ecOverwrites == "true" {
				pools[i].ECOverwrites = true
			}
			ecProfile, _ := uc.cluster.GetParameter(ctx, config, pools[i].Name, "erasure_code_profile")
			if ecProfile != "" {
				k, m, _ := uc.cluster.GetECProfile(ctx, config, ecProfile)
				pools[i].DataChunks, _ = strconv.ParseUint(k, 10, 64)
				pools[i].CodingChunks, _ = strconv.ParseUint(m, 10, 64)
			}
		}
		maxBytes, maxObjects, _ := uc.cluster.GetQuota(ctx, config, pools[i].Name)
		pools[i].QuotaBytes = maxBytes
		pools[i].QuotaObjects = maxObjects
	}
	return pools, nil
}

func (uc *StorageUseCase) CreatePool(ctx context.Context, uuid, facility, pool, poolType string, ecOverwrites bool, replicatedSize, quotaMaxBytes, quotaMaxObjects uint64, applications []string) (*Pool, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	if err := uc.cluster.CreatePool(ctx, config, pool, poolType); err != nil {
		return nil, err
	}
	if poolType == "erasure" && ecOverwrites {
		if err := uc.cluster.SetParameter(ctx, config, pool, "allow_ec_overwrites", "true"); err != nil {
			return nil, err
		}
	}
	if poolType == "replicated" && replicatedSize > 0 {
		if err := uc.cluster.SetParameter(ctx, config, pool, "size", strconv.FormatUint(replicatedSize, 10)); err != nil {
			return nil, err
		}
	}
	for _, app := range applications {
		if err := uc.cluster.EnableApplication(ctx, config, pool, app); err != nil {
			return nil, err
		}
	}
	if err := uc.cluster.SetQuota(ctx, config, pool, quotaMaxBytes, quotaMaxObjects); err != nil {
		return nil, err
	}
	return &Pool{
		Name: pool,
	}, nil
}

func (uc *StorageUseCase) UpdatePool(ctx context.Context, uuid, facility, pool string, quotaMaxBytes, quotaMaxObjects uint64) (*Pool, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	if err := uc.cluster.SetQuota(ctx, config, pool, quotaMaxBytes, quotaMaxObjects); err != nil {
		return nil, err
	}
	return &Pool{
		Name: pool,
	}, nil
}

func (uc *StorageUseCase) DeletePool(ctx context.Context, uuid, facility, pool string) error {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return err
	}
	return uc.cluster.DeletePool(ctx, config, pool)
}
