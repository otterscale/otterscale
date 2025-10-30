package mux

import (
	"net/http"

	"connectrpc.com/connect"

	bootstrapv1 "github.com/otterscale/otterscale/api/bootstrap/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/app"
)

type Bootstrap struct {
	*http.ServeMux

	svc *app.BootstrapService
}

func (b *Bootstrap) serviceNames() []string {
	return []string{bootstrapv1.BootstrapServiceName}
}

func (b *Bootstrap) RegisterHandlers(_ []connect.HandlerOption) {
	b.Handle(bootstrapv1.NewBootstrapServiceHandler(b.svc))
}

func NewBootstrap(svc *app.BootstrapService) *Bootstrap {
	// Initialize Bootstrap and register all handlers
	bootstrap := &Bootstrap{
		ServeMux: &http.ServeMux{},
		svc:      svc,
	}

	// Register gRPC reflection and health check
	registerGRPCReflection(bootstrap.ServeMux, bootstrap.serviceNames()...)
	registerGRPCHealth(bootstrap.ServeMux, bootstrap.serviceNames()...)

	return bootstrap
}
