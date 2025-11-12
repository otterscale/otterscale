package file

import "context"

type SubvolumeSnapshot struct{}

// Note: Ceph create and update operations only return error status.
type SubvolumeSnapshotRepo interface {
	List(ctx context.Context, scope string, volume, subvolume, group string) ([]SubvolumeSnapshot, error)
	Get(ctx context.Context, scope string, volume, subvolume, group, snapshot string) (*SubvolumeSnapshot, error)
	Create(ctx context.Context, scope string, volume, subvolume, group, snapshot string) error
	Delete(ctx context.Context, scope string, volume, subvolume, group, snapshot string) error
}

func (uc *FileUseCase) CreateSubvolumeSnapshot(ctx context.Context, scope, volume, subvolume, group, snapshot string) (*SubvolumeSnapshot, error) {
	if err := uc.subvolumeSnapshot.Create(ctx, scope, volume, subvolume, group, snapshot); err != nil {
		return nil, err
	}
	return uc.subvolumeSnapshot.Get(ctx, scope, volume, subvolume, group, snapshot)
}

func (uc *FileUseCase) DeleteSubvolumeSnapshot(ctx context.Context, scope, volume, subvolume, group, snapshot string) error {
	return uc.subvolumeSnapshot.Delete(ctx, scope, volume, subvolume, group, snapshot)
}
