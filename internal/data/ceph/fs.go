package ceph

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/ceph/go-ceph/rados"
	"golang.org/x/sync/errgroup"

	"github.com/openhdc/otterscale/internal/core"
)

var (
	validPath    = regexp.MustCompile(`Path = "(?<Path>.+)";`)
	validClients = regexp.MustCompile(`Clients = (?<Clients>.+);`)
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
	fsDump, err := dumpFS(conn)
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
	subs, err := listSubvolumes(conn, volume, group)
	if err != nil {
		return nil, err
	}
	subvolumes := []core.Subvolume{}
	for _, sub := range subs {
		info, err := getSubvolume(conn, volume, sub.Name, group)
		if err != nil {
			return nil, err
		}
		subvolumes = append(subvolumes, *r.toSubvolume(sub.Name, info))
	}
	return subvolumes, nil
}

func (r *fs) GetSubvolume(ctx context.Context, config *core.StorageConfig, volume, subvolume, group string) (*core.Subvolume, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	info, err := getSubvolume(conn, volume, subvolume, group)
	if err != nil {
		return nil, err
	}
	return r.toSubvolume(subvolume, info), nil
}

func (r *fs) CreateSubvolume(ctx context.Context, config *core.StorageConfig, volume, subvolume, group string, size uint64) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return createSubvolume(conn, volume, subvolume, group, size)
}

func (r *fs) ResizeSubvolume(ctx context.Context, config *core.StorageConfig, volume, subvolume, group string, size uint64) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return resizeSubvolume(conn, volume, subvolume, group, size)
}

func (r *fs) DeleteSubvolume(ctx context.Context, config *core.StorageConfig, volume, subvolume, group string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return removeSubvolume(conn, volume, subvolume, group)
}

func (r *fs) ListSubvolumeSnapshots(ctx context.Context, config *core.StorageConfig, volume, subvolume, group string) ([]core.SubvolumeSnapshot, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	ss, err := listSubvolumeSnapshots(conn, volume, subvolume, group)
	if err != nil {
		return nil, err
	}
	var snapshots []core.SubvolumeSnapshot
	for _, s := range ss {
		snapshot, err := getSubvolumeSnapshot(conn, volume, subvolume, group, s.Name)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, *r.toSubvolumeSnapshot(s.Name, snapshot))
	}
	return snapshots, nil
}

func (r *fs) GetSubvolumeSnapshot(ctx context.Context, config *core.StorageConfig, volume, subvolume, group, snapshot string) (*core.SubvolumeSnapshot, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	info, err := getSubvolumeSnapshot(conn, volume, subvolume, group, snapshot)
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
	return createSubvolumeSnapshot(conn, volume, subvolume, group, snapshot)
}

func (r *fs) DeleteSubvolumeSnapshot(ctx context.Context, config *core.StorageConfig, volume, subvolume, group, snapshot string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return removeSubvolumeSnapshot(conn, volume, subvolume, group, snapshot)
}

func (r *fs) ListSubvolumeGroups(ctx context.Context, config *core.StorageConfig, volume string) ([]core.SubvolumeGroup, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	groups, err := listSubvolumeGroups(conn, volume)
	if err != nil {
		return nil, err
	}
	subvolumeGroups := []core.SubvolumeGroup{}
	for _, group := range groups {
		info, err := getSubvolumeGroup(conn, volume, group.Name)
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
	info, err := getSubvolumeGroup(conn, volume, group)
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
	return createSubvolumeGroup(conn, volume, group, size)
}

func (r *fs) ResizeSubvolumeGroup(ctx context.Context, config *core.StorageConfig, volume, group string, size uint64) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return resizeSubvolumeGroup(conn, volume, group, size)
}

func (r *fs) DeleteSubvolumeGroup(ctx context.Context, config *core.StorageConfig, volume, group string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return removeSubvolumeGroup(conn, volume, group)
}

func (r *fs) ListPathToExportClients(ctx context.Context, config *core.StorageConfig, pool string) (map[string][]string, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}

	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		return nil, err
	}

	indices, err := r.exportIndice(ioctx)
	if err != nil {
		return nil, err
	}

	eg, _ := errgroup.WithContext(ctx)
	results := make([][]string, len(indices))
	for i, index := range indices {
		eg.Go(func() error {
			results[i], _ = r.exportIndex(ioctx, index) // skip error
			return nil
		})
	}
	_ = eg.Wait() // skip error

	m := map[string][]string{}
	for _, result := range results {
		if len(result) == 0 {
			continue
		}
		m[result[0]] = result[1:]
	}
	return m, nil
}

func (r *fs) toVolumes(d *fsDump) []core.Volume {
	ret := []core.Volume{}
	for i := range d.FileSystems {
		ret = append(ret, core.Volume{
			ID:   d.FileSystems[i].ID,
			Name: d.FileSystems[i].MDSMap.FileSystemName,
		})
	}
	return ret
}

func (r *fs) toSubvolume(name string, info *subvolumeInfo) *core.Subvolume {
	ret := &core.Subvolume{
		Name:      name,
		Path:      info.Path,
		Mode:      fmt.Sprintf("%o", info.Mode),
		PoolName:  info.DataPool,
		Quota:     info.BytesUsed,
		Used:      info.BytesUsed,
		CreatedAt: info.CreatedAt.Time,
	}
	return ret
}

func (r *fs) toSubvolumeSnapshot(name string, info *subvolumeSnapshotInfo) *core.SubvolumeSnapshot {
	ret := &core.SubvolumeSnapshot{
		Name:             name,
		HasPendingClones: info.HasPendingClones,
		CreatedAt:        info.CreatedAt.Time,
	}
	return ret
}

func (r *fs) toSubvolumeGroups(name string, info *subvolumeGroupInfo) *core.SubvolumeGroup {
	ret := &core.SubvolumeGroup{
		Name:      name,
		Mode:      fmt.Sprintf("%06o", info.Mode),
		PoolName:  info.DataPool,
		Quota:     info.BytesUsed,
		Used:      info.BytesUsed,
		CreatedAt: info.CreatedAt.Time,
	}
	return ret
}

func (r *fs) exportIndice(ioctx *rados.IOContext) ([]string, error) {
	buffer := make([]byte, 1024)
	n, err := ioctx.Read("ganesha-export-index", buffer, 0)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(buffer[:n]), "\n%url rados://ceph-nfs/")
	return slices.DeleteFunc(lines, func(s string) bool { return s == "" }), nil
}

func (r *fs) exportIndex(ioctx *rados.IOContext, index string) ([]string, error) {
	buffer := make([]byte, 1024)
	n, err := ioctx.Read(index, buffer, 0)
	if err != nil {
		return nil, err
	}

	if !validPath.Match(buffer[:n]) || !validClients.Match(buffer[:n]) {
		return nil, fmt.Errorf("export index %q not found", index)
	}

	ret := []string{
		validPath.FindStringSubmatch(string(buffer[:n]))[1],
	}

	clients := validClients.FindStringSubmatch(string(buffer[:n]))[1]
	for _, client := range strings.Split(clients, ",") {
		ret = append(ret, strings.TrimSpace(client))
	}

	return ret, nil
}
