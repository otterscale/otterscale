package object

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ceph/go-ceph/rgw/admin"
)

const (
	BucketCannedACLPrivate           = types.BucketCannedACLPrivate
	BucketCannedACLPublicRead        = types.BucketCannedACLPublicRead
	BucketCannedACLPublicReadWrite   = types.BucketCannedACLPublicReadWrite
	BucketCannedACLAuthenticatedRead = types.BucketCannedACLAuthenticatedRead
)

type (
	// Bucket represents a Ceph RGW Bucket resource.
	Bucket = admin.Bucket

	// BucketCannedACL represents a AWS BucketCannedACL resource.
	BucketCannedACL = types.BucketCannedACL

	// Grant represents a AWS Grant resource.
	Grant = types.Grant
)

type BucketData struct {
	*Bucket
	Policy *string
	Grants []Grant
}

// Note: Ceph create and update operations only return error status.
type BucketRepo interface {
	List(ctx context.Context, scope string) ([]Bucket, error)
	Get(ctx context.Context, scope, bucket string) (*Bucket, error)
	Create(ctx context.Context, scope, bucket string, acl BucketCannedACL) error
	Delete(ctx context.Context, scope, bucket string) error
	UpdateOwner(ctx context.Context, scope, bucket, owner string) error
	GetPolicy(ctx context.Context, scope, bucket string) (*string, error)
	UpdatePolicy(ctx context.Context, scope, bucket, policy string) error
	GetACL(ctx context.Context, scope, bucket string) ([]Grant, error)
	UpdateACL(ctx context.Context, scope, bucket string, acl BucketCannedACL) error
	Endpoint(scope string) string
	Key(scope string) (accessKey string, secretKey string)
}

func (uc *UseCase) ListBuckets(ctx context.Context, scope string) ([]BucketData, error) {
	buckets, err := uc.bucket.List(ctx, scope)
	if err != nil {
		return nil, err
	}

	ret := []BucketData{}

	for i := range buckets {
		policy, _ := uc.bucket.GetPolicy(ctx, scope, buckets[i].Bucket)
		grants, _ := uc.bucket.GetACL(ctx, scope, buckets[i].Bucket)

		ret = append(ret, BucketData{
			Bucket: &buckets[i],
			Policy: policy,
			Grants: grants,
		})
	}

	return ret, nil
}

func (uc *UseCase) CreateBucket(ctx context.Context, scope, bucket, owner, policy string, acl BucketCannedACL) (*BucketData, error) {
	if err := uc.bucket.Create(ctx, scope, bucket, acl); err != nil {
		return nil, err
	}

	if err := uc.bucket.UpdateOwner(ctx, scope, bucket, owner); err != nil {
		return nil, err
	}

	if policy != "" {
		if err := uc.bucket.UpdatePolicy(ctx, scope, bucket, policy); err != nil {
			return nil, err
		}
	}

	if err := uc.bucket.UpdateACL(ctx, scope, bucket, acl); err != nil {
		return nil, err
	}

	b, err := uc.bucket.Get(ctx, scope, bucket)
	if err != nil {
		return nil, err
	}

	p, err := uc.bucket.GetPolicy(ctx, scope, bucket)
	if err != nil {
		return nil, err
	}

	a, err := uc.bucket.GetACL(ctx, scope, bucket)
	if err != nil {
		return nil, err
	}

	return &BucketData{
		Bucket: b,
		Policy: p,
		Grants: a,
	}, nil
}

func (uc *UseCase) UpdateBucket(ctx context.Context, scope, bucket, owner, policy string, acl BucketCannedACL) (*BucketData, error) {
	if err := uc.bucket.UpdateOwner(ctx, scope, bucket, owner); err != nil {
		return nil, err
	}

	if policy != "" {
		if err := uc.bucket.UpdatePolicy(ctx, scope, bucket, policy); err != nil {
			return nil, err
		}
	}

	if err := uc.bucket.UpdateACL(ctx, scope, bucket, acl); err != nil {
		return nil, err
	}

	b, err := uc.bucket.Get(ctx, scope, bucket)
	if err != nil {
		return nil, err
	}

	p, err := uc.bucket.GetPolicy(ctx, scope, bucket)
	if err != nil {
		return nil, err
	}

	a, err := uc.bucket.GetACL(ctx, scope, bucket)
	if err != nil {
		return nil, err
	}

	return &BucketData{
		Bucket: b,
		Policy: p,
		Grants: a,
	}, nil
}

func (uc *UseCase) DeleteBucket(ctx context.Context, scope, bucket string) error {
	return uc.bucket.Delete(ctx, scope, bucket)
}
