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

func (r *rbd) ListImages(ctx context.Context, config *core.StorageConfig, pools ...string) ([]core.RBDImage, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	images := []core.RBDImage{}
	for _, pool := range pools {
		imgs, err := r.listImages(conn, pool)
		if err != nil {
			return nil, err
		}
		images = append(images, imgs...)
	}
	return images, nil
}

func (r *rbd) listImages(conn *rados.Conn, pool string) ([]core.RBDImage, error) {
	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		return nil, err
	}
	defer ioctx.Destroy()

	imgs := []core.RBDImage{}
	imgNames, err := cephrbd.GetImageNames(ioctx)
	if err != nil {
		return nil, err
	}
	for _, imgName := range imgNames {
		img, err := r.openImage(ioctx, imgName)
		if err != nil {
			return nil, err
		}
		imgs = append(imgs, *img)
	}
	return imgs, nil
}

func (r *rbd) openImage(ioctx *rados.IOContext, imgName string) (*core.RBDImage, error) {
	img, err := cephrbd.OpenImage(ioctx, imgName, cephrbd.NoSnapshot)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	size, _ := img.GetSize()

	// fmt.Println(img.GetName())
	// fmt.Println(img.GetSize())
	// fmt.Println(img.GetFeatures())
	// fmt.Println(img.GetImageMirrorMode())
	// fmt.Println(img.GetMirrorImageInfo())

	return &core.RBDImage{
		Name: img.GetName(),
		Size: size,
	}, nil
}
