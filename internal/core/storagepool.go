package core

import (
	"context"
)

type Pool struct {
	Name string
}

type CephPoolRepo interface {
	List(ctx context.Context, config *StorageConfig) ([]Pool, error)
	Create(ctx context.Context, config *StorageConfig, name string) (*Pool, error)
	Delete(ctx context.Context, config *StorageConfig, name string) error
}

func (uc *StorageUseCase) ListPools(ctx context.Context, uuid, facility string) ([]Pool, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.pool.List(ctx, config)
}

func (uc *StorageUseCase) CreatePool(ctx context.Context, uuid, facility, name string) (*Pool, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.pool.Create(ctx, config, name)
}

func (uc *StorageUseCase) DeletePool(ctx context.Context, uuid, facility, name string) error {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return err
	}
	return uc.pool.Delete(ctx, config, name)
}
