package ceph

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ceph/go-ceph/rgw/admin"

	"github.com/otterscale/otterscale/internal/core"
)

type rgw struct {
	ceph *Ceph
}

func NewRGW(ceph *Ceph) core.CephRGWRepo {
	return &rgw{
		ceph: ceph,
	}
}

var _ core.CephRGWRepo = (*rgw)(nil)

// ceph api
func (r *rgw) ListBuckets(ctx context.Context, config *core.CephConfig) ([]core.RGWBucket, error) {
	client, err := r.ceph.client(config)
	if err != nil {
		return nil, err
	}
	bs, err := client.ListBucketsWithStat(ctx)
	if err != nil {
		return nil, err
	}
	buckets := []core.RGWBucket{}
	for i := range bs {
		bucket, err := r.toBucket(ctx, config, &bs[i])
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, *bucket)
	}
	return buckets, nil
}

// ceph api
func (r *rgw) GetBucket(ctx context.Context, config *core.CephConfig, bucket string) (*core.RGWBucket, error) {
	client, err := r.ceph.client(config)
	if err != nil {
		return nil, err
	}
	bs, err := client.ListBucketsWithStat(ctx)
	if err != nil {
		return nil, err
	}
	for i := range bs {
		if bs[i].Bucket == bucket {
			return r.toBucket(ctx, config, &bs[i])
		}
	}
	return nil, fmt.Errorf("bucket %q not found", bucket)
}

// s3 api
func (r *rgw) CreateBucket(ctx context.Context, config *core.CephConfig, bucket string, acl types.BucketCannedACL) error {
	s3Client, err := r.s3Client(ctx, config)
	if err != nil {
		return err
	}
	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucket,
		ACL:    acl,
	})
	return err
}

// ceph api
func (r *rgw) UpdateBucketOwner(ctx context.Context, config *core.CephConfig, bucket, owner string) error {
	client, err := r.ceph.client(config)
	if err != nil {
		return err
	}
	return client.LinkBucket(ctx, admin.BucketLinkInput{
		Bucket: bucket,
		UID:    owner,
	})
}

// s3 api
func (r *rgw) UpdateBucketACL(ctx context.Context, config *core.CephConfig, bucket string, acl types.BucketCannedACL) error {
	s3Client, err := r.s3Client(ctx, config)
	if err != nil {
		return err
	}
	_, err = s3Client.PutBucketAcl(ctx, &s3.PutBucketAclInput{
		Bucket: &bucket,
		ACL:    acl,
	})
	return err
}

// s3 api
func (r *rgw) UpdateBucketPolicy(ctx context.Context, config *core.CephConfig, bucket, policy string) error {
	s3Client, err := r.s3Client(ctx, config)
	if err != nil {
		return err
	}
	_, err = s3Client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: &bucket,
		Policy: &policy,
	})
	return err
}

// s3 api
func (r *rgw) DeleteBucket(ctx context.Context, config *core.CephConfig, bucket string) error {
	s3Client, err := r.s3Client(ctx, config)
	if err != nil {
		return err
	}
	_, err = s3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: &bucket,
	})
	return err
}

// ceph api
func (r *rgw) ListUsers(ctx context.Context, config *core.CephConfig) ([]core.RGWUser, error) {
	client, err := r.ceph.client(config)
	if err != nil {
		return nil, err
	}
	ids, err := client.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	if ids == nil {
		return nil, fmt.Errorf("empty user ids")
	}
	users := []core.RGWUser{}
	for _, id := range *ids {
		user, err := client.GetUser(ctx, admin.User{
			ID: id,
			Keys: []admin.UserKeySpec{{
				AccessKey: config.AccessKey,
				SecretKey: config.SecretKey,
			}},
		})
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// ceph api
func (r *rgw) CreateUser(ctx context.Context, config *core.CephConfig, id, name string, suspended bool) (*core.RGWUser, error) {
	client, err := r.ceph.client(config)
	if err != nil {
		return nil, err
	}
	user, err := client.CreateUser(ctx, admin.User{
		ID:          id,
		DisplayName: name,
		Suspended:   r.boolToIntPointer(suspended),
	})
	if err != nil {
		return nil, err
	}
	return &user, err
}

// ceph api
func (r *rgw) UpdateUser(ctx context.Context, config *core.CephConfig, id, name string, suspended bool) (*core.RGWUser, error) {
	client, err := r.ceph.client(config)
	if err != nil {
		return nil, err
	}
	user, err := client.ModifyUser(ctx, admin.User{
		ID:          id,
		DisplayName: name,
		Suspended:   r.boolToIntPointer(suspended),
	})
	if err != nil {
		return nil, err
	}
	return &user, err
}

// ceph api
func (r *rgw) DeleteUser(ctx context.Context, config *core.CephConfig, id string) error {
	client, err := r.ceph.client(config)
	if err != nil {
		return err
	}
	return client.RemoveUser(ctx, admin.User{ID: id})
}

// ceph api
func (r *rgw) CreateUserKey(ctx context.Context, config *core.CephConfig, id string) (*core.RGWUserKey, error) {
	client, err := r.ceph.client(config)
	if err != nil {
		return nil, err
	}
	keys, err := client.CreateKey(ctx, admin.UserKeySpec{UID: id})
	if err != nil {
		return nil, err
	}
	if keys != nil && len(*keys) > 0 {
		return &(*keys)[0], nil
	}
	return nil, fmt.Errorf("create key failed")
}

// ceph api
func (r *rgw) DeleteUserKey(ctx context.Context, config *core.CephConfig, id, accessKey string) error {
	client, err := r.ceph.client(config)
	if err != nil {
		return err
	}
	return client.RemoveKey(ctx, admin.UserKeySpec{UID: id, AccessKey: accessKey})
}

func (r *rgw) s3Client(ctx context.Context, conf *core.CephConfig) (*s3.Client, error) {
	const (
		cephRegion       = "us-east-1"
		retryMaxAttempts = 5
	)

	client, err := r.ceph.client(conf)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cephRegion),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     client.AccessKey,
				SecretAccessKey: client.SecretKey,
				SessionToken:    "",
			},
		}),
		config.WithBaseEndpoint(client.Endpoint),
		config.WithRetryMaxAttempts(retryMaxAttempts),
		config.WithHTTPClient(client.HTTPClient),
	)
	if err != nil {
		return nil, err
	}

	svc := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	return svc, nil
}

func (r *rgw) toBucket(ctx context.Context, config *core.CephConfig, b *admin.Bucket) (*core.RGWBucket, error) {
	bucket := &core.RGWBucket{
		Bucket: b,
	}
	s3Client, err := r.s3Client(ctx, config)
	if err != nil {
		return nil, err
	}
	policy, _ := s3Client.GetBucketPolicy(ctx, &s3.GetBucketPolicyInput{
		Bucket: &b.Bucket,
	})
	if policy != nil {
		bucket.Policy = policy.Policy
	}
	acl, _ := s3Client.GetBucketAcl(ctx, &s3.GetBucketAclInput{
		Bucket: &b.Bucket,
	})
	if acl != nil {
		bucket.Grants = acl.Grants
	}
	return bucket, nil
}

func (r *rgw) boolToIntPointer(b bool) *int {
	if b {
		return aws.Int(1)
	}
	return nil
}
