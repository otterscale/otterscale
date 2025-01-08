package service

import (
	"context"

	v1 "github.com/openhdc/openhdc/api/user/v1"
	"github.com/openhdc/openhdc/internal/biz"
)

type UserService struct {
	v1.UnimplementedUserServer

	uc *biz.UserUseCase
}

func NewSampleService(uc *biz.UserUseCase) *UserService {
	return &UserService{
		uc: uc,
	}
}

var _ v1.UserServer = (*UserService)(nil)

func (s *UserService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	_, err := s.uc.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &v1.GetUserResponse{}, nil
}
