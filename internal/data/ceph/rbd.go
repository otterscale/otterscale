package ceph

import (
	"context"
	"slices"
	"time"

	"github.com/ceph/go-ceph/rados"
	cephrbd "github.com/ceph/go-ceph/rbd"

	"github.com/otterscale/otterscale/internal/core"
)

type rbd struct {
	ceph *Ceph
}

func NewRBD(ceph *Ceph) core.CephRBDRepo {
	return &rbd{
		ceph: ceph,
	}
}

var _ core.CephRBDRepo = (*rbd)(nil)

func (r *rbd) ListImages(ctx context.Context, config *core.StorageConfig, pool string) ([]core.RBDImage, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}

	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		return nil, err
	}
	defer ioctx.Destroy()

	imgNames, err := cephrbd.GetImageNames(ioctx)
	if err != nil {
		return nil, err
	}

	imgs := []core.RBDImage{}
	for _, imgName := range imgNames {
		img, err := r.openImage(ioctx, pool, imgName)
		if err != nil {
			return nil, err
		}
		imgs = append(imgs, *img)
	}
	return imgs, nil
}

func (r *rbd) GetImage(ctx context.Context, config *core.StorageConfig, pool, image string) (*core.RBDImage, error) {
	conn, err := r.ceph.connection(config)
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

func (r *rbd) CreateImage(ctx context.Context, config *core.StorageConfig, pool, image string, order int, stripeUnit, stripeCount, size, features uint64) (*core.RBDImage, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}

	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		return nil, err
	}
	defer ioctx.Destroy()

	if _, err := cephrbd.Create3(ioctx, image, size, features, order, stripeUnit, stripeCount); err != nil {
		return nil, err
	}
	return r.openImage(ioctx, pool, image)
}

func (r *rbd) UpdateImageSize(ctx context.Context, config *core.StorageConfig, pool, image string, size uint64) error {
	conn, err := r.ceph.connection(config)
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

	return img.Resize(size)
}

func (r *rbd) DeleteImage(ctx context.Context, config *core.StorageConfig, pool, image string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}

	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		return err
	}
	defer ioctx.Destroy()

	return cephrbd.RemoveImage(ioctx, image)
}

func (r *rbd) CreateImageSnapshot(ctx context.Context, config *core.StorageConfig, pool, image, snapshot string) error {
	conn, err := r.ceph.connection(config)
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

func (r *rbd) DeleteImageSnapshot(ctx context.Context, config *core.StorageConfig, pool, image, snapshot string) error {
	conn, err := r.ceph.connection(config)
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

func (r *rbd) RollbackImageSnapshot(ctx context.Context, config *core.StorageConfig, pool, image, snapshot string) error {
	conn, err := r.ceph.connection(config)
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

func (r *rbd) ProtectImageSnapshot(ctx context.Context, config *core.StorageConfig, pool, image, snapshot string) error {
	conn, err := r.ceph.connection(config)
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

func (r *rbd) UnprotectImageSnapshot(ctx context.Context, config *core.StorageConfig, pool, image, snapshot string) error {
	conn, err := r.ceph.connection(config)
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

func (r *rbd) openImage(ioctx *rados.IOContext, pool, image string) (*core.RBDImage, error) {
	img, err := cephrbd.OpenImage(ioctx, image, cephrbd.NoSnapshot)
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

	// disk usage
	snaps = r.sortAndAppendSnapInfo(snaps, info.Size)
	snapshots, err := r.imageDiskUsage(img, mirrorMode, features, snaps)
	if err != nil {
		return nil, err
	}
	du := snapshots[len(snapshots)-1].Used
	snapshots = snapshots[:len(snapshots)-1]

	return &core.RBDImage{
		Name:                 img.GetName(),
		PoolName:             pool,
		ObjectSize:           info.Obj_size,
		StripeUnit:           stripeUnit,
		StripeCount:          stripeCount,
		Quota:                info.Size,
		Used:                 du,
		ObjectCount:          info.Num_objs,
		FeatureLayering:      r.featureOn(features, cephrbd.FeatureLayering),
		FeatureExclusiveLock: r.featureOn(features, cephrbd.FeatureExclusiveLock),
		FeatureObjectMap:     r.featureOn(features, cephrbd.FeatureObjectMap),
		FeatureFastDiff:      r.featureOn(features, cephrbd.FeatureFastDiff),
		FeatureDeepFlatten:   r.featureOn(features, cephrbd.FeatureDeepFlatten),
		CreatedAt:            time.Unix(timestamp.Sec, timestamp.Nsec),
		Snapshots:            snapshots,
	}, nil
}

func (r *rbd) featureOn(features, feature uint64) bool {
	return features&feature == feature
}

func (r *rbd) imageDiskUsage(img *cephrbd.Image, mirrorMode cephrbd.ImageMirrorMode, features uint64, snaps []cephrbd.SnapInfo) ([]core.RBDImageSnapshot, error) {
	snapshots := []core.RBDImageSnapshot{}
	previous := ""
	for _, info := range snaps {
		snap := img.GetSnapshot(info.Name)
		protected, _ := snap.IsProtected()
		snapshot := core.RBDImageSnapshot{
			Name:      info.Name,
			Quota:     info.Size,
			Protected: protected,
		}
		if mirrorMode != cephrbd.ImageMirrorModeSnapshot && r.featureOn(features, cephrbd.FeatureFastDiff) {
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

func (r *rbd) sortAndAppendSnapInfo(snaps []cephrbd.SnapInfo, size uint64) []cephrbd.SnapInfo {
	slices.SortFunc(snaps, func(a, b cephrbd.SnapInfo) int {
		if a.Id < b.Id {
			return -1
		} else if a.Id > b.Id {
			return 1
		}
		return 0
	})

	id := uint64(0)
	if len(snaps) > 0 {
		id = uint64(len(snaps) + 1) //nolint:gosec
	}

	return append(snaps, cephrbd.SnapInfo{
		Id:   id,
		Size: size,
	})
}

func (r *rbd) diskUsage(img *cephrbd.Image, previous, current string, size uint64) (uint64, error) {
	du := uint64(0)
	if err := img.SetSnapshot(current); err != nil {
		return 0, err
	}
	if err := img.DiffIterate(cephrbd.DiffIterateConfig{
		SnapName:    previous,
		Length:      size,
		WholeObject: cephrbd.EnableWholeObject,
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
