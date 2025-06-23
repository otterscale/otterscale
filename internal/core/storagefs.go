package core

import (
	"context"
)

type Volume struct {
	Name string
}

type Subvolume struct {
	Name string
}

type SubvolumeGroup struct {
	Name string
}

type CephFSRepo interface {
	ListVolumes(ctx context.Context, config *StorageConfig) ([]Volume, error)
	ListSubvolumes(ctx context.Context, config *StorageConfig, volume, group string) ([]Subvolume, error)
	ListSubvolumeGroups(ctx context.Context, config *StorageConfig, volume string) ([]SubvolumeGroup, error)
}

func (uc *StorageUseCase) ListVolumes(ctx context.Context, uuid, facility string) ([]Volume, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.fs.ListVolumes(ctx, config)
}

func (uc *StorageUseCase) ListSubvolumes(ctx context.Context, uuid, facility, volume, group string) ([]Subvolume, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.fs.ListSubvolumes(ctx, config, volume, group)
}

func (uc *StorageUseCase) ListSubvolumeGroups(ctx context.Context, uuid, facility, volume string) ([]SubvolumeGroup, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.fs.ListSubvolumeGroups(ctx, config, volume)
}
