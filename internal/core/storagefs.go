package core

import (
	"context"
	"fmt"
	"time"
)

type Volume struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}

type SubvolumeSnapshot struct {
	Name             string
	HasPendingClones string
	CreatedAt        time.Time
}

type SubvolumeExport struct {
	IP      string
	Path    string
	Clients []string
	Command string
}

type Subvolume struct {
	Name      string
	Path      string
	Mode      string
	PoolName  string
	Quota     uint64
	Used      uint64
	CreatedAt time.Time
	Export    *SubvolumeExport
	Snapshots []SubvolumeSnapshot
}

type SubvolumeGroup struct {
	Name      string
	Mode      string
	PoolName  string
	Quota     uint64
	Used      uint64
	CreatedAt time.Time
}

type CephFSRepo interface {
	ListVolumes(ctx context.Context, config *StorageConfig) ([]Volume, error)
	ListSubvolumes(ctx context.Context, config *StorageConfig, volume, group string) ([]Subvolume, error)
	GetSubvolume(ctx context.Context, config *StorageConfig, volume, subvolume, group string) (*Subvolume, error)
	ListSubvolumeSnapshots(ctx context.Context, config *StorageConfig, volume, subvolume, group string) ([]SubvolumeSnapshot, error)
	GetSubvolumeSnapshot(ctx context.Context, config *StorageConfig, volume, subvolume, group, snapshot string) (*SubvolumeSnapshot, error)
	CreateSubvolume(ctx context.Context, config *StorageConfig, volume, subvolume, group string, size uint64) error
	ResizeSubvolume(ctx context.Context, config *StorageConfig, volume, subvolume, group string, size uint64) error
	DeleteSubvolume(ctx context.Context, config *StorageConfig, volume, subvolume, group string) error
	CreateSubvolumeSnapshot(ctx context.Context, config *StorageConfig, volume, subvolume, group, snapshot string) error
	DeleteSubvolumeSnapshot(ctx context.Context, config *StorageConfig, volume, subvolume, group, snapshot string) error
	ListSubvolumeGroups(ctx context.Context, config *StorageConfig, volume string) ([]SubvolumeGroup, error)
	GetSubvolumeGroup(ctx context.Context, config *StorageConfig, volume, group string) (*SubvolumeGroup, error)
	CreateSubvolumeGroup(ctx context.Context, config *StorageConfig, volume, group string, size uint64) error
	ResizeSubvolumeGroup(ctx context.Context, config *StorageConfig, volume, group string, size uint64) error
	DeleteSubvolumeGroup(ctx context.Context, config *StorageConfig, volume, group string) error
	ListPathToExportClients(ctx context.Context, config *StorageConfig, pool string) (map[string][]string, error)
}

func (uc *StorageUseCase) ListVolumes(ctx context.Context, uuid, facility string) ([]Volume, error) {
	config, err := storageConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.fs.ListVolumes(ctx, config)
}

func (uc *StorageUseCase) ListSubvolumes(ctx context.Context, uuid, facility, volume, group string) ([]Subvolume, error) {
	config, err := storageConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	subs, err := uc.fs.ListSubvolumes(ctx, config, volume, group)
	if err != nil {
		return nil, err
	}
	pool := nfsName(facility) // TODO: multiple ceph-nfs charms
	m, err := uc.fs.ListPathToExportClients(ctx, config, pool)
	if err != nil {
		return nil, err
	}
	ip, err := uc.exportIP(ctx, uuid, nfsName(facility))
	if err != nil {
		return nil, err
	}
	for i := range subs {
		subs[i].Snapshots, err = uc.fs.ListSubvolumeSnapshots(ctx, config, volume, subs[i].Name, group)
		if err != nil {
			return nil, err
		}
		subs[i].Export = &SubvolumeExport{
			IP:      ip,
			Path:    subs[i].Path,
			Command: fmt.Sprintf("mount -t nfs4 %s:%s /mnt/%s", ip, subs[i].Path, subs[i].Name),
			Clients: m[subs[i].Path],
		}
	}
	return subs, nil
}

func (uc *StorageUseCase) CreateSubvolume(ctx context.Context, uuid, facility, volume, subvolume, group string, size uint64, export bool) (*Subvolume, error) {
	config, err := storageConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	if export {
		leader, err := uc.facility.GetLeader(ctx, uuid, nfsName(facility))
		if err != nil {
			return nil, err
		}
		action := "create-share"
		params := map[string]any{
			"name": subvolume,
			"size": size / 1024 / 1024 / 1024, // gb
		}
		if _, err := runAction(ctx, uc.action, uuid, leader, action, params); err != nil {
			return nil, err
		}
		return uc.fs.GetSubvolume(ctx, config, volume, subvolume, "")
	}

	if err := uc.fs.CreateSubvolume(ctx, config, volume, subvolume, group, size); err != nil {
		return nil, err
	}
	return uc.fs.GetSubvolume(ctx, config, volume, subvolume, group)
}

func (uc *StorageUseCase) UpdateSubvolume(ctx context.Context, uuid, facility, volume, subvolume, group string, size uint64) (*Subvolume, error) {
	config, err := storageConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	export, err := uc.subvolumeExport(ctx, config, facility, volume, subvolume, group)
	if err != nil {
		return nil, err
	}
	if export {
		leader, err := uc.facility.GetLeader(ctx, uuid, nfsName(facility))
		if err != nil {
			return nil, err
		}
		action := "resize-share"
		params := map[string]any{
			"name": subvolume,
			"size": size / 1024 / 1024 / 1024, // gb
		}
		if _, err := runAction(ctx, uc.action, uuid, leader, action, params); err != nil {
			return nil, err
		}
		return uc.fs.GetSubvolume(ctx, config, volume, subvolume, "")
	}

	if err := uc.fs.ResizeSubvolume(ctx, config, volume, subvolume, group, size); err != nil {
		return nil, err
	}
	return uc.fs.GetSubvolume(ctx, config, volume, subvolume, group)
}

func (uc *StorageUseCase) DeleteSubvolume(ctx context.Context, uuid, facility, volume, subvolume, group string) error {
	config, err := storageConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	export, err := uc.subvolumeExport(ctx, config, facility, volume, subvolume, group)
	if err != nil {
		return err
	}
	if export {
		leader, err := uc.facility.GetLeader(ctx, uuid, nfsName(facility))
		if err != nil {
			return err
		}
		action := "delete-share"
		params := map[string]any{
			"name":  subvolume,
			"purge": true,
		}
		if _, err := runAction(ctx, uc.action, uuid, leader, action, params); err != nil {
			return err
		}
		return nil
	}
	return uc.fs.DeleteSubvolume(ctx, config, volume, subvolume, group)
}

func (uc *StorageUseCase) GrantSubvolumeClient(ctx context.Context, uuid, facility, volume, subvolume, clientIP string) error {
	leader, err := uc.facility.GetLeader(ctx, uuid, nfsName(facility))
	if err != nil {
		return err
	}
	action := "grant-access"
	params := map[string]any{
		"name":   subvolume,
		"client": clientIP,
	}
	if _, err := runAction(ctx, uc.action, uuid, leader, action, params); err != nil {
		return err
	}
	return nil
}

func (uc *StorageUseCase) RevokeSubvolumeClient(ctx context.Context, uuid, facility, volume, subvolume, clientIP string) error {
	leader, err := uc.facility.GetLeader(ctx, uuid, nfsName(facility))
	if err != nil {
		return err
	}
	action := "revoke-access"
	params := map[string]any{
		"name":   subvolume,
		"client": clientIP,
	}
	if _, err := runAction(ctx, uc.action, uuid, leader, action, params); err != nil {
		return err
	}
	return nil
}

func (uc *StorageUseCase) CreateSubvolumeSnapshot(ctx context.Context, uuid, facility, volume, subvolume, group, snapshot string) (*SubvolumeSnapshot, error) {
	config, err := storageConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	if err := uc.fs.CreateSubvolumeSnapshot(ctx, config, volume, subvolume, group, snapshot); err != nil {
		return nil, err
	}
	return uc.fs.GetSubvolumeSnapshot(ctx, config, volume, subvolume, group, snapshot)
}

func (uc *StorageUseCase) DeleteSubvolumeSnapshot(ctx context.Context, uuid, facility, volume, subvolume, group, snapshot string) error {
	config, err := storageConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.fs.DeleteSubvolumeSnapshot(ctx, config, volume, subvolume, group, snapshot)
}

func (uc *StorageUseCase) ListSubvolumeGroups(ctx context.Context, uuid, facility, volume string) ([]SubvolumeGroup, error) {
	config, err := storageConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.fs.ListSubvolumeGroups(ctx, config, volume)
}

func (uc *StorageUseCase) CreateSubvolumeGroup(ctx context.Context, uuid, facility, volume, group string, size uint64) (*SubvolumeGroup, error) {
	config, err := storageConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	if err := uc.fs.CreateSubvolumeGroup(ctx, config, volume, group, size); err != nil {
		return nil, err
	}
	return uc.fs.GetSubvolumeGroup(ctx, config, volume, group)
}

func (uc *StorageUseCase) UpdateSubvolumeGroup(ctx context.Context, uuid, facility, volume, group string, size uint64) (*SubvolumeGroup, error) {
	config, err := storageConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	if err := uc.fs.ResizeSubvolumeGroup(ctx, config, volume, group, size); err != nil {
		return nil, err
	}
	return uc.fs.GetSubvolumeGroup(ctx, config, volume, group)
}

func (uc *StorageUseCase) DeleteSubvolumeGroup(ctx context.Context, uuid, facility, volume, group string) error {
	config, err := storageConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.fs.DeleteSubvolumeGroup(ctx, config, volume, group)
}

func (uc *StorageUseCase) subvolumeExport(ctx context.Context, config *StorageConfig, facility, volume, subvolume, group string) (bool, error) {
	pool := nfsName(facility) // TODO: multiple ceph-nfs charms
	m, err := uc.fs.ListPathToExportClients(ctx, config, pool)
	if err != nil {
		return false, err
	}
	sub, err := uc.fs.GetSubvolume(ctx, config, volume, subvolume, group)
	if err != nil {
		return false, err
	}
	_, ok := m[sub.Path]
	return ok, nil
}

func (uc *StorageUseCase) exportIP(ctx context.Context, uuid, facility string) (string, error) {
	config, err := uc.facility.GetConfig(ctx, uuid, facility)
	if err != nil {
		return "", err
	}
	vip, ok := config["vip"]
	if !ok {
		return "", nil
	}
	value, ok := vip.(map[string]any)["value"]
	if !ok {
		return "", nil
	}
	return value.(string), nil
}
