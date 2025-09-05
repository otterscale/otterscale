// fs_test.go
package ceph // same package as the implementation – allows calling unexported methods

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

/* ----------------------------------------------------------
 * Helper – build a minimal *core.StorageConfig for the tests
 * ---------------------------------------------------------- */
func newFSTestConfig(fsid, endpoint string) *core.StorageConfig {
	return &core.StorageConfig{
		StorageCephConfig: &core.StorageCephConfig{
			FSID:    fsid,
			MONHost: "invalid-monhost", // forces rados.Connect to fail
			Key:     "dummy-key",
		},
		StorageRGWConfig: &core.StorageRGWConfig{
			Endpoint:  endpoint,
			AccessKey: "dummy-ak",
			SecretKey: "dummy-sk",
		},
	}
}

/* ----------------------------------------------------------
 * 1️⃣  Connection‑error path – every public method ends up
 *     calling r.ceph.connection(), therefore an invalid MONHost
 *     makes the method return an error without panicking.
 * ---------------------------------------------------------- */
func TestListVolumes_ConnectionError(t *testing.T) {
	ctx := context.Background()
	cfg := newFSTestConfig("fsid-1", "")
	cephInst := New(&config.Config{})
	fsRepo := NewFS(cephInst)

	_, err := fsRepo.ListVolumes(ctx, cfg)
	assert.Error(t, err)
}

func TestListSubvolumes_ConnectionError(t *testing.T) {
	ctx := context.Background()
	cfg := newFSTestConfig("fsid-2", "")
	cephInst := New(&config.Config{})
	fsRepo := NewFS(cephInst)

	_, err := fsRepo.ListSubvolumes(ctx, cfg, "vol", "grp")
	assert.Error(t, err)
}

func TestGetSubvolume_ConnectionError(t *testing.T) {
	ctx := context.Background()
	cfg := newFSTestConfig("fsid-3", "")
	cephInst := New(&config.Config{})
	fsRepo := NewFS(cephInst)

	_, err := fsRepo.GetSubvolume(ctx, cfg, "vol", "sub", "grp")
	assert.Error(t, err)
}

func TestCreateSubvolume_ConnectionError(t *testing.T) {
	ctx := context.Background()
	cfg := newFSTestConfig("fsid-4", "")
	cephInst := New(&config.Config{})
	fsRepo := NewFS(cephInst)

	err := fsRepo.CreateSubvolume(ctx, cfg, "vol", "sub", "grp", 1<<30)
	assert.Error(t, err)
}

func TestResizeSubvolume_ConnectionError(t *testing.T) {
	ctx := context.Background()
	cfg := newFSTestConfig("fsid-5", "")
	cephInst := New(&config.Config{})
	fsRepo := NewFS(cephInst)

	err := fsRepo.ResizeSubvolume(ctx, cfg, "vol", "sub", "grp", 2<<30)
	assert.Error(t, err)
}

func TestDeleteSubvolume_ConnectionError(t *testing.T) {
	ctx := context.Background()
	cfg := newFSTestConfig("fsid-6", "")
	cephInst := New(&config.Config{})
	fsRepo := NewFS(cephInst)

	err := fsRepo.DeleteSubvolume(ctx, cfg, "vol", "sub", "grp")
	assert.Error(t, err)
}

/* toVolumes ---------------------------------------------------- */
func TestToVolumes(t *testing.T) {
	now := time.Now()

	d := &fsDump{
		FileSystems: []struct {
			MDSMap struct {
				FileSystemName string   `json:"fs_name,omitempty"`
				Created        CephTime `json:"created,omitempty"`
			} `json:"mdsmap,omitempty"`
			ID int64 `json:"id,omitempty"`
		}{
			{
				MDSMap: struct {
					FileSystemName string   `json:"fs_name,omitempty"`
					Created        CephTime `json:"created,omitempty"`
				}{
					FileSystemName: "myvol",
					Created:        CephTime{Time: now},
				},
				ID: 123,
			},
		},
	}
	r := &fs{}
	vols := r.toVolumes(d)
	// Validate the result.
	assert.Len(t, vols, 1)
	assert.Equal(t, int64(123), vols[0].ID)
	assert.Equal(t, "myvol", vols[0].Name)
	assert.Equal(t, now, vols[0].CreatedAt)
}

func TestToSubvolume(t *testing.T) {
	created := time.Now()
	info := &subvolumeInfo{
		BytesQuota: "10737418240", // 10 GiB
		Path:       "/my/vol/sub",
		Mode:       0o755,
		DataPool:   "rbdpool",
		BytesUsed:  5 << 30,
		CreatedAt: CephSubvolumeTime{
			Time: created,
		},
	}
	r := &fs{}
	sub := r.toSubvolume("mysub", info)
	quota, _ := strconv.ParseUint("10737418240", 10, 64)

	assert.Equal(t, "mysub", sub.Name)
	assert.Equal(t, "/my/vol/sub", sub.Path)
	assert.Equal(t, "755", sub.Mode) // octal string
	assert.Equal(t, "rbdpool", sub.PoolName)
	assert.Equal(t, quota, sub.Quota)
	assert.Equal(t, uint64(5<<30), sub.Used)
	assert.Equal(t, created, sub.CreatedAt)
}

func TestToSubvolumeSnapshot(t *testing.T) {
	now := time.Now()
	info := &subvolumeSnapshotInfo{
		HasPendingClones: "false",
		CreatedAt: CephSubvolumeTime{
			Time: now,
		},
	}
	r := &fs{}
	snap := r.toSubvolumeSnapshot("snap1", info)

	assert.Equal(t, "snap1", snap.Name)
	assert.Equal(t, "false", snap.HasPendingClones)
	assert.Equal(t, now, snap.CreatedAt)
}

func TestToSubvolumeGroups(t *testing.T) {
	created := time.Now()
	info := &subvolumeGroupInfo{
		BytesQuota: 21474836480,
		Mode:       0o770, // octal literal – 504 decimal
		DataPool:   "data-pool",
		BytesUsed:  12 << 30, // 12 GiB
		CreatedAt: CephSubvolumeTime{
			Time: created,
		},
	}

	r := &fs{}
	group := r.toSubvolumeGroups("mygroup", info)
	quota, _ := strconv.ParseUint("21474836480", 10, 64)
	assert.Equal(t, "mygroup", group.Name)
	assert.Equal(t, "000770", group.Mode)
	assert.Equal(t, "data-pool", group.PoolName)
	assert.Equal(t, quota, group.Quota)
	assert.Equal(t, uint64(12<<30), group.Used)
	assert.Equal(t, created, group.CreatedAt)
}
