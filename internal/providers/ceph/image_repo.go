package ceph

import (
	"context"
	"slices"
	"time"

	"github.com/ceph/go-ceph/rados"
	"github.com/ceph/go-ceph/rbd"

	"github.com/otterscale/otterscale/internal/core/storage/block"
)

type imageRepo struct {
	ceph *Ceph
}

func NewImageRepo(ceph *Ceph) block.ImageRepo {
	return &imageRepo{
		ceph: ceph,
	}
}

var _ block.ImageRepo = (*imageRepo)(nil)

func (r *imageRepo) List(_ context.Context, scope, pool string) ([]block.Image, error) {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return nil, err
	}

	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		return nil, err
	}
	defer ioctx.Destroy()

	imgNames, err := rbd.GetImageNames(ioctx)
	if err != nil {
		return nil, err
	}

	images := []block.Image{}

	for _, imgName := range imgNames {
		image, err := r.openImage(ioctx, pool, imgName)
		if err != nil {
			return nil, err
		}
		images = append(images, *image)
	}

	return images, nil
}

func (r *imageRepo) Get(_ context.Context, scope, pool, image string) (*block.Image, error) {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return nil, err
	}

	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		return nil, err
	}
	defer ioctx.Destroy()

	return r.openImage(ioctx, pool, image)
}

func (r *imageRepo) Create(_ context.Context, scope, pool, image string, order int, stripeUnit, stripeCount, size, features uint64) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		return err
	}
	defer ioctx.Destroy()

	_, err = rbd.Create3(ioctx, image, size, features, order, stripeUnit, stripeCount)
	return err
}

func (r *imageRepo) Resize(_ context.Context, scope, pool, image string, size uint64) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		return err
	}
	defer ioctx.Destroy()

	img, err := rbd.OpenImage(ioctx, image, rbd.NoSnapshot)
	if err != nil {
		return err
	}
	defer img.Close()

	return img.Resize(size)
}

func (r *imageRepo) Delete(_ context.Context, scope, pool, image string) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		return err
	}
	defer ioctx.Destroy()

	return rbd.RemoveImage(ioctx, image)
}

func (r *imageRepo) sortASnapInfo(snaps []rbd.SnapInfo) {
	slices.SortFunc(snaps, func(a, b rbd.SnapInfo) int {
		if a.Id < b.Id {
			return -1
		} else if a.Id > b.Id {
			return 1
		}
		return 0
	})
}

func (r *imageRepo) appendSnapInfo(snaps []rbd.SnapInfo, size uint64) []rbd.SnapInfo {
	id := uint64(0)

	if len(snaps) > 0 {
		id = snaps[len(snaps)-1].Id + 1
	}

	return append(snaps, rbd.SnapInfo{
		Id:   id,
		Size: size,
	})
}

func (r *imageRepo) diskUsage(img *rbd.Image, previous, current string, size uint64) (uint64, error) {
	if err := img.SetSnapshot(current); err != nil {
		return 0, err
	}

	du := uint64(0)

	if err := img.DiffIterate(rbd.DiffIterateConfig{
		SnapName:    previous,
		Length:      size,
		WholeObject: rbd.EnableWholeObject,
		Callback: func(_, length uint64, exists int, _ any) int {
			if exists > 0 {
				du += length
			}
			return 0
		},
	}); err != nil {
		return 0, err
	}

	return du, nil
}

func (r *imageRepo) imageDiskUsage(img *rbd.Image, mirrorMode rbd.ImageMirrorMode, features uint64, snaps []rbd.SnapInfo) ([]block.ImageSnapshot, error) {
	snapshots := []block.ImageSnapshot{}
	previous := ""

	for _, info := range snaps {
		snap := img.GetSnapshot(info.Name)
		protected, _ := snap.IsProtected()

		snapshot := block.ImageSnapshot{
			Name:      info.Name,
			Quota:     info.Size,
			Protected: protected,
		}

		if mirrorMode != rbd.ImageMirrorModeSnapshot && r.featureOn(features, block.ImageFeatureFastDiff) {
			du, err := r.diskUsage(img, previous, info.Name, info.Size)
			if err != nil {
				return nil, err
			}
			snapshot.Used = du
		}

		snapshots = append(snapshots, snapshot)
		previous = info.Name
	}

	return snapshots, nil
}

func (r *imageRepo) openImage(ioctx *rados.IOContext, pool, image string) (*block.Image, error) {
	img, err := rbd.OpenImage(ioctx, image, rbd.NoSnapshot)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	info, err := img.Stat()
	if err != nil {
		return nil, err
	}

	stripeUnit, _ := img.GetStripeUnit()
	stripeCount, _ := img.GetStripeCount()
	features, _ := img.GetFeatures()
	timestamp, _ := img.GetCreateTimestamp()
	snaps, _ := img.GetSnapshotNames()
	mirrorMode, _ := img.GetImageMirrorMode()

	// // disk usage
	r.sortASnapInfo(snaps)
	snaps = r.appendSnapInfo(snaps, info.Size)
	snapshots, err := r.imageDiskUsage(img, mirrorMode, features, snaps)
	if err != nil {
		return nil, err
	}
	du := snapshots[len(snapshots)-1].Used
	snapshots = snapshots[:len(snapshots)-1]

	return &block.Image{
		Name:                 img.GetName(),
		PoolName:             pool,
		ObjectSize:           info.Obj_size,
		StripeUnit:           stripeUnit,
		StripeCount:          stripeCount,
		Quota:                info.Size,
		Used:                 du,
		ObjectCount:          info.Num_objs,
		FeatureLayering:      r.featureOn(features, block.ImageFeatureLayering),
		FeatureExclusiveLock: r.featureOn(features, block.ImageFeatureExclusiveLock),
		FeatureObjectMap:     r.featureOn(features, block.ImageFeatureObjectMap),
		FeatureFastDiff:      r.featureOn(features, block.ImageFeatureFastDiff),
		FeatureDeepFlatten:   r.featureOn(features, block.ImageFeatureDeepFlatten),
		CreatedAt:            time.Unix(timestamp.Sec, timestamp.Nsec),
		Snapshots:            snapshots,
	}, nil
}

func (r *imageRepo) featureOn(features, feature uint64) bool {
	return features&feature == feature
}
