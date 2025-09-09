// rgb_test.go
package ceph // same package as the implementation – allows calling unexported methods

import (
	"context"
	"testing"
	"time"

	"github.com/ceph/go-ceph/rados"
	"github.com/stretchr/testify/assert"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

/* ---------------------------------------------------------- *
 * Helper – build a minimal *core.StorageConfig for the tests *
 * ---------------------------------------------------------- */
func rbd_newTestConfig(fsid, endpoint string) *core.StorageConfig {
	return &core.StorageConfig{
		StorageCephConfig: &core.StorageCephConfig{
			FSID:    fsid,
			MONHost: "invalid-monhost", // forces rados.Connect to fail
			Key:     "dummy-key",
		},
		StorageRGWConfig: &core.StorageRGWConfig{
			Endpoint:  endpoint, // not used by RBD, can be empty
			AccessKey: "dummy-ak",
			SecretKey: "dummy-sk",
		},
	}
}

/* ---------------------------------------------------------- *
 * Test Bool functions – pure helpers                               *
 * ---------------------------------------------------------- */
func TestFeatureOn(t *testing.T) {
	r := &rbd{}
	const (
		featA uint64 = 0x1
		featB uint64 = 0x2
	)
	assert.True(t, r.featureOn(featA|featB, featA))
	assert.False(t, r.featureOn(featA, featB))
}

/* ---------------------------------------------------------- *
 * Every public method ultimately calls r.ceph.connection().
 * By providing an invalid MONHost the connection fails, so the
 * method must return an error without panic.
 * ---------------------------------------------------------- */
func TestListImages_ConnectionError(t *testing.T) {
	ctx := context.Background()
	cfg := rbd_newTestConfig("fsid-1", "")
	cephInst := New(&config.Config{})
	rbdInst := NewRBD(cephInst)

	_, err := rbdInst.ListImages(ctx, cfg, "any-pool")
	assert.Error(t, err)
}

func TestGetImage_ConnectionError(t *testing.T) {
	ctx := context.Background()
	cfg := rbd_newTestConfig("fsid-2", "")
	cephInst := New(&config.Config{})
	rbdInst := NewRBD(cephInst)

	_, err := rbdInst.GetImage(ctx, cfg, "any-pool", "img")
	assert.Error(t, err)
}

func TestCreateImage_ConnectionError(t *testing.T) {
	ctx := context.Background()
	cfg := rbd_newTestConfig("fsid-3", "")
	cephInst := New(&config.Config{})
	rbdInst := NewRBD(cephInst)

	_, err := rbdInst.CreateImage(ctx, cfg, "any-pool", "img", 0, 0, 0, 0, 0)
	assert.Error(t, err)
}

func TestUpdateImageSize_ConnectionError(t *testing.T) {
	ctx := context.Background()
	cfg := rbd_newTestConfig("fsid-4", "")
	cephInst := New(&config.Config{})
	rbdInst := NewRBD(cephInst)

	err := rbdInst.UpdateImageSize(ctx, cfg, "any-pool", "img", 1024)
	assert.Error(t, err)
}

func TestDeleteImage_ConnectionError(t *testing.T) {
	ctx := context.Background()
	cfg := rbd_newTestConfig("fsid-5", "")
	cephInst := New(&config.Config{})
	rbdInst := NewRBD(cephInst)

	err := rbdInst.DeleteImage(ctx, cfg, "any-pool", "img")
	assert.Error(t, err)
}

/* ---------------------------------------------------------- *
 * Snapshot‑related methods – also depend on a valid connection *
 * ---------------------------------------------------------- */
func TestCreateImageSnapshot_ConnectionError(t *testing.T) {
	ctx := context.Background()
	cfg := rbd_newTestConfig("fsid-6", "")
	cephInst := New(&config.Config{})
	rbdInst := NewRBD(cephInst)

	err := rbdInst.CreateImageSnapshot(ctx, cfg, "any-pool", "img", "snap")
	assert.Error(t, err)
}

func TestDeleteImageSnapshot_ConnectionError(t *testing.T) {
	ctx := context.Background()
	cfg := rbd_newTestConfig("fsid-7", "")
	cephInst := New(&config.Config{})
	rbdInst := NewRBD(cephInst)

	err := rbdInst.DeleteImageSnapshot(ctx, cfg, "any-pool", "img", "snap")
	assert.Error(t, err)
}

func TestRollbackImageSnapshot_ConnectionError(t *testing.T) {
	ctx := context.Background()
	cfg := rbd_newTestConfig("fsid-8", "")
	cephInst := New(&config.Config{})
	rbdInst := NewRBD(cephInst)

	err := rbdInst.RollbackImageSnapshot(ctx, cfg, "any-pool", "img", "snap")
	assert.Error(t, err)
}

func TestProtectImageSnapshot_ConnectionError(t *testing.T) {
	ctx := context.Background()
	cfg := rbd_newTestConfig("fsid-9", "")
	cephInst := New(&config.Config{})
	rbdInst := NewRBD(cephInst)

	err := rbdInst.ProtectImageSnapshot(ctx, cfg, "any-pool", "img", "snap")
	assert.Error(t, err)
}

func TestUnprotectImageSnapshot_ConnectionError(t *testing.T) {
	ctx := context.Background()
	cfg := rbd_newTestConfig("fsid-10", "")
	cephInst := New(&config.Config{})
	rbdInst := NewRBD(cephInst)

	err := rbdInst.UnprotectImageSnapshot(ctx, cfg, "any-pool", "img", "snap")
	assert.Error(t, err)
}

/* ---------------------------------------------------------- *
 * The private method `openImage` is exercised indirectly
 * via the public methods.  To prove it does not panic when the
 * underlying Ceph calls fail, we call it directly after forcing
 * a connection error.
 * ---------------------------------------------------------- */
func TestOpenImage_ErrorPropagation(t *testing.T) {
	ctx := context.Background()
	cfg := rbd_newTestConfig("fsid-11", "")
	cephInst := New(&config.Config{})
	rbdInst := NewRBD(cephInst)

	// Obtain a connection that will fail on OpenIOContext
	_, err := rbdInst.GetImage(ctx, cfg, "any-pool", "non‑existent")
	assert.Error(t, err)

	// Direct call to the unexported helper - we need the concrete type.
	r := rbdInst.(*rbd)

	// Build a dummy ioctx – it will be nil because connection failed.
	var ioctx *rados.IOContext = nil
	_, err = r.openImage(ioctx, "any-pool", "non‑existent")
	assert.Error(t, err)
}

/* ---------------------------------------------------------- *
 * Verify that a successfully built RBDImage contains the
 * expected fields (only a few sanity checks – the heavy
 * Ceph interaction is not exercised in unit tests).
 * ---------------------------------------------------------- */
func TestRBDImageStruct(t *testing.T) {
	// Construct a minimal RBDImage manually – this does not hit Ceph.
	img := core.RBDImage{
		Name:                 "img1",
		PoolName:             "pool1",
		ObjectSize:           4096,
		StripeUnit:           0,
		StripeCount:          0,
		Quota:                10 << 30, // 10 GiB
		Used:                 5 << 30,
		ObjectCount:          2560,
		FeatureLayering:      true,
		FeatureExclusiveLock: false,
		FeatureObjectMap:     true,
		FeatureFastDiff:      false,
		FeatureDeepFlatten:   true,
		CreatedAt:            time.Now(),
		Snapshots: []core.RBDImageSnapshot{
			{
				Name:      "snap1",
				Quota:     5 << 30,
				Used:      2 << 30,
				Protected: false,
			},
		},
	}
	assert.Equal(t, "img1", img.Name)
	assert.Equal(t, "pool1", img.PoolName)
	assert.Len(t, img.Snapshots, 1)
	assert.True(t, img.FeatureLayering)
}
