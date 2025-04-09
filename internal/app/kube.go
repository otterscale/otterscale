package app

import (
	"github.com/openhdc/openhdc/api/kube/v1/v1connect"
	"github.com/openhdc/openhdc/internal/domain/service"
)

// KubeApp implements the KubeServiceHandler interface
type KubeApp struct {
	v1connect.UnimplementedKubeServiceHandler
	svc *service.KubeService
}

// NewKubeApp creates a new KubeApp instance
func NewKubeApp(svc *service.KubeService) *KubeApp {
	return &KubeApp{svc: svc}
}

// Ensure KubeApp implements the KubeServiceHandler interface
var _ v1connect.KubeServiceHandler = (*KubeApp)(nil)
