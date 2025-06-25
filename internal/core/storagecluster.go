package core

import (
	"context"
	"strconv"
)

type MON struct {
	Name string
}

type OSD struct {
	Name        string
	DeviceClass string
}

type Pool struct {
	Name         string
	Applications []string
}

type CephClusterRepo interface {
	ListMONs(ctx context.Context, config *StorageConfig) ([]MON, error)
	ListOSDs(ctx context.Context, config *StorageConfig) ([]OSD, error)
	DoSMART(ctx context.Context, config *StorageConfig, who string) (map[string][]string, error)
	ListPools(ctx context.Context, config *StorageConfig) ([]Pool, error)
	ListPoolsByApplication(ctx context.Context, config *StorageConfig, application string) ([]Pool, error)
	CreatePool(ctx context.Context, config *StorageConfig, poolName, poolType string) error
	DeletePool(ctx context.Context, config *StorageConfig, poolName string) error
	EnableApplication(ctx context.Context, config *StorageConfig, poolName, application string) error
	SetParameter(ctx context.Context, config *StorageConfig, poolName, key, value string) error
	SetQuota(ctx context.Context, config *StorageConfig, poolName string, maxBytes, maxObjects uint64) error
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
	return uc.cluster.ListPools(ctx, config)
}

func (uc *StorageUseCase) CreatePool(ctx context.Context, uuid, facility, pool, poolType string, ecOverwrites bool, replicatedSize uint32, quotaMaxBytes, quotaMaxObjects uint64, applications []string) (*Pool, error) {
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
		if err := uc.cluster.SetParameter(ctx, config, pool, "size", strconv.Itoa(int(replicatedSize))); err != nil {
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
