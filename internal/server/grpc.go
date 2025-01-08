package server

import (
	"github.com/openhdc/openhdc"
	v1 "github.com/openhdc/openhdc/api/user/v1"
	"github.com/openhdc/openhdc/internal/service"
)

func NewGRPCServer(opts []openhdc.ServerOption, us *service.UserService) *openhdc.Server {
	srv := openhdc.NewDefaultServer(opts...)
	gs := srv.GRPCServer()
	v1.RegisterUserServer(gs, us)
	return srv
}
