package core

import (
	"context"
)

type Volume struct {
	Name string
}

type SubvolumeSnapshot struct {
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
	GetSubvolumeSnapshot(ctx context.Context, config *StorageConfig, volume, subvolume, group, snapshot string) (*SubvolumeSnapshot, error)
	CreateSubvolumeSnapshot(ctx context.Context, config *StorageConfig, volume, subvolume, group, snapshot string) error
	DeleteSubvolumeSnapshot(ctx context.Context, config *StorageConfig, volume, subvolume, group, snapshot string) error
	ListSubvolumeGroups(ctx context.Context, config *StorageConfig, volume string) ([]SubvolumeGroup, error)
	GetSubvolumeGroup(ctx context.Context, config *StorageConfig, volume, group string) (*SubvolumeGroup, error)
	CreateSubvolumeGroup(ctx context.Context, config *StorageConfig, volume, group string, size uint64) error
	ResizeSubvolumeGroup(ctx context.Context, config *StorageConfig, volume, group string, size uint64) error
	DeleteSubvolumeGroup(ctx context.Context, config *StorageConfig, volume, group string) error
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

func (uc *StorageUseCase) CreateSubvolumeSnapshot(ctx context.Context, uuid, facility, volume, subvolume, group, snapshot string) (*SubvolumeSnapshot, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	if err := uc.fs.CreateSubvolumeSnapshot(ctx, config, volume, subvolume, group, snapshot); err != nil {
		return nil, err
	}
	return uc.fs.GetSubvolumeSnapshot(ctx, config, volume, subvolume, group, snapshot)
}

func (uc *StorageUseCase) DeleteSubvolumeSnapshot(ctx context.Context, uuid, facility, volume, subvolume, group, snapshot string) error {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return err
	}
	return uc.fs.DeleteSubvolumeSnapshot(ctx, config, volume, subvolume, group, snapshot)
}

func (uc *StorageUseCase) ListSubvolumeGroups(ctx context.Context, uuid, facility, volume string) ([]SubvolumeGroup, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.fs.ListSubvolumeGroups(ctx, config, volume)
}

func (uc *StorageUseCase) CreateSubvolumeGroup(ctx context.Context, uuid, facility, volume, group string, size uint64) (*SubvolumeGroup, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	if err := uc.fs.CreateSubvolumeGroup(ctx, config, volume, group, size); err != nil {
		return nil, err
	}
	return uc.fs.GetSubvolumeGroup(ctx, config, volume, group)
}

func (uc *StorageUseCase) UpdateSubvolumeGroup(ctx context.Context, uuid, facility, volume, group string, size uint64) (*SubvolumeGroup, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	if err := uc.fs.ResizeSubvolumeGroup(ctx, config, volume, group, size); err != nil {
		return nil, err
	}
	return uc.fs.GetSubvolumeGroup(ctx, config, volume, group)
}

func (uc *StorageUseCase) DeleteSubvolumeGroup(ctx context.Context, uuid, facility, volume, group string) error {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return err
	}
	return uc.fs.DeleteSubvolumeGroup(ctx, config, volume, group)
}
