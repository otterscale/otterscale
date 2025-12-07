package block

import (
	"context"
	"math"
	"time"

	"github.com/otterscale/otterscale/internal/core/storage"
)

const (
	ImageFeatureLayering uint64 = 1 << iota
	ImageFeatureStripingV2
	ImageFeatureExclusiveLock
	ImageFeatureObjectMap
	ImageFeatureFastDiff
	ImageFeatureDeepFlatten
)

type Image struct {
	Name                 string
	PoolName             string
	ObjectSize           uint64
	StripeUnit           uint64
	StripeCount          uint64
	Quota                uint64
	Used                 uint64
	ObjectCount          uint64
	FeatureLayering      bool
	FeatureExclusiveLock bool
	FeatureObjectMap     bool
	FeatureFastDiff      bool
	FeatureDeepFlatten   bool
	CreatedAt            time.Time
	Snapshots            []ImageSnapshot
}

// Note: Ceph create and update operations only return error status.
type ImageRepo interface {
	List(ctx context.Context, scope string, pool string) ([]Image, error)
	Get(ctx context.Context, scope string, pool, image string) (*Image, error)
	Create(ctx context.Context, scope string, pool, image string, order int, stripeUnit, stripeCount, size, features uint64) error
	Resize(ctx context.Context, scope string, pool, image string, size uint64) error
	Delete(ctx context.Context, scope string, pool, image string) error
}

func (uc *UseCase) ListImages(ctx context.Context, scope string) ([]Image, error) {
	pools, err := uc.pool.List(ctx, scope, storage.PoolApplicationBlock)
	if err != nil {
		return nil, err
	}

	images := []Image{}

	for i := range pools {
		poolImages, err := uc.image.List(ctx, scope, pools[i].Name)
		if err != nil {
			return nil, err
		}
		images = append(images, poolImages...)
	}

	return images, nil
}

func (uc *UseCase) CreateImage(ctx context.Context, scope, pool, image string, objectSizeBytes, stripeUnitBytes, stripeCount, size uint64, layering, exclusiveLock, objectMap, fastDiff, deepFlatten bool) (*Image, error) {
	order := int(math.Round(math.Log2(float64(objectSizeBytes))))
	features := uc.toImageFeatures(layering, exclusiveLock, objectMap, fastDiff, deepFlatten)

	if err := uc.image.Create(ctx, scope, pool, image, order, stripeUnitBytes, stripeCount, size, features); err != nil {
		return nil, err
	}

	return uc.image.Get(ctx, scope, pool, image)
}

func (uc *UseCase) UpdateImage(ctx context.Context, scope, pool, image string, size uint64) (*Image, error) {
	if err := uc.image.Resize(ctx, scope, pool, image, size); err != nil {
		return nil, err
	}

	return uc.image.Get(ctx, scope, pool, image)
}

func (uc *UseCase) DeleteImage(ctx context.Context, scope, pool, image string) error {
	return uc.image.Delete(ctx, scope, pool, image)
}

func (uc *UseCase) toImageFeatures(layering, exclusiveLock, objectMap, fastDiff, deepFlatten bool) uint64 {
	var features uint64

	if layering {
		features |= ImageFeatureLayering
	}

	if exclusiveLock {
		features |= ImageFeatureExclusiveLock
	}

	if objectMap {
		features |= ImageFeatureObjectMap
	}

	if fastDiff {
		features |= ImageFeatureFastDiff
	}

	if deepFlatten {
		features |= ImageFeatureDeepFlatten
	}

	return features
}
