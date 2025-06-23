package core

import (
	"context"
)

type Snapshot struct {
	Name string
}

type SnapshotSchedule struct {
	Name string
}

type CephSnapshotRepo interface {
	List(ctx context.Context, config *StorageConfig) ([]Snapshot, error)
}

type CephSnapshotScheduleRepo interface {
	List(ctx context.Context, config *StorageConfig) ([]SnapshotSchedule, error)
}

// func (uc *StorageUseCase) ListSnapshots(ctx context.Context, uuid, facility string) ([]Snapshot, error) {
// 	config, err := uc.config(ctx, uuid, facility)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return uc.pool.List(ctx, config)
// }

// func (uc *StorageUseCase) ListSnapshotSchedules(ctx context.Context, uuid, facility string) ([]SnapshotSchedule, error) {
// 	config, err := uc.config(ctx, uuid, facility)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return uc.pool.List(ctx, config)
// }
