package app

import (
	"github.com/openhdc/openhdc/api/nexus/v1/pbconnect"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type NexusApp struct {
	pbconnect.UnimplementedNexusHandler
	svc *service.NexusService
}

func NewNexusApp(svc *service.NexusService) *NexusApp {
	return &NexusApp{svc: svc}
}

var _ pbconnect.NexusHandler = (*NexusApp)(nil)
