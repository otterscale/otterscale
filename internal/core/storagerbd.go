package core

import (
	"context"
	"math"

	"github.com/ceph/go-ceph/rbd"
)

type RBDImage struct {
	Name            string
	Size            uint64
	Features        string
	ImageMirrorMode string
	MirrorImageInfo string
}

type CephRBDRepo interface {
	ListImages(ctx context.Context, config *StorageConfig, poolName string) ([]RBDImage, error)
	GetImage(ctx context.Context, config *StorageConfig, poolName, imageName string) (*RBDImage, error)
	CreateImage(ctx context.Context, config *StorageConfig, poolName, imageName string, order int, stripeUnit, stripeCount, size, features uint64) (*RBDImage, error)
	UpdateImageSize(ctx context.Context, config *StorageConfig, poolName, imageName string, size uint64) error
	UpdateImageFeatures(ctx context.Context, config *StorageConfig, poolName, imageName string, features uint64, enabled bool) error
	DeleteImage(ctx context.Context, config *StorageConfig, poolName, imageName string) error
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

func (uc *StorageUseCase) UpdateImage(ctx context.Context, uuid, facility, pool, image string, size uint64, exclusiveLock, objectMap, fastDiff, deepFlatten bool) (*RBDImage, error) {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}

	// the 'layering' feature cannot be updated
	trueFeatures := uc.convertToRBDImageFeatures(false, exclusiveLock, objectMap, fastDiff, deepFlatten)
	if err := uc.rbd.UpdateImageFeatures(ctx, config, pool, image, trueFeatures, true); err != nil {
		return nil, err
	}

	allFeatures := uc.convertToRBDImageFeatures(true, true, true, true, true)
	falseFeatures := allFeatures ^ trueFeatures
	if err := uc.rbd.UpdateImageFeatures(ctx, config, pool, image, falseFeatures, false); err != nil {
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
