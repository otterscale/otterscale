package app

import (
	"context"

	v1 "github.com/openhdc/openhdc/api/user/v1"
	"github.com/openhdc/openhdc/internal/service/domain/service"
)

type UserApp struct {
	v1.UnimplementedUserServer

	svc *service.UserService
}

func NewUserApp(svc *service.UserService) *UserApp {
	return &UserApp{
		svc: svc,
	}
}

var _ v1.UserServer = (*UserApp)(nil)

func (a *UserApp) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	_, err := a.svc.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &v1.GetUserResponse{}, nil
}
