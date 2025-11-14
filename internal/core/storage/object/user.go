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

func (uc *ObjectUseCase) ListUsers(ctx context.Context, scope string) ([]User, error) {
	return uc.user.List(ctx, scope)
}

func (uc *ObjectUseCase) CreateUser(ctx context.Context, scope, id, name string, suspended bool) (*User, error) {
	return uc.user.Create(ctx, scope, id, name, suspended)
}

func (uc *ObjectUseCase) UpdateUser(ctx context.Context, scope, id, name string, suspended bool) (*User, error) {
	return uc.user.Update(ctx, scope, id, name, suspended)
}

func (uc *ObjectUseCase) DeleteUser(ctx context.Context, scope, id string) error {
	return uc.user.Delete(ctx, scope, id)
}

func (uc *ObjectUseCase) CreateUserKey(ctx context.Context, scope, id string) (*UserKey, error) {
	return uc.user.CreateKey(ctx, scope, id)
}

func (uc *ObjectUseCase) DeleteUserKey(ctx context.Context, scope, id, accessKey string) error {
	return uc.user.DeleteKey(ctx, scope, id, accessKey)
}
