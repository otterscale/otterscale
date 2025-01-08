package data

import (
	"context"

	"github.com/openhdc/openhdc/internal/biz"
)

type userRepo struct {
	data *Data
}

func NewUserRepo(d *Data) biz.UserRepo {
	return &userRepo{
		data: d,
	}
}

var _ biz.UserRepo = (*userRepo)(nil)

func (r *userRepo) Get(ctx context.Context) (*biz.User, error) {
	u := &biz.User{}
	return u, nil
}
