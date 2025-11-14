package ceph

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/otterscale/otterscale/internal/core/storage/object"
)

type userRepo struct {
	ceph *Ceph
}

func NewUserRepo(ceph *Ceph) object.UserRepo {
	return &userRepo{
		ceph: ceph,
	}
}

var _ object.UserRepo = (*userRepo)(nil)

func (r *userRepo) List(ctx context.Context, scope string) ([]object.User, error) {
	client, err := r.ceph.client(scope)
	if err != nil {
		return nil, err
	}

	ids, err := client.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	if ids == nil {
		return []object.User{}, nil
	}

	return r.list(ctx, client, *ids)
}

func (r *userRepo) Create(ctx context.Context, scope, id, name string, suspended bool) (*object.User, error) {
	client, err := r.ceph.client(scope)
	if err != nil {
		return nil, err
	}

	user, err := client.CreateUser(ctx, admin.User{
		ID:          id,
		DisplayName: name,
		Suspended:   r.intPtr(suspended),
	})
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (r *userRepo) Update(ctx context.Context, scope, id, name string, suspended bool) (*object.User, error) {
	client, err := r.ceph.client(scope)
	if err != nil {
		return nil, err
	}

	user, err := client.ModifyUser(ctx, admin.User{
		ID:          id,
		DisplayName: name,
		Suspended:   r.intPtr(suspended),
	})
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (r *userRepo) Delete(ctx context.Context, scope, id string) error {
	client, err := r.ceph.client(scope)
	if err != nil {
		return err
	}

	user := admin.User{
		ID: id,
	}

	return client.RemoveUser(ctx, user)
}

func (r *userRepo) CreateKey(ctx context.Context, scope, id string) (*object.UserKey, error) {
	client, err := r.ceph.client(scope)
	if err != nil {
		return nil, err
	}

	keys, err := client.CreateKey(ctx, admin.UserKeySpec{UID: id})
	if err != nil {
		return nil, err
	}

	if keys == nil {
		return nil, fmt.Errorf("create key returned nil")
	}

	if len(*keys) == 0 {
		return nil, fmt.Errorf("create key returned empty list")
	}

	return &(*keys)[0], nil
}

func (r *userRepo) DeleteKey(ctx context.Context, scope, id, accessKey string) error {
	client, err := r.ceph.client(scope)
	if err != nil {
		return err
	}

	key := admin.UserKeySpec{
		UID:       id,
		AccessKey: accessKey,
	}

	return client.RemoveKey(ctx, key)
}

func (r *userRepo) list(ctx context.Context, client *admin.API, ids []string) ([]object.User, error) {
	keys := []admin.UserKeySpec{{
		AccessKey: client.AccessKey,
		SecretKey: client.SecretKey,
	}}

	users := []object.User{}

	for _, id := range ids {
		user, err := client.GetUser(ctx, admin.User{
			ID:   id,
			Keys: keys,
		})
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepo) intPtr(b bool) *int {
	if b {
		return aws.Int(1)
	}
	return nil
}
