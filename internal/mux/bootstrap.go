package mux

import (
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"

	bootstrapv1 "github.com/otterscale/otterscale/api/bootstrap/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/app"
)

type Bootstrap struct {
	*http.ServeMux

	bootstrap *app.BootstrapService
}

func NewBootstrap(bootstrap *app.BootstrapService) *Bootstrap {
	return &Bootstrap{
		ServeMux:  &http.ServeMux{},
		bootstrap: bootstrap,
	}
}

func (b *Bootstrap) RegisterHandlers(opts []connect.HandlerOption) error {
	// Prepare service names for reflection and health check
	services := []string{bootstrapv1.BootstrapServiceName}

	// Register gRPC reflection
	reflector := grpcreflect.NewStaticReflector(services...)
	b.Handle(grpcreflect.NewHandlerV1(reflector))
	b.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// Register gRPC health check
	checker := grpchealth.NewStaticChecker(services...)
	b.Handle(grpchealth.NewHandler(checker))

	// Register service handlers
	b.Handle(bootstrapv1.NewBootstrapServiceHandler(b.bootstrap, opts...))

	return nil
}
