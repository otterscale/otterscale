package ceph

import (
	"context"
	"encoding/json"

	"github.com/ceph/go-ceph/rados"

	"github.com/openhdc/otterscale/internal/core"
)

type fs struct {
	ceph *Ceph
}

func NewFS(ceph *Ceph) core.CephFSRepo {
	return &fs{
		ceph: ceph,
	}
}

var _ core.CephFSRepo = (*fs)(nil)

func (r *fs) ListVolumes(ctx context.Context, config *core.StorageConfig) ([]core.Volume, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	fsDump, err := r.fsDump(conn)
	if err != nil {
		return nil, err
	}
	return r.toVolumes(fsDump), nil
}

func (r *fs) ListSubvolumes(ctx context.Context, config *core.StorageConfig, volume, group string) ([]core.Subvolume, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	subs, err := r.fsSubvolumes(conn, volume, group)
	if err != nil {
		return nil, err
	}
	subvolumes := []core.Subvolume{}
	for _, sub := range subs {
		info, err := r.fsSubvolumeInfo(conn, volume, sub.Name, group)
		if err != nil {
			return nil, err
		}
		subvolumes = append(subvolumes, *r.toSubvolume(sub.Name, info))
	}
	return subvolumes, nil
}

func (r *fs) GetSubvolumeSnapshot(ctx context.Context, config *core.StorageConfig, volume, subvolume, group, snapshot string) (*core.SubvolumeSnapshot, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	info, err := r.fsSubvolumeSnapshotInfo(conn, volume, subvolume, group, snapshot)
	if err != nil {
		return nil, err
	}
	return r.toSubvolumeSnapshot(snapshot, info), nil
}

func (r *fs) CreateSubvolumeSnapshot(ctx context.Context, config *core.StorageConfig, volume, subvolume, group, snapshot string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return r.fsSubvolumeSnapshotCreate(conn, volume, subvolume, group, snapshot)
}

func (r *fs) DeleteSubvolumeSnapshot(ctx context.Context, config *core.StorageConfig, volume, subvolume, group, snapshot string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return r.fsSubvolumeSnapshotRemove(conn, volume, subvolume, group, snapshot)
}

func (r *fs) ListSubvolumeGroups(ctx context.Context, config *core.StorageConfig, volume string) ([]core.SubvolumeGroup, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	groups, err := r.fsSubvolumeGroups(conn, volume)
	if err != nil {
		return nil, err
	}
	subvolumeGroups := []core.SubvolumeGroup{}
	for _, group := range groups {
		info, err := r.fsSubvolumeGroupInfo(conn, volume, group.Name)
		if err != nil {
			return nil, err
		}
		subvolumeGroups = append(subvolumeGroups, *r.toSubvolumeGroups(group.Name, info))
	}
	return subvolumeGroups, nil
}

func (r *fs) GetSubvolumeGroup(ctx context.Context, config *core.StorageConfig, volume, group string) (*core.SubvolumeGroup, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	info, err := r.fsSubvolumeGroupInfo(conn, volume, group)
	if err != nil {
		return nil, err
	}
	return r.toSubvolumeGroups(group, info), nil
}

func (r *fs) CreateSubvolumeGroup(ctx context.Context, config *core.StorageConfig, volume, group string, size uint64) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return r.fsSubvolumeGroupCreate(conn, volume, group, size)
}

func (r *fs) ResizeSubvolumeGroup(ctx context.Context, config *core.StorageConfig, volume, group string, size uint64) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return r.fsSubvolumeGroupResize(conn, volume, group, size)
}

func (r *fs) DeleteSubvolumeGroup(ctx context.Context, config *core.StorageConfig, volume, group string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return r.fsSubvolumeGroupRemove(conn, volume, group)
}

func (r *fs) toVolumes(d *fsDump) []core.Volume {
	ret := []core.Volume{}
	for i := range d.Filesystems {
		ret = append(ret, core.Volume{
			Name: d.Filesystems[i].Mdsmap.FsName,
		})
	}
	return ret
}

func (r *fs) toSubvolume(name string, info *fsSubvolumeInfo) *core.Subvolume {
	ret := &core.Subvolume{
		Name: name,
	}
	return ret
}

func (r *fs) toSubvolumeSnapshot(name string, info *fsSubvolumeSnapshotInfo) *core.SubvolumeSnapshot {
	ret := &core.SubvolumeSnapshot{
		Name: name,
	}
	return ret
}

func (r *fs) toSubvolumeGroups(name string, info *fsSubvolumeGroupInfo) *core.SubvolumeGroup {
	ret := &core.SubvolumeGroup{
		Name: name,
	}
	return ret
}

func (r *fs) fsDump(conn *rados.Conn) (*fsDump, error) {
	cmd, err := json.Marshal(map[string]string{
		"prefix": "fs dump",
		"format": "json",
	})
	if err != nil {
		return nil, err
	}
	resp, _, err := conn.MonCommand(cmd)
	if err != nil {
		return nil, err
	}
	var fsDump fsDump
	if err := json.Unmarshal(resp, &fsDump); err != nil {
		return nil, err
	}
	return &fsDump, nil
}

func (r *fs) fsSubvolumes(conn *rados.Conn, volume, group string) ([]fsSubvolume, error) {
	cmd, err := json.Marshal(map[string]string{
		"prefix":     "fs subvolume ls",
		"vol_name":   volume,
		"group_name": group,
		"format":     "json",
	})
	if err != nil {
		return nil, err
	}
	resp, _, err := conn.MonCommand(cmd)
	if err != nil {
		return nil, err
	}
	var fsSubvolumes []fsSubvolume
	if err := json.Unmarshal(resp, &fsSubvolumes); err != nil {
		return nil, err
	}
	return fsSubvolumes, nil
}

func (r *fs) fsSubvolumeInfo(conn *rados.Conn, volume, subvolume, group string) (*fsSubvolumeInfo, error) {
	cmd, err := json.Marshal(map[string]string{
		"prefix":     "fs subvolume info",
		"vol_name":   volume,
		"sub_name":   subvolume,
		"group_name": group,
		"format":     "json",
	})
	if err != nil {
		return nil, err
	}
	resp, _, err := conn.MonCommand(cmd)
	if err != nil {
		return nil, err
	}
	var fsSubvolumeInfo fsSubvolumeInfo
	if err := json.Unmarshal(resp, &fsSubvolumeInfo); err != nil {
		return nil, err
	}
	return &fsSubvolumeInfo, nil
}

func (r *fs) fsSubvolumeSnapshotInfo(conn *rados.Conn, volume, subvolume, group, snapshot string) (*fsSubvolumeSnapshotInfo, error) {
	m := map[string]string{
		"prefix":    "fs subvolume snapshot info",
		"vol_name":  volume,
		"sub_name":  subvolume,
		"snap_name": snapshot,
		"format":    "json",
	}
	if group != "" {
		m["group_name"] = group
	}
	cmd, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	resp, _, err := conn.MonCommand(cmd)
	if err != nil {
		return nil, err
	}
	var fsSubvolumeSnapshotInfo fsSubvolumeSnapshotInfo
	if err := json.Unmarshal(resp, &fsSubvolumeSnapshotInfo); err != nil {
		return nil, err
	}
	return &fsSubvolumeSnapshotInfo, nil
}

func (r *fs) fsSubvolumeGroups(conn *rados.Conn, volume string) ([]fsSubvolumeGroup, error) {
	cmd, err := json.Marshal(map[string]string{
		"prefix":   "fs subvolumegroup ls",
		"vol_name": volume,
		"format":   "json",
	})
	if err != nil {
		return nil, err
	}
	resp, _, err := conn.MonCommand(cmd)
	if err != nil {
		return nil, err
	}
	var fsSubvolumeGroups []fsSubvolumeGroup
	if err := json.Unmarshal(resp, &fsSubvolumeGroups); err != nil {
		return nil, err
	}
	return fsSubvolumeGroups, nil
}

func (r *fs) fsSubvolumeGroupInfo(conn *rados.Conn, volume, group string) (*fsSubvolumeGroupInfo, error) {
	cmd, err := json.Marshal(map[string]string{
		"prefix":     "fs subvolumegroup info",
		"vol_name":   volume,
		"group_name": group,
		"format":     "json",
	})
	if err != nil {
		return nil, err
	}
	resp, _, err := conn.MonCommand(cmd)
	if err != nil {
		return nil, err
	}
	var fsSubvolumeGroupInfo fsSubvolumeGroupInfo
	if err := json.Unmarshal(resp, &fsSubvolumeGroupInfo); err != nil {
		return nil, err
	}
	return &fsSubvolumeGroupInfo, nil
}

func (r *fs) fsSubvolumeGroupCreate(conn *rados.Conn, volume, group string, size uint64) error {
	cmd, err := json.Marshal(map[string]any{
		"prefix":     "fs subvolumegroup create",
		"vol_name":   volume,
		"group_name": group,
		"size":       size,
		"format":     "json",
	})
	if err != nil {
		return err
	}
	if _, _, err := conn.MonCommand(cmd); err != nil {
		return err
	}
	return nil
}

func (r *fs) fsSubvolumeGroupResize(conn *rados.Conn, volume, group string, size uint64) error {
	cmd, err := json.Marshal(map[string]any{
		"prefix":     "fs subvolumegroup resize",
		"vol_name":   volume,
		"group_name": group,
		"new_size":   size,
		"format":     "json",
	})
	if err != nil {
		return err
	}
	if _, _, err := conn.MonCommand(cmd); err != nil {
		return err
	}
	return nil
}

func (r *fs) fsSubvolumeGroupRemove(conn *rados.Conn, volume, group string) error {
	cmd, err := json.Marshal(map[string]any{
		"prefix":     "fs subvolumegroup rm",
		"vol_name":   volume,
		"group_name": group,
		"format":     "json",
	})
	if err != nil {
		return err
	}
	if _, _, err := conn.MonCommand(cmd); err != nil {
		return err
	}
	return nil
}

func (r *fs) fsSubvolumeSnapshotCreate(conn *rados.Conn, volume, subvolume, group, snapshot string) error {
	m := map[string]any{
		"prefix":    "fs subvolume snapshot create",
		"vol_name":  volume,
		"sub_name":  subvolume,
		"snap_name": snapshot,
		"format":    "json",
	}
	if group != "" {
		m["group_name"] = group
	}
	cmd, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if _, _, err := conn.MonCommand(cmd); err != nil {
		return err
	}
	return nil
}

func (r *fs) fsSubvolumeSnapshotRemove(conn *rados.Conn, volume, subvolume, group, snapshot string) error {
	m := map[string]any{
		"prefix":    "fs subvolume snapshot rm",
		"vol_name":  volume,
		"sub_name":  subvolume,
		"snap_name": snapshot,
		"format":    "json",
	}
	if group != "" {
		m["group_name"] = group
	}
	cmd, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if _, _, err := conn.MonCommand(cmd); err != nil {
		return err
	}
	return nil
}
