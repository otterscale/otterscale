package core

import (
	"context"
)

type RBDImage struct {
	Name            string
	Size            uint64
	Features        string
	ImageMirrorMode string
	MirrorImageInfo string
}

type CephRBDRepo interface {
	ListImages(ctx context.Context, config *StorageConfig, pools ...string) ([]RBDImage, error)
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
	poolNames := []string{}
	for i := range pools {
		poolNames = append(poolNames, pools[i].Name)
	}
	return uc.rbd.ListImages(ctx, config, poolNames...)
}
