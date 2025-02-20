package repo

import (
	"context"

	"github.com/openhdc/openhdc/internal/service/domain/model"
	"github.com/openhdc/openhdc/internal/service/domain/service"
	"github.com/openhdc/openhdc/internal/service/infra/ent"
)

type userRepo struct {
	client *ent.Client
}

func NewUserRepo(client *ent.Client) service.UserRepo {
	return &userRepo{
		client: client,
	}
}

var _ service.UserRepo = (*userRepo)(nil)

func (r *userRepo) Get(ctx context.Context) (*model.User, error) {
	return &model.User{}, nil
}
