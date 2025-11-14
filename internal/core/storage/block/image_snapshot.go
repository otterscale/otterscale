package block

import (
	"context"
)

type ImageSnapshot struct {
	Name      string
	Quota     uint64
	Used      uint64
	Protected bool
}

// Note: Ceph create and update operations only return error status.
type ImageSnapshotRepo interface {
	Create(ctx context.Context, scope string, pool, image, snapshot string) error
	Delete(ctx context.Context, scope string, pool, image, snapshot string) error
	Rollback(ctx context.Context, scope string, pool, image, snapshot string) error
	Protect(ctx context.Context, scope string, pool, image, snapshot string) error
	Unprotect(ctx context.Context, scope string, pool, image, snapshot string) error
}

func (uc *BlockUseCase) CreateImageSnapshot(ctx context.Context, scope, pool, image, snapshot string) (*ImageSnapshot, error) {
	if err := uc.imageSnapshot.Create(ctx, scope, pool, image, snapshot); err != nil {
		return nil, err
	}

	return &ImageSnapshot{
		Name: snapshot,
	}, nil
}

func (uc *BlockUseCase) DeleteImageSnapshot(ctx context.Context, scope, pool, image, snapshot string) error {
	return uc.imageSnapshot.Delete(ctx, scope, pool, image, snapshot)
}

func (uc *BlockUseCase) RollbackImageSnapshot(ctx context.Context, scope, pool, image, snapshot string) error {
	return uc.imageSnapshot.Rollback(ctx, scope, pool, image, snapshot)
}

func (uc *BlockUseCase) ProtectImageSnapshot(ctx context.Context, scope, pool, image, snapshot string) error {
	return uc.imageSnapshot.Protect(ctx, scope, pool, image, snapshot)
}

func (uc *BlockUseCase) UnprotectImageSnapshot(ctx context.Context, scope, pool, image, snapshot string) error {
	return uc.imageSnapshot.Unprotect(ctx, scope, pool, image, snapshot)
}
