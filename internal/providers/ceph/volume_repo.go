package ceph

import (
	"context"

	"github.com/otterscale/otterscale/internal/core/storage/file"
)

type volumeRepo struct {
	ceph *Ceph
}

func NewVolumeRepo(ceph *Ceph) file.VolumeRepo {
	return &volumeRepo{
		ceph: ceph,
	}
}

var _ file.VolumeRepo = (*volumeRepo)(nil)

func (r *volumeRepo) List(_ context.Context, scope string) ([]file.Volume, error) {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return nil, err
	}

	dump, err := dumpFS(conn)
	if err != nil {
		return nil, err
	}

	return r.toVolumes(dump), nil
}

func (r *volumeRepo) toVolumes(d *fsDump) []file.Volume {
	ret := []file.Volume{}

	for i := range d.FileSystems {
		ret = append(ret, file.Volume{
			ID:        d.FileSystems[i].ID,
			Name:      d.FileSystems[i].MDSMap.FileSystemName,
			CreatedAt: d.FileSystems[i].MDSMap.Created.Time,
		})
	}

	return ret
}
