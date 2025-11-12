package file

import "context"

type SubvolumeGroup struct{}

// Note: Ceph create and update operations only return error status.
type SubvolumeGroupRepo interface {
	List(ctx context.Context, scope string, volume string) ([]SubvolumeGroup, error)
	Get(ctx context.Context, scope string, volume, group string) (*SubvolumeGroup, error)
	Create(ctx context.Context, scope string, volume, group string, size uint64) error
	Resize(ctx context.Context, scope string, volume, group string, size uint64) error
	Delete(ctx context.Context, scope string, volume, group string) error
}

func (uc *FileUseCase) ListSubvolumeGroups(ctx context.Context, scope, volume string) ([]SubvolumeGroup, error) {
	return uc.subvolumeGroup.List(ctx, scope, volume)
}

func (uc *FileUseCase) CreateSubvolumeGroup(ctx context.Context, scope, volume, group string, size uint64) (*SubvolumeGroup, error) {
	if err := uc.subvolumeGroup.Create(ctx, scope, volume, group, size); err != nil {
		return nil, err
	}
	return uc.subvolumeGroup.Get(ctx, scope, volume, group)
}

func (uc *FileUseCase) UpdateSubvolumeGroup(ctx context.Context, scope, volume, group string, size uint64) (*SubvolumeGroup, error) {
	if err := uc.subvolumeGroup.Resize(ctx, scope, volume, group, size); err != nil {
		return nil, err
	}
	return uc.subvolumeGroup.Get(ctx, scope, volume, group)
}

func (uc *FileUseCase) DeleteSubvolumeGroup(ctx context.Context, scope, volume, group string) error {
	return uc.subvolumeGroup.Delete(ctx, scope, volume, group)
}
