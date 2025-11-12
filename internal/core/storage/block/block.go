package block

import (
	"github.com/otterscale/otterscale/internal/core/storage/pool"
)

type BlockUseCase struct {
	image         ImageRepo
	imageSnapshot ImageSnapshotRepo

	pool pool.PoolRepo
}

func NewBlockUseCase(image ImageRepo, imageSnapshot ImageSnapshotRepo, pool pool.PoolRepo) *BlockUseCase {
	return &BlockUseCase{
		image:         image,
		imageSnapshot: imageSnapshot,
		pool:          pool,
	}
}
