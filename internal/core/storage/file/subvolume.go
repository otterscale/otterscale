package file

import (
	"context"
	"fmt"
	"time"
)

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

type SubvolumeExport struct {
	IP      string
	Path    string
	Clients []string
	Command string
}

// Note: Ceph create and update operations only return error status.
type SubvolumeRepo interface {
	List(ctx context.Context, scope string, volume, group string) ([]Subvolume, error)
	Get(ctx context.Context, scope string, volume, subvolume, group string) (*Subvolume, error)
	Create(ctx context.Context, scope string, volume, subvolume, group string, size uint64) error
	Resize(ctx context.Context, scope string, volume, subvolume, group string, size uint64) error
	Delete(ctx context.Context, scope string, volume, subvolume, group string) error
	ListExportClients(ctx context.Context, scope string, pool string) (map[string][]string, error)
}

func (uc *FileUseCase) ListSubvolumes(ctx context.Context, scope, volume, group string) ([]Subvolume, error) {
	subvolumes, err := uc.subvolume.List(ctx, scope, volume, group)
	if err != nil {
		return nil, err
	}

	clients, err := uc.subvolume.ListExportClients(ctx, scope, uc.nfsAppName(scope)) // TODO: multiple ceph-nfs charms
	if err != nil {
		return nil, err
	}

	vip, err := uc.vip(ctx, scope)
	if err != nil {
		return nil, err
	}

	for i := range subvolumes {
		subvolumes[i].Snapshots, err = uc.subvolumeSnapshot.List(ctx, scope, volume, subvolumes[i].Name, group)
		if err != nil {
			return nil, err
		}

		subvolumes[i].Export = &SubvolumeExport{
			IP:      vip,
			Path:    subvolumes[i].Path,
			Command: fmt.Sprintf("mount -t nfs4 %s:%s /mnt/%s", vip, subvolumes[i].Path, subvolumes[i].Name),
			Clients: clients[subvolumes[i].Path],
		}
	}

	return subvolumes, nil
}

func (uc *FileUseCase) CreateSubvolume(ctx context.Context, scope, volume, subvolume, group string, size uint64, export bool) (*Subvolume, error) {
	if !export {
		if err := uc.subvolume.Create(ctx, scope, volume, subvolume, group, size); err != nil {
			return nil, err
		}

		return uc.subvolume.Get(ctx, scope, volume, subvolume, group)
	}

	params := map[string]any{
		"name": subvolume,
		"size": size / 1024 / 1024 / 1024, // gb
	}

	if _, err := uc.action.Run(ctx, scope, uc.nfsAppName(scope), "create-share", params); err != nil {
		return nil, err
	}

	return uc.subvolume.Get(ctx, scope, volume, subvolume, "")
}

func (uc *FileUseCase) UpdateSubvolume(ctx context.Context, scope, volume, subvolume, group string, size uint64) (*Subvolume, error) {
	export, err := uc.isSubvolumeExport(ctx, scope, volume, subvolume, group)
	if err != nil {
		return nil, err
	}

	if !export {
		if err := uc.subvolume.Resize(ctx, scope, volume, subvolume, group, size); err != nil {
			return nil, err
		}

		return uc.subvolume.Get(ctx, scope, volume, subvolume, group)
	}

	params := map[string]any{
		"name": subvolume,
		"size": size / 1024 / 1024 / 1024, // gb
	}

	if _, err := uc.action.Run(ctx, scope, uc.nfsAppName(scope), "resize-share", params); err != nil {
		return nil, err
	}

	return uc.subvolume.Get(ctx, scope, volume, subvolume, "")
}

func (uc *FileUseCase) DeleteSubvolume(ctx context.Context, scope, volume, subvolume, group string) error {
	export, err := uc.isSubvolumeExport(ctx, scope, volume, subvolume, group)
	if err != nil {
		return err
	}
	if !export {
		return uc.subvolume.Delete(ctx, scope, volume, subvolume, group)
	}

	params := map[string]any{
		"name":  subvolume,
		"purge": true,
	}

	_, err = uc.action.Run(ctx, scope, uc.nfsAppName(scope), "delete-share", params)
	return err
}

func (uc *FileUseCase) GrantSubvolumeClient(ctx context.Context, scope, subvolume, clientIP string) error {
	params := map[string]any{
		"name":   subvolume,
		"client": clientIP,
	}

	_, err := uc.action.Run(ctx, scope, uc.nfsAppName(scope), "grant-access", params)
	return err
}
func (uc *FileUseCase) RevokeSubvolumeClient(ctx context.Context, scope, subvolume, clientIP string) error {
	params := map[string]any{
		"name":   subvolume,
		"client": clientIP,
	}

	_, err := uc.action.Run(ctx, scope, uc.nfsAppName(scope), "revoke-access", params)
	return err
}

// TODO: multiple ceph-nfs charms
func (uc *FileUseCase) isSubvolumeExport(ctx context.Context, scope, volume, subvolume, group string) (bool, error) {
	m, err := uc.subvolume.ListExportClients(ctx, scope, uc.nfsAppName(scope))
	if err != nil {
		return false, err
	}

	sub, err := uc.subvolume.Get(ctx, scope, volume, subvolume, group)
	if err != nil {
		return false, err
	}

	_, ok := m[sub.Path]
	return ok, nil
}

func (uc *FileUseCase) vip(ctx context.Context, scope string) (string, error) {
	config, err := uc.facility.Config(ctx, scope, uc.nfsAppName(scope))
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

func (uc *FileUseCase) nfsAppName(scope string) string {
	return scope + "-ceph-nfs"
}
