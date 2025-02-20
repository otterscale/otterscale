package server

import (
	"github.com/openhdc/openhdc"
	v1 "github.com/openhdc/openhdc/api/user/v1"
	"github.com/openhdc/openhdc/internal/service/app"
)

func NewGRPCServer(opts []openhdc.ServerOption, ua *app.UserApp) *openhdc.Server {
	srv := openhdc.NewDefaultServer(opts...)
	gs := srv.GRPCServer()
	v1.RegisterUserServer(gs, ua)
	return srv
}
