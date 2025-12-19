package object

import (
	"context"

	"github.com/ceph/go-ceph/rgw/admin"
)

type (
	// User represents a Ceph User resource.
	User = admin.User

	// UserKey represents a Ceph UserKey resource.
	UserKey = admin.UserKeySpec
)

type UserRepo interface {
	List(ctx context.Context, scope string) ([]User, error)
	Create(ctx context.Context, scope, id, name string, suspended bool) (*User, error)
	Update(ctx context.Context, scope, id, name string, suspended bool) (*User, error)
	Delete(ctx context.Context, scope, id string) error
	CreateKey(ctx context.Context, scope, id string) (*UserKey, error)
	DeleteKey(ctx context.Context, scope, id, accessKey string) error
}

func (uc *UseCase) ListUsers(ctx context.Context, scope string) (users []User, uri string, err error) {
	users, err = uc.user.List(ctx, scope)
	if err != nil {
		return nil, "", err
	}
	return users, uc.bucket.Endpoint(scope), nil
}

func (uc *UseCase) CreateUser(ctx context.Context, scope, id, name string, suspended bool) (*User, error) {
	u, err := uc.user.Create(ctx, scope, id, name, suspended)
	if err != nil {
		return nil, err
	}

	// remove default keys
	for _, k := range u.Keys {
		if err := uc.DeleteUserKey(ctx, scope, id, k.AccessKey); err != nil {
			return nil, err
		}
	}

	return u, nil
}

func (uc *UseCase) UpdateUser(ctx context.Context, scope, id, name string, suspended bool) (*User, error) {
	return uc.user.Update(ctx, scope, id, name, suspended)
}

func (uc *UseCase) DeleteUser(ctx context.Context, scope, id string) error {
	return uc.user.Delete(ctx, scope, id)
}

func (uc *UseCase) CreateUserKey(ctx context.Context, scope, id string) (*UserKey, error) {
	return uc.user.CreateKey(ctx, scope, id)
}

func (uc *UseCase) DeleteUserKey(ctx context.Context, scope, id, accessKey string) error {
	return uc.user.DeleteKey(ctx, scope, id, accessKey)
}
