package ceph

import (
	"context"

	"github.com/ceph/go-ceph/rgw/admin"

	"github.com/otterscale/otterscale/internal/core/storage"
)

type userRepo struct {
	ceph *Ceph
}

func NewUserRepo(ceph *Ceph) storage.UserRepo {
	return &userRepo{
		ceph: ceph,
	}
}

var _ storage.UserRepo = (*userRepo)(nil)

func (r *userRepo) List(ctx context.Context, scope string) ([]storage.User, error) {
	client, err := r.ceph.Client(scope)
	if err != nil {
		return nil, err
	}

	ids, err := client.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	if ids == nil {
		return []storage.User{}, nil
	}

	return r.list(ctx, client, *ids)
}

func (r *userRepo) list(ctx context.Context, client *admin.API, ids []string) ([]storage.User, error) {
	keys := []admin.UserKeySpec{{
		AccessKey: client.AccessKey,
		SecretKey: client.SecretKey,
	}}

	users := []storage.User{}

	for _, id := range ids {
		user, err := client.GetUser(ctx, admin.User{
			ID:   id,
			Keys: keys,
		})
		if err != nil {
			return nil, err
		}
		users = append(users, *r.toUser(&user))
	}

	return users, nil
}

func (r *userRepo) toUser(u *admin.User) *storage.User {
	return &storage.User{
		ID: u.ID,
	}
}
