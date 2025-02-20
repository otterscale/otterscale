package service

import (
	"context"

	"github.com/openhdc/openhdc/internal/service/domain/model"
)

type UserRepo interface {
	Get(ctx context.Context) (*model.User, error)
}

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Get(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.Get(ctx)
}
