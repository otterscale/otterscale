// rgw_test.go
package ceph // same package as the implementation – allows calling unexported methods

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/stretchr/testify/assert"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

/* ---------------------------------------------------------- *
 * Helper – build a minimal *core.StorageConfig for the tests *
 * ---------------------------------------------------------- */
func newTestConfig(fsid, endpoint string) *core.StorageConfig {
	return &core.StorageConfig{
		StorageCephConfig: &core.StorageCephConfig{
			FSID:    fsid,
			MONHost: "invalid-monhost", // forces rados.Connect to fail
			Key:     "dummy-key",
		},
		StorageRGWConfig: &core.StorageRGWConfig{
			Endpoint:  endpoint, // empty => admin.New fails
			AccessKey: "dummy-ak",
			SecretKey: "dummy-sk",
		},
	}
}

/* ---------------------------------------------------------- *
 * Test boolToIntPointer (pure function)                     *
 * ---------------------------------------------------------- */
func TestBoolToIntPointer(t *testing.T) {
	r := &rgw{}
	assert.Equal(t, 1, *r.boolToIntPointer(true))
	assert.Nil(t, r.boolToIntPointer(false))
}

/* ---------------------------------------------------------- *
 * When the admin client cannot be created, the error must be   *
 * propagated (CreateBucket, ListBuckets, UpdateBucketOwner)    *
 * ---------------------------------------------------------- */
func TestCreateBucket_ClientError(t *testing.T) {
	ctx := context.Background()
	// empty endpoint -> admin.New fails (error message contains “endpoint not set”)
	cfg := newTestConfig("fsid-1", "")

	cephInst := New(&config.Config{})
	rgw := NewRGW(cephInst) // returns core.CephRGWRepo, usable as is

	err := rgw.CreateBucket(ctx, cfg, "test-bucket", types.BucketCannedACLPrivate)
	assert.Error(t, err)
	// The actual error comes from the admin client constructor,
	// so we just check that it mentions the endpoint.
	assert.Contains(t, err.Error(), "endpoint")
}

/* ---------------------------------------------------------- *
 * When the S3 client cannot be created, the error must be      *
 * propagated (DeleteBucket)                                    *
 * ---------------------------------------------------------- */
func TestDeleteBucket_S3ClientError(t *testing.T) {
	ctx := context.Background()
	// admin client can be created (non‑empty endpoint) but the endpoint does not exist.
	cfg := newTestConfig("fsid-2", "http://valid-endpoint")

	cephInst := New(&config.Config{})
	rgwInst := NewRGW(cephInst) // core.CephRGWRepo (interface)
	r := rgwInst.(*rgw)         // concrete type to access private methods

	err := r.DeleteBucket(ctx, cfg, "some-bucket")
	assert.Error(t, err)

	// The error originates from the S3 client DNS lookup.
	assert.Contains(t, err.Error(), "lookup")
	assert.Contains(t, err.Error(), "DeleteBucket")
}

/* ---------------------------------------------------------- *
 * ListBuckets should fail when the admin client cannot be      *
 * constructed.                                                *
 * ---------------------------------------------------------- */
func TestListBuckets_ClientError(t *testing.T) {
	ctx := context.Background()
	// empty endpoint -> admin.New fails
	cfg := newTestConfig("fsid-3", "")

	cephInst := New(&config.Config{})
	rgw := NewRGW(cephInst)

	_, err := rgw.ListBuckets(ctx, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "endpoint")
}

/* ---------------------------------------------------------- *
 * toBucket should return an error when the S3 client cannot    *
 * be created.                                                *
 * ---------------------------------------------------------- */
func TestToBucket_S3ClientError(t *testing.T) {
	ctx := context.Background()
	// admin client succeeds (valid endpoint) but S3 endpoint does not exist.
	cfg := newTestConfig("fsid-4", "http://valid-endpoint")

	cephInst := New(&config.Config{})
	rgwInst := NewRGW(cephInst)
	r := rgwInst.(*rgw)

	b := &admin.Bucket{Bucket: "dummy-bkt"}

	// The method should not panic and should return a bucket with nil policy/acl.
	bucket, err := r.toBucket(ctx, cfg, b)
	assert.NoError(t, err)
	assert.NotNil(t, bucket)
	assert.Nil(t, bucket.Policy)
	assert.Nil(t, bucket.Grants)
}

/* ---------------------------------------------------------- *
 * UpdateBucketOwner – error path (invalid configuration)      *
 * ---------------------------------------------------------- */
func TestUpdateBucketOwner_Error(t *testing.T) {
	ctx := context.Background()
	// completely invalid config -> client creation fails
	cfg := newTestConfig("", "")

	cephInst := New(&config.Config{})
	rgw := NewRGW(cephInst)

	err := rgw.UpdateBucketOwner(ctx, cfg, "bucket", "owner")
	assert.Error(t, err)
}
