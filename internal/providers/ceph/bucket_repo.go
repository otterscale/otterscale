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
	"github.com/otterscale/otterscale/internal/core/storage/object"
)

type bucketRepo struct {
	ceph *Ceph
}

func NewBucketRepo(ceph *Ceph) object.BucketRepo {
	return &bucketRepo{
		ceph: ceph,
	}
}

var _ object.BucketRepo = (*bucketRepo)(nil)

func (r *bucketRepo) List(ctx context.Context, scope string) ([]object.Bucket, error) {
	client, err := r.ceph.client(scope)
	if err != nil {
		return nil, err
	}

	return client.ListBucketsWithStat(ctx)
}

func (r *bucketRepo) Get(ctx context.Context, scope, bucket string) (*object.Bucket, error) {
	buckets, err := r.List(ctx, scope)
	if err != nil {
		return nil, err
	}

	for i := range buckets {
		if buckets[i].Bucket == bucket {
			return &buckets[i], nil
		}
	}

	return nil, fmt.Errorf("bucket %q not found", bucket)
}

func (r *bucketRepo) Create(ctx context.Context, scope, bucket string, acl object.BucketCannedACL) error {
	s3Client, err := r.s3Client(ctx, scope)
	if err != nil {
		return err
	}

	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucket,
		ACL:    types.BucketCannedACL(acl),
	})
	return err
}

func (r *bucketRepo) Delete(ctx context.Context, scope, bucket string) error {
	s3Client, err := r.s3Client(ctx, scope)
	if err != nil {
		return err
	}

	_, err = s3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: &bucket,
	})
	return err
}

func (r *bucketRepo) UpdateOwner(ctx context.Context, scope, bucket, owner string) error {
	client, err := r.ceph.client(scope)
	if err != nil {
		return err
	}

	return client.LinkBucket(ctx, admin.BucketLinkInput{
		Bucket: bucket,
		UID:    owner,
	})
}

func (r *bucketRepo) GetPolicy(ctx context.Context, scope, bucket string) (*string, error) {
	s3Client, err := r.s3Client(ctx, scope)
	if err != nil {
		return nil, err
	}

	resp, err := s3Client.GetBucketPolicy(ctx, &s3.GetBucketPolicyInput{
		Bucket: &bucket,
	})
	if err != nil {
		return nil, err
	}

	return resp.Policy, nil
}

func (r *bucketRepo) UpdatePolicy(ctx context.Context, scope, bucket, policy string) error {
	s3Client, err := r.s3Client(ctx, scope)
	if err != nil {
		return err
	}

	_, err = s3Client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: &bucket,
		Policy: &policy,
	})
	return err
}

func (r *bucketRepo) GetACL(ctx context.Context, scope, bucket string) ([]object.Grant, error) {
	s3Client, err := r.s3Client(ctx, scope)
	if err != nil {
		return nil, err
	}

	resp, err := s3Client.GetBucketAcl(ctx, &s3.GetBucketAclInput{
		Bucket: &bucket,
	})
	if err != nil {
		return nil, err
	}

	return resp.Grants, nil
}

func (r *bucketRepo) UpdateACL(ctx context.Context, scope, bucket string, acl object.BucketCannedACL) error {
	s3Client, err := r.s3Client(ctx, scope)
	if err != nil {
		return err
	}

	_, err = s3Client.PutBucketAcl(ctx, &s3.PutBucketAclInput{
		Bucket: &bucket,
		ACL:    types.BucketCannedACL(acl),
	})
	return err
}

func (r *bucketRepo) Endpoint(scope string) string {
	client, err := r.ceph.client(scope)
	if err != nil {
		return ""
	}

	return client.Endpoint
}

func (r *bucketRepo) Key(scope string) (accessKey string, secretKey string) {
	client, err := r.ceph.client(scope)
	if err != nil {
		return "", ""
	}

	return client.AccessKey, client.SecretKey
}

func (r *bucketRepo) s3Client(ctx context.Context, scope string) (*s3.Client, error) {
	const (
		cephRegion       = "us-east-1"
		retryMaxAttempts = 5
	)

	client, err := r.ceph.client(scope)
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

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	}), nil
}
