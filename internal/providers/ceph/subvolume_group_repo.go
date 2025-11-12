package ceph

import (
	"context"

	"github.com/otterscale/otterscale/internal/core/storage/file"
)

type subvolumeGroupRepo struct {
	ceph *Ceph
}

func NewSubvolumeGroupRepo(ceph *Ceph) file.SubvolumeGroupRepo {
	return &subvolumeGroupRepo{
		ceph: ceph,
	}
}

var _ file.SubvolumeGroupRepo = (*subvolumeGroupRepo)(nil)

func (r *subvolumeGroupRepo) List(_ context.Context, scope, volume string) ([]file.SubvolumeGroup, error) {
	conn, err := r.ceph.Connection(scope)
	if err != nil {
		return nil, err
	}

	dumpNames, err := listSubvolumeGroups(conn, volume)
	if err != nil {
		return nil, err
	}

	subvolumeGroups := []file.SubvolumeGroup{}

	for _, dumpName := range dumpNames {
		dump, err := getSubvolumeGroup(conn, volume, dumpName.Name)
		if err != nil {
			return nil, err
		}
		subvolumeGroups = append(subvolumeGroups, *r.toSubvolumeGroups(dumpName.Name, dump))
	}

	return subvolumeGroups, nil
}

func (r *subvolumeGroupRepo) Get(_ context.Context, scope, volume, group string) (*file.SubvolumeGroup, error) {
	conn, err := r.ceph.Connection(scope)
	if err != nil {
		return nil, err
	}

	dump, err := getSubvolumeGroup(conn, volume, group)
	if err != nil {
		return nil, err
	}

	return r.toSubvolumeGroups(group, dump), nil
}

func (r *subvolumeGroupRepo) Create(_ context.Context, scope, volume, group string, size uint64) error {
	conn, err := r.ceph.Connection(scope)
	if err != nil {
		return err
	}

	return createSubvolumeGroup(conn, volume, group, size)
}

func (r *subvolumeGroupRepo) Resize(_ context.Context, scope, volume, group string, size uint64) error {
	conn, err := r.ceph.Connection(scope)
	if err != nil {
		return err
	}

	return resizeSubvolumeGroup(conn, volume, group, size)
}

func (r *subvolumeGroupRepo) Delete(_ context.Context, scope, volume, group string) error {
	conn, err := r.ceph.Connection(scope)
	if err != nil {
		return err
	}

	return removeSubvolumeGroup(conn, volume, group)
}

func (r *subvolumeGroupRepo) toSubvolumeGroups(name string, info *subvolumeGroupInfo) *file.SubvolumeGroup {
	// quota, _ := parseQuota(info.BytesQuota)
	ret := &file.SubvolumeGroup{
		// 	Name:      name,
		// 	Mode:      fmt.Sprintf("%06o", info.Mode),
		// 	PoolName:  info.DataPool,
		// 	Quota:     quota,
		// 	Used:      info.BytesUsed,
		// 	CreatedAt: info.CreatedAt.Time,
	}
	return ret
}
