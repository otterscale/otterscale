package ceph

import (
	"context"

	"github.com/otterscale/otterscale/internal/core/storage/file"
)

// Note: Ceph API do not support context.
type subvolumeSnapshotRepo struct {
	ceph *Ceph
}

func NewSubvolumeSnapshotRepo(ceph *Ceph) file.SubvolumeSnapshotRepo {
	return &subvolumeSnapshotRepo{
		ceph: ceph,
	}
}

var _ file.SubvolumeSnapshotRepo = (*subvolumeSnapshotRepo)(nil)

func (r *subvolumeSnapshotRepo) List(_ context.Context, scope, volume, subvolume, group string) ([]file.SubvolumeSnapshot, error) {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return nil, err
	}

	dumpNames, err := listSubvolumeSnapshots(conn, volume, subvolume, group)
	if err != nil {
		return nil, err
	}

	var snapshots []file.SubvolumeSnapshot

	for _, dumpName := range dumpNames {
		dump, err := getSubvolumeSnapshot(conn, volume, subvolume, group, dumpName.Name)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, *r.toSubvolumeSnapshot(dumpName.Name, dump))
	}

	return snapshots, nil
}

func (r *subvolumeSnapshotRepo) Get(_ context.Context, scope, volume, subvolume, group, snapshot string) (*file.SubvolumeSnapshot, error) {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return nil, err
	}

	dump, err := getSubvolumeSnapshot(conn, volume, subvolume, group, snapshot)
	if err != nil {
		return nil, err
	}

	return r.toSubvolumeSnapshot(snapshot, dump), nil
}

func (r *subvolumeSnapshotRepo) Create(_ context.Context, scope, volume, subvolume, group, snapshot string) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	return createSubvolumeSnapshot(conn, volume, subvolume, group, snapshot)
}

func (r *subvolumeSnapshotRepo) Delete(_ context.Context, scope, volume, subvolume, group, snapshot string) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	return removeSubvolumeSnapshot(conn, volume, subvolume, group, snapshot)
}

func (r *subvolumeSnapshotRepo) toSubvolumeSnapshot(name string, info *subvolumeSnapshotInfo) *file.SubvolumeSnapshot {
	return &file.SubvolumeSnapshot{
		Name:             name,
		HasPendingClones: info.HasPendingClones,
		CreatedAt:        info.CreatedAt.Time,
	}
}
