package app

import (
	pb "github.com/openhdc/openhdc/api/stack/v1"
	"github.com/openhdc/openhdc/internal/domain/service"
)

type StackApp struct {
	pb.UnimplementedStackServiceServer

	ks *service.KubeService
}

func NewStackApp(ks *service.KubeService) *StackApp {
	return &StackApp{
		ks: ks,
	}
}
