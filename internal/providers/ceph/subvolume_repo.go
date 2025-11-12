package ceph

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/ceph/go-ceph/rados"
	"github.com/otterscale/otterscale/internal/core/storage/file"
)

const indexSeparator = "\n%url rados://ceph-nfs/"

var (
	validPath    = regexp.MustCompile(`Path = "(?<Path>.+)";`)
	validClients = regexp.MustCompile(`Clients = (?<Clients>.+);`)
)

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
	conn, err := r.ceph.Connection(scope)
	if err != nil {
		return nil, err
	}

	dumpNames, err := listSubvolumes(conn, volume, group)
	if err != nil {
		return nil, err
	}

	subvolumes := []file.Subvolume{}

	for _, dumpName := range dumpNames {
		info, err := getSubvolume(conn, volume, dumpName.Name, group)
		if err != nil {
			return nil, err
		}
		subvolumes = append(subvolumes, *r.toSubvolume(dumpName.Name, info))
	}

	return subvolumes, nil
}

func (r *subvolumeRepo) Get(_ context.Context, scope, volume, subvolume, group string) (*file.Subvolume, error) {
	conn, err := r.ceph.Connection(scope)
	if err != nil {
		return nil, err
	}

	dump, err := getSubvolume(conn, volume, subvolume, group)
	if err != nil {
		return nil, err
	}

	return r.toSubvolume(subvolume, dump), nil
}

func (r *subvolumeRepo) Create(_ context.Context, scope, volume, subvolume, group string, size uint64) error {
	conn, err := r.ceph.Connection(scope)
	if err != nil {
		return err
	}

	return createSubvolume(conn, volume, subvolume, group, size)
}

func (r *subvolumeRepo) Resize(_ context.Context, scope, volume, subvolume, group string, size uint64) error {
	conn, err := r.ceph.Connection(scope)
	if err != nil {
		return err
	}

	return resizeSubvolume(conn, volume, subvolume, group, size)
}

func (r *subvolumeRepo) Delete(_ context.Context, scope, volume, subvolume, group string) error {
	conn, err := r.ceph.Connection(scope)
	if err != nil {
		return err
	}

	return removeSubvolume(conn, volume, subvolume, group)
}

func (r *subvolumeRepo) ListExportClients(ctx context.Context, scope, pool string) (map[string][]string, error) {
	conn, err := r.ceph.Connection(scope)
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
}

func (r *subvolumeRepo) toSubvolume(name string, info *subvolumeInfo) *file.Subvolume {
	// quota, _ := parseQuota(info.BytesQuota)
	ret := &file.Subvolume{
		// Name:      name,
		// Path:      info.Path,
		// Mode:      fmt.Sprintf("%o", info.Mode),
		// PoolName:  info.DataPool,
		// Quota:     quota,
		// Used:      info.BytesUsed,
		// CreatedAt: info.CreatedAt.Time,
	}
	return ret
}
