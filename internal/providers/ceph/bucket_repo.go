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
	client, err := r.ceph.Client(scope)
	if err != nil {
		return nil, err
	}

	statBuckets, err := client.ListBucketsWithStat(ctx)
	if err != nil {
		return nil, err
	}

	buckets := []object.Bucket{}

	for i := range statBuckets {
		bucket, err := r.toBucket(ctx, scope, &statBuckets[i])
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, *bucket)
	}

	return buckets, nil
}

func (r *bucketRepo) Get(ctx context.Context, scope, bucket string) (*object.Bucket, error) {
	buckets, err := r.List(ctx, scope)
	if err != nil {
		return nil, err
	}

	for i := range buckets {
		if buckets[i].Name == bucket {
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

func (r *bucketRepo) UpdateOwner(ctx context.Context, scope, bucket, owner string) error {
	client, err := r.ceph.Client(scope)
	if err != nil {
		return err
	}

	return client.LinkBucket(ctx, admin.BucketLinkInput{
		Bucket: bucket,
		UID:    owner,
	})
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

func (r *bucketRepo) s3Client(ctx context.Context, scope string) (*s3.Client, error) {
	const (
		cephRegion       = "us-east-1"
		retryMaxAttempts = 5
	)

	client, err := r.ceph.Client(scope)
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

func (r *bucketRepo) toBucket(ctx context.Context, scope string, b *admin.Bucket) (*object.Bucket, error) {
	bucket := &object.Bucket{
		Name: b.Bucket,
	}

	s3Client, err := r.s3Client(ctx, scope)
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
		bucket.Grants = r.toGrants(acl.Grants)
	}

	return bucket, nil
}

func (r *bucketRepo) toGrants(gs []types.Grant) []object.BucketGrant {
	grants := []object.BucketGrant{}

	for _, g := range gs {
		grants = append(grants, object.BucketGrant{
			Grantee: &object.Grantee{
				Type:         object.Type(g.Grantee.Type),
				DisplayName:  g.Grantee.DisplayName,
				EmailAddress: g.Grantee.EmailAddress,
				ID:           g.Grantee.ID,
				URI:          g.Grantee.URI,
			},
			Permission: object.Permission(g.Permission),
		})
	}

	return grants
}
