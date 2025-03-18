package app

import (
	pb "github.com/openhdc/openhdc/api/stack/v1"
	"github.com/openhdc/openhdc/internal/domain/service"
)

type StackApp struct {
	pb.UnimplementedStackServiceServer

	svc *service.StackService
}

func NewStackApp(svc *service.StackService) *StackApp {
	return &StackApp{
		svc: svc,
	}
}
