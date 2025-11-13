package object

import (
	"context"
)

type Type string

const (
	TypeCanonicalUser         Type = "CanonicalUser"
	TypeAmazonCustomerByEmail Type = "AmazonCustomerByEmail"
	TypeGroup                 Type = "Group"
)

type Permission string

const (
	PermissionFullControl Permission = "FULL_CONTROL"
	PermissionWrite       Permission = "WRITE"
	PermissionWriteAcp    Permission = "WRITE_ACP"
	PermissionRead        Permission = "READ"
	PermissionReadAcp     Permission = "READ_ACP"
)

type BucketCannedACL string

const (
	BucketCannedACLPrivate           BucketCannedACL = "private"
	BucketCannedACLPublicRead        BucketCannedACL = "public-read"
	BucketCannedACLPublicReadWrite   BucketCannedACL = "public-read-write"
	BucketCannedACLAuthenticatedRead BucketCannedACL = "authenticated-read"
)

type Bucket struct {
	Name   string
	Policy *string
	Grants []BucketGrant
}

type BucketGrant struct {
	Grantee    *Grantee
	Permission Permission
}

type Grantee struct {
	Type         Type
	DisplayName  *string
	EmailAddress *string
	ID           *string
	URI          *string
}

// Note: Ceph create and update operations only return error status.
type BucketRepo interface {
	List(ctx context.Context, scope string) ([]Bucket, error)
	Get(ctx context.Context, scope, bucket string) (*Bucket, error)
	Create(ctx context.Context, scope, bucket string, acl BucketCannedACL) error
	UpdateOwner(ctx context.Context, scope, bucket, owner string) error
	UpdateACL(ctx context.Context, scope, bucket string, acl BucketCannedACL) error
	UpdatePolicy(ctx context.Context, scope, bucket, policy string) error
	Delete(ctx context.Context, scope, bucket string) error
	Endpoint(scope string) string
	Key(scope string) (accessKey string, secretKey string)
}

func (uc *ObjectUseCase) ListBuckets(ctx context.Context, scope string) ([]Bucket, error) {
	return uc.bucket.List(ctx, scope)
}

func (uc *ObjectUseCase) CreateBucket(ctx context.Context, scope, bucket, owner, policy string, acl BucketCannedACL) (*Bucket, error) {
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

	return uc.bucket.Get(ctx, scope, bucket)
}

func (uc *ObjectUseCase) UpdateBucket(ctx context.Context, scope, bucket, owner, policy string, acl BucketCannedACL) (*Bucket, error) {
	if err := uc.bucket.UpdateACL(ctx, scope, bucket, acl); err != nil {
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

	return uc.bucket.Get(ctx, scope, bucket)
}

func (uc *ObjectUseCase) DeleteBucket(ctx context.Context, scope, bucket string) error {
	return uc.bucket.Delete(ctx, scope, bucket)
}
