package block

import "github.com/otterscale/otterscale/internal/core/storage"

type BlockUseCase struct {
	image         ImageRepo
	imageSnapshot ImageSnapshotRepo

	pool storage.PoolRepo
}

func NewBlockUseCase(image ImageRepo, imageSnapshot ImageSnapshotRepo, pool storage.PoolRepo) *BlockUseCase {
	return &BlockUseCase{
		image:         image,
		imageSnapshot: imageSnapshot,
		pool:          pool,
	}
}
