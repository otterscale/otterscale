package core

import (
	"context"
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
	ListPools(ctx context.Context, config *StorageConfig) ([]Pool, error)
	ListPoolsByApplication(ctx context.Context, config *StorageConfig, application string) ([]Pool, error)
	// Create(ctx context.Context, config *StorageConfig, name string) (*Pool, error)
	// Delete(ctx context.Context, config *StorageConfig, name string) error
}

func (uc *StorageUseCase) ListPools(ctx context.Context, uuid, facility string) ([]Pool, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.cluster.ListPools(ctx, config)
}

// func (uc *StorageUseCase) CreatePool(ctx context.Context, uuid, facility, name string) (*Pool, error) {
// 	config, err := uc.config(ctx, uuid, facility)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return uc.pool.Create(ctx, config, name)
// }

// func (uc *StorageUseCase) DeletePool(ctx context.Context, uuid, facility, name string) error {
// 	config, err := uc.config(ctx, uuid, facility)
// 	if err != nil {
// 		return err
// 	}
// 	return uc.pool.Delete(ctx, config, name)
// }

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
