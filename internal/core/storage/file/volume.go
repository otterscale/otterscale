package file

import (
	"context"
	"time"
)

type Volume struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}

// Note: Ceph create and update operations only return error status.
type VolumeRepo interface {
	List(ctx context.Context, scope string) ([]Volume, error)
}

func (uc *UseCase) ListVolumes(ctx context.Context, scope string) ([]Volume, error) {
	return uc.volume.List(ctx, scope)
}
