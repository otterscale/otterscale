package ceph

import (
	"context"

	cephrbd "github.com/ceph/go-ceph/rbd"

	"github.com/otterscale/otterscale/internal/core/storage/block"
)

type imageSnapshotRepo struct {
	ceph *Ceph
}

func NewImageSnapshotRepo(ceph *Ceph) block.ImageSnapshotRepo {
	return &imageSnapshotRepo{
		ceph: ceph,
	}
}

var _ block.ImageSnapshotRepo = (*imageSnapshotRepo)(nil)

func (r *imageSnapshotRepo) Create(_ context.Context, scope, pool, image, snapshot string) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		return err
	}
	defer ioctx.Destroy()

	img, err := cephrbd.OpenImage(ioctx, image, cephrbd.NoSnapshot)
	if err != nil {
		return err
	}
	defer img.Close()

	_, err = img.CreateSnapshot(snapshot)
	return err
}

func (r *imageSnapshotRepo) Delete(_ context.Context, scope, pool, image, snapshot string) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		return err
	}
	defer ioctx.Destroy()

	img, err := cephrbd.OpenImage(ioctx, image, cephrbd.NoSnapshot)
	if err != nil {
		return err
	}
	defer img.Close()

	return img.GetSnapshot(snapshot).Remove()
}

func (r *imageSnapshotRepo) Rollback(_ context.Context, scope, pool, image, snapshot string) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		return err
	}
	defer ioctx.Destroy()

	img, err := cephrbd.OpenImage(ioctx, image, cephrbd.NoSnapshot)
	if err != nil {
		return err
	}
	defer img.Close()

	return img.GetSnapshot(snapshot).Rollback()
}

func (r *imageSnapshotRepo) Protect(_ context.Context, scope, pool, image, snapshot string) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		return err
	}
	defer ioctx.Destroy()

	img, err := cephrbd.OpenImage(ioctx, image, cephrbd.NoSnapshot)
	if err != nil {
		return err
	}
	defer img.Close()

	return img.GetSnapshot(snapshot).Protect()
}

func (r *imageSnapshotRepo) Unprotect(_ context.Context, scope, pool, image, snapshot string) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		return err
	}
	defer ioctx.Destroy()

	img, err := cephrbd.OpenImage(ioctx, image, cephrbd.NoSnapshot)
	if err != nil {
		return err
	}
	defer img.Close()

	return img.GetSnapshot(snapshot).Unprotect()
}
