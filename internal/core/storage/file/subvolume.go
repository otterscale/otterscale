package file

import (
	"context"
	"time"
)

type Bytes uint64
type UnixMode uint32

type SubvolumeState string
type Feature string

type SubvolumeKey struct {
	VolumeName    string
	GroupName     string
	SubvolumeName string
}

type SubvolumeInfo struct {
	Path         string
	State        SubvolumeState
	UID          uint32
	GID          uint32
	Mode         UnixMode
	BytesPercent float64
	BytesUsed    Bytes
	BytesQuota   *Bytes

	DataPool      string
	PoolNamespace string

	Atime, Mtime, Ctime, CreatedAt time.Time
	Features                       []Feature
}

type Subvolume struct {
	Key  SubvolumeKey
	Info SubvolumeInfo
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
	Get(ctx context.Context, scope string, volume, group, subvolume string) (*Subvolume, error)
	Create(ctx context.Context, scope string, volume, group, subvolume string, size *Bytes, uid, gid *uint32, mode *UnixMode, poolLayout *string, isNamespaceIsolated *bool) error
	Resize(ctx context.Context, scope string, volume, group, subvolume string, newSize Bytes, noShrink *bool) error
	Delete(ctx context.Context, scope string, volume, group, subvolume string, isForce *bool) error
}

func (uc *UseCase) ListSubvolumes(ctx context.Context, scope, volume, group string) ([]Subvolume, error) {
	subvolumes, err := uc.subvolume.List(ctx, scope, volume, group)
	if err != nil {
		return nil, err
	}

	return subvolumes, nil
}

func (uc *UseCase) CreateSubvolume(ctx context.Context, scope, volume, group, subvolume string, size *Bytes, uid, gid *uint32, mode *UnixMode, poolLayout *string, isNamespaceIsolated *bool) (*Subvolume, error) {
	if err := uc.subvolume.Create(ctx, scope, volume, group, subvolume, size, uid, gid, mode, poolLayout, isNamespaceIsolated); err != nil {
		return nil, err
	}

	return uc.subvolume.Get(ctx, scope, volume, group, subvolume)
}

func (uc *UseCase) GetSubvolume(ctx context.Context, scope, volume, group, subvolume string) (*Subvolume, error) {
	return uc.subvolume.Get(ctx, scope, volume, group, subvolume)
}

func (uc *UseCase) UpdateSubvolume(ctx context.Context, scope, volume, group, subvolume string, newSize Bytes, noShrink *bool) (*Subvolume, error) {
	if err := uc.subvolume.Resize(ctx, scope, volume, group, subvolume, newSize, noShrink); err != nil {
		return nil, err
	}
	return uc.subvolume.Get(ctx, scope, volume, group, subvolume)
}

func (uc *UseCase) DeleteSubvolume(ctx context.Context, scope, volume, group, subvolume string, isForce *bool) error {
	return uc.subvolume.Delete(ctx, scope, volume, group, subvolume, isForce)
}
