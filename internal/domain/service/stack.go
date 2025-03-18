package service

import (
	"context"

	"github.com/openhdc/openhdc/internal/domain/model"
)

type StackService struct {
	user UserRepo
}

type UserRepo interface {
	Get(ctx context.Context, id string) (*model.User, error)
}

func NewStackService(user UserRepo) *StackService {
	return &StackService{
		user: user,
	}
}
