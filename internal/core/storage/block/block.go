package block

import "github.com/otterscale/otterscale/internal/core/storage"

type UseCase struct {
	image         ImageRepo
	imageSnapshot ImageSnapshotRepo

	pool storage.PoolRepo
}

func NewUseCase(image ImageRepo, imageSnapshot ImageSnapshotRepo, pool storage.PoolRepo) *UseCase {
	return &UseCase{
		image:         image,
		imageSnapshot: imageSnapshot,
		pool:          pool,
	}
}
