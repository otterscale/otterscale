package ceph

import (
	"context"

	"github.com/ceph/go-ceph/rados"
	cephrbd "github.com/ceph/go-ceph/rbd"

	"github.com/openhdc/otterscale/internal/core"
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

func (r *rbd) ListImages(ctx context.Context, config *core.StorageConfig, poolName string) ([]core.RBDImage, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}

	ioctx, err := conn.OpenIOContext(poolName)
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
		img, err := r.openImage(ioctx, imgName)
		if err != nil {
			return nil, err
		}
		imgs = append(imgs, *img)
	}
	return imgs, nil
}

func (r *rbd) GetImage(ctx context.Context, config *core.StorageConfig, poolName, imageName string) (*core.RBDImage, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}

	ioctx, err := conn.OpenIOContext(poolName)
	if err != nil {
		return nil, err
	}
	defer ioctx.Destroy()

	return r.openImage(ioctx, imageName)
}

func (r *rbd) CreateImage(ctx context.Context, config *core.StorageConfig, poolName, imageName string, order int, stripeUnit, stripeCount, size, features uint64) (*core.RBDImage, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}

	ioctx, err := conn.OpenIOContext(poolName)
	if err != nil {
		return nil, err
	}
	defer ioctx.Destroy()

	if _, err := cephrbd.Create3(ioctx, imageName, size, features, order, stripeUnit, stripeCount); err != nil {
		return nil, err
	}
	return r.openImage(ioctx, imageName)
}

func (r *rbd) UpdateImageSize(ctx context.Context, config *core.StorageConfig, poolName, imageName string, size uint64) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}

	ioctx, err := conn.OpenIOContext(poolName)
	if err != nil {
		return err
	}
	defer ioctx.Destroy()

	img, err := cephrbd.OpenImage(ioctx, imageName, cephrbd.NoSnapshot)
	if err != nil {
		return err
	}
	return img.Resize(size)
}

func (r *rbd) DeleteImage(ctx context.Context, config *core.StorageConfig, poolName, imageName string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}

	ioctx, err := conn.OpenIOContext(poolName)
	if err != nil {
		return err
	}
	defer ioctx.Destroy()

	return cephrbd.RemoveImage(ioctx, imageName)
}

// func (r *rbd) ListSnapshots(ctx context.Context, config *core.StorageConfig, pools ...string) ([]core.RBDImage, error) {
// 	conn, err := r.ceph.connection(config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	images := []core.RBDImage{}
// 	for _, pool := range pools {
// 		imgs, err := r.listImages(conn, pool)
// 		if err != nil {
// 			return nil, err
// 		}
// 		images = append(images, imgs...)
// 	}
// 	return images, nil
// }

func (r *rbd) openImage(ioctx *rados.IOContext, imgName string) (*core.RBDImage, error) {
	img, err := cephrbd.OpenImage(ioctx, imgName, cephrbd.NoSnapshot)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	size, _ := img.GetSize()
	features, _ := img.GetFeatures()

	// img.GetSnapshotNames()

	// fmt.Println(img.GetName())
	// fmt.Println(img.GetSize())
	// fmt.Println(img.GetFeatures())
	// fmt.Println(img.GetImageMirrorMode())
	// fmt.Println(img.GetMirrorImageInfo())

	return &core.RBDImage{
		Name:     img.GetName(),
		Size:     size,
		Features: features,
	}, nil
}
