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

	"github.com/openhdc/otterscale/internal/core"
)

const cephRegion = "us-east-1"

type rgw struct {
	ceph *Ceph
}

func NewRGW(ceph *Ceph) core.CephRGWRepo {
	return &rgw{
		ceph: ceph,
	}
}

var _ core.CephRGWRepo = (*rgw)(nil)

func (r *rgw) ListBuckets(ctx context.Context, config *core.StorageConfig) ([]core.RGWBucket, error) {
	client, err := r.ceph.client(config)
	if err != nil {
		return nil, err
	}
	return client.ListBucketsWithStat(ctx)
}

func (r *rgw) GetBucket(ctx context.Context, config *core.StorageConfig, bucket string) (*core.RGWBucket, error) {
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
			return &bs[i], nil
		}
	}
	return nil, fmt.Errorf("bucket %q not found", bucket)
}

// TODO: AuthN & AuthZ
func (r *rgw) CreateBucket(ctx context.Context, config *core.StorageConfig, bucket string, acl types.BucketCannedACL) error {
	s3Client, err := r.s3Client(config)
	if err != nil {
		return err
	}
	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucket,
		ACL:    acl,
	})
	return err
}

// TODO: AuthN & AuthZ
func (r *rgw) UpdateBucketOwner(ctx context.Context, config *core.StorageConfig, bucket, owner string) error {
	client, err := r.ceph.client(config)
	if err != nil {
		return err
	}
	return client.LinkBucket(ctx, admin.BucketLinkInput{
		Bucket: bucket,
		UID:    owner,
	})
}

// TODO: AuthN & AuthZ
func (r *rgw) UpdateBucketACL(ctx context.Context, config *core.StorageConfig, bucket string, acl types.BucketCannedACL) error {
	s3Client, err := r.s3Client(config)
	if err != nil {
		return err
	}
	_, err = s3Client.PutBucketAcl(ctx, &s3.PutBucketAclInput{
		Bucket: &bucket,
		ACL:    acl,
	})
	return err
}

// TODO: AuthN & AuthZ
func (r *rgw) UpdateBucketPolicy(ctx context.Context, config *core.StorageConfig, bucket, policy string) error {
	s3Client, err := r.s3Client(config)
	if err != nil {
		return err
	}
	_, err = s3Client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: &bucket,
		Policy: &policy,
	})
	return err
}

// TODO: AuthN & AuthZ
func (r *rgw) DeleteBucket(ctx context.Context, config *core.StorageConfig, bucket string) error {
	s3Client, err := r.s3Client(config)
	if err != nil {
		return err
	}
	_, err = s3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: &bucket,
	})
	return err
}

func (r *rgw) ListUsers(ctx context.Context, config *core.StorageConfig) ([]core.RGWUser, error) {
	client, err := r.ceph.client(config)
	if err != nil {
		return nil, err
	}
	users := []core.RGWUser{}
	userNames, err := client.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	if userNames != nil {
		for _, userName := range *userNames {
			users = append(users, core.RGWUser{
				DisplayName: userName,
			})
		}
	}
	return users, nil
}

// TODO: AuthN & AuthZ
func (r *rgw) CreateUser(ctx context.Context, config *core.StorageConfig, id, name string, suspended bool) (*core.RGWUser, error) {
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

// TODO: AuthN & AuthZ
func (r *rgw) UpdateUser(ctx context.Context, config *core.StorageConfig, id, name string, suspended bool) (*core.RGWUser, error) {
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

// TODO: AuthN & AuthZ
func (r *rgw) DeleteUser(ctx context.Context, config *core.StorageConfig, id string) error {
	client, err := r.ceph.client(config)
	if err != nil {
		return err
	}
	return client.RemoveUser(ctx, admin.User{ID: id})
}

// TODO: AuthN & AuthZ
func (r *rgw) CreateUserKey(ctx context.Context, config *core.StorageConfig, id string) (*core.RGWUserKey, error) {
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

// TODO: AuthN & AuthZ
func (r *rgw) DeleteUserKey(ctx context.Context, config *core.StorageConfig, id, accessKey string) error {
	client, err := r.ceph.client(config)
	if err != nil {
		return err
	}
	return client.RemoveKey(ctx, admin.UserKeySpec{UID: id, AccessKey: accessKey})
}

func (r *rgw) s3Client(conf *core.StorageConfig) (*s3.Client, error) {
	client, err := r.ceph.client(conf)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cephRegion),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     client.AccessKey,
				SecretAccessKey: client.SecretKey,
				SessionToken:    "",
			},
		}),
		config.WithBaseEndpoint(client.Endpoint),
		config.WithRetryMaxAttempts(5),
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

func (r *rgw) boolToIntPointer(b bool) *int {
	if b {
		i := 1
		return &i
	}
	return nil
}
