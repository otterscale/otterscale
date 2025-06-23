package core

import (
	"context"
	"math"

	"github.com/ceph/go-ceph/rbd"
)

type RBDImageSnapshot struct {
	Name string
}

type RBDImage struct {
	Name            string
	Size            uint64
	Features        uint64
	ImageMirrorMode string
	MirrorImageInfo string
	Snapshots       []RBDImageSnapshot
}

type CephRBDRepo interface {
	ListImages(ctx context.Context, config *StorageConfig, poolName string) ([]RBDImage, error)
	GetImage(ctx context.Context, config *StorageConfig, poolName, imageName string) (*RBDImage, error)
	CreateImage(ctx context.Context, config *StorageConfig, poolName, imageName string, order int, stripeUnit, stripeCount, size, features uint64) (*RBDImage, error)
	UpdateImageSize(ctx context.Context, config *StorageConfig, poolName, imageName string, size uint64) error
	DeleteImage(ctx context.Context, config *StorageConfig, poolName, imageName string) error
	CreateImageSnapshot(ctx context.Context, config *StorageConfig, poolName, imageName, snapshotName string) error
	DeleteImageSnapshot(ctx context.Context, config *StorageConfig, poolName, imageName, snapshotName string) error
	RollbackImageSnapshot(ctx context.Context, config *StorageConfig, poolName, imageName, snapshotName string) error
	ProtectImageSnapshot(ctx context.Context, config *StorageConfig, poolName, imageName, snapshotName string) error
	UnprotectImageSnapshot(ctx context.Context, config *StorageConfig, poolName, imageName, snapshotName string) error
}

func (uc *StorageUseCase) ListImages(ctx context.Context, uuid, facility string) ([]RBDImage, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}

	pools, err := uc.cluster.ListPoolsByApplication(ctx, config, "rbd")
	if err != nil {
		return nil, err
	}

	images := []RBDImage{}
	for i := range pools {
		imgs, err := uc.rbd.ListImages(ctx, config, pools[i].Name)
		if err != nil {
			return nil, err
		}
		images = append(images, imgs...)
	}
	return images, nil
}

func (uc *StorageUseCase) CreateImage(ctx context.Context, uuid, facility, pool, image string, objectSizeBytes, stripeUnitBytes, stripeCount, size uint64, layering, exclusiveLock, objectMap, fastDiff, deepFlatten bool) (*RBDImage, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}

	order := int(math.Round(math.Log(float64(objectSizeBytes))*100) / 100)
	features := uc.convertToRBDImageFeatures(layering, exclusiveLock, objectMap, fastDiff, deepFlatten)
	return uc.rbd.CreateImage(ctx, config, pool, image, order, stripeUnitBytes, stripeCount, size, features)
}

func (uc *StorageUseCase) UpdateImage(ctx context.Context, uuid, facility, pool, image string, size uint64) (*RBDImage, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	if err := uc.rbd.UpdateImageSize(ctx, config, pool, image, size); err != nil {
		return nil, err
	}
	return uc.rbd.GetImage(ctx, config, pool, image)
}

func (uc *StorageUseCase) DeleteImage(ctx context.Context, uuid, facility, pool, image string) error {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return err
	}
	return uc.rbd.DeleteImage(ctx, config, pool, image)
}

func (uc *StorageUseCase) CreateImageSnapshot(ctx context.Context, uuid, facility, pool, image, snapshot string) (*RBDImageSnapshot, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	if err := uc.rbd.CreateImageSnapshot(ctx, config, pool, image, snapshot); err != nil {
		return nil, err
	}
	return &RBDImageSnapshot{
		Name: snapshot,
	}, nil
}

func (uc *StorageUseCase) DeleteImageSnapshot(ctx context.Context, uuid, facility, pool, image, snapshot string) error {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return err
	}
	return uc.rbd.DeleteImageSnapshot(ctx, config, pool, image, snapshot)
}

func (uc *StorageUseCase) RollbackImageSnapshot(ctx context.Context, uuid, facility, pool, image, snapshot string) error {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return err
	}
	return uc.rbd.RollbackImageSnapshot(ctx, config, pool, image, snapshot)
}

func (uc *StorageUseCase) ProtectImageSnapshot(ctx context.Context, uuid, facility, pool, image, snapshot string) error {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return err
	}
	return uc.rbd.ProtectImageSnapshot(ctx, config, pool, image, snapshot)
}

func (uc *StorageUseCase) UnprotectImageSnapshot(ctx context.Context, uuid, facility, pool, image, snapshot string) error {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return err
	}
	return uc.rbd.UnprotectImageSnapshot(ctx, config, pool, image, snapshot)
}

func (uc *StorageUseCase) convertToRBDImageFeatures(layering, exclusiveLock, objectMap, fastDiff, deepFlatten bool) uint64 {
	var fs uint64
	if layering {
		fs |= rbd.FeatureLayering
	}
	if exclusiveLock {
		fs |= rbd.FeatureExclusiveLock
	}
	if objectMap {
		fs |= rbd.FeatureObjectMap
	}
	if fastDiff {
		fs |= rbd.FeatureFastDiff
	}
	if deepFlatten {
		fs |= rbd.FeatureDeepFlatten
	}
	return fs
}
