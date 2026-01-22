package ceph

import (
	"context"
	"regexp"

	"github.com/otterscale/otterscale/internal/core/storage/file"
)

const indexSeparator = "\n%url rados://ceph-nfs/"

var (
	validPath    = regexp.MustCompile(`Path = "(?<Path>.+)";`)
	validClients = regexp.MustCompile(`Clients = (?<Clients>.+);`)
)

// Note: Ceph API do not support context.
type subvolumeRepo struct {
	ceph *Ceph
}

func NewSubvolumeRepo(ceph *Ceph) file.SubvolumeRepo {
	return &subvolumeRepo{
		ceph: ceph,
	}
}

var _ file.SubvolumeRepo = (*subvolumeRepo)(nil)

func (r *subvolumeRepo) List(_ context.Context, scope, volume, group string) ([]file.Subvolume, error) {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return nil, err
	}

	dumpNames, err := listSubvolumes(conn, volume, group)
	if err != nil {
		return nil, err
	}

	subvolumes := []file.Subvolume{}

	for i := range dumpNames {
		info, err := getSubvolume(conn, volume, group, dumpNames[i].Name)
		if err != nil {
			return nil, err
		}
		subvolumes = append(subvolumes, *r.toSubvolume(volume, group, dumpNames[i].Name, info))
	}

	return subvolumes, nil
}

func (r *subvolumeRepo) Get(_ context.Context, scope, volume, group, subvolume string) (*file.Subvolume, error) {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return nil, err
	}

	info, err := getSubvolume(conn, volume, group, subvolume)
	if err != nil {
		return nil, err
	}

	return r.toSubvolume(volume, group, subvolume, info), nil
}

func (r *subvolumeRepo) Create(_ context.Context, scope, volume, group, subvolume string, size *file.Bytes, uid, gid *uint32, mode *file.UnixMode, poolLayout *string, isNamespaceIsolated *bool) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}
	var cephSize *uint64
	if size != nil {
		v := uint64(*size)
		cephSize = &v
	}
	var subvolumeMode *uint32
	if mode != nil {
		v := uint32(*mode)
		subvolumeMode = &v
	}

	return createSubvolume(conn, volume, group, subvolume, cephSize, uid, gid, subvolumeMode, poolLayout, isNamespaceIsolated)
}

func (r *subvolumeRepo) Resize(_ context.Context, scope, volume, group, subvolume string, newSize file.Bytes, noShrink *bool) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	return resizeSubvolume(conn, volume, group, subvolume, uint64(newSize), noShrink)
}

func (r *subvolumeRepo) Delete(_ context.Context, scope, volume, group, subvolume string, isForce *bool) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	return removeSubvolume(conn, volume, group, subvolume, isForce)
}

/*func (r *subvolumeRepo) ListExportClients(_ context.Context, scope, pool string) (map[string][]string, error) {
	conn, err := r.ceph.connection(scope)
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

	wg := sync.WaitGroup{}
	results := sync.Map{}

	for _, index := range indices {
		wg.Go(func() {
			path, clients, err := r.exportIndex(ioctx, index)
			if err != nil {
				return // skip error
			}
			results.Store(path, clients)
		})
	}

	wg.Wait()

	m := map[string][]string{}

	results.Range(func(key, value interface{}) bool {
		m[key.(string)] = value.([]string)
		return true
	})

	return m, nil
}

func (r *subvolumeRepo) exportIndice(ioctx *rados.IOContext) ([]string, error) {
	buffer := make([]byte, 1024) //nolint:mnd // default
	n, err := ioctx.Read("ganesha-export-index", buffer, 0)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(buffer[:n]), indexSeparator)

	return slices.DeleteFunc(lines, func(s string) bool { return s == "" }), nil
}

func (r *subvolumeRepo) exportIndex(ioctx *rados.IOContext, index string) (path string, clients []string, err error) {
	buffer := make([]byte, 1024) //nolint:mnd // default
	n, err := ioctx.Read(index, buffer, 0)
	if err != nil {
		return "", nil, err
	}

	if !validPath.Match(buffer[:n]) || !validClients.Match(buffer[:n]) {
		return "", nil, fmt.Errorf("export index %q not found", index)
	}

	path = validPath.FindStringSubmatch(string(buffer[:n]))[1]
	clients = []string{}

	for _, client := range strings.Split(validClients.FindStringSubmatch(string(buffer[:n]))[1], ",") {
		clients = append(clients, strings.TrimSpace(client))
	}

	return path, clients, nil
}*/

func (r *subvolumeRepo) toSubvolume(volume, group, subvolume string, info *subvolumeInfo) *file.Subvolume {
	if info == nil {
		return nil
	}

	used := uint64(info.BytesUsed)
	var quotaPtr *file.Bytes
	var quota uint64
	if info.BytesQuota != nil && uint64(*info.BytesQuota) > 0 {
		quota = uint64(*info.BytesQuota)
		q := file.Bytes(quota)
		quotaPtr = &q
	}

	var usage float64
	if quota > 0 {
		usage = float64(used) / float64(quota)
	}

	features := make([]file.Feature, 0, len(info.Features))
	for _, f := range info.Features {
		features = append(features, file.Feature(f))
	}

	return &file.Subvolume{
		Key: file.SubvolumeKey{
			VolumeName:    volume,
			GroupName:     group,
			SubvolumeName: subvolume,
		},
		Info: file.SubvolumeInfo{
			Path:         info.Path,
			State:        file.SubvolumeState(info.State),
			UID:          info.UID,
			GID:          info.GID,
			Mode:         file.UnixMode(info.Mode),
			BytesPercent: usage,
			BytesUsed:    file.Bytes(used),
			BytesQuota:   quotaPtr,

			DataPool:      info.DataPool,
			PoolNamespace: info.PoolNamespace,

			Atime:     info.Atime.Time,
			Mtime:     info.Mtime.Time,
			Ctime:     info.Ctime.Time,
			CreatedAt: info.CreatedAt.Time,

			Features: features,
		},
	}
}
