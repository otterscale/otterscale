package file

import "context"

type Volume struct{}

// Note: Ceph create and update operations only return error status.
type VolumeRepo interface {
	List(ctx context.Context, scope string) ([]Volume, error)
}

func (uc *FileUseCase) ListVolumes(ctx context.Context, scope string) ([]Volume, error) {
	return uc.volume.List(ctx, scope)
}
