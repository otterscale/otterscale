package app

import (
	"github.com/openhdc/openhdc/api/stack/v1/v1connect"
	"github.com/openhdc/openhdc/internal/domain/service"
)

// StackApp implements the StackServiceServer interface
type StackApp struct {
	v1connect.UnimplementedStackServiceHandler
	svc *service.StackService
}

// NewStackApp creates a new StackApp instance
func NewStackApp(svc *service.StackService) *StackApp {
	return &StackApp{svc: svc}
}

// Ensure StackApp implements the StackServiceServer interface
var _ v1connect.StackServiceHandler = (*StackApp)(nil)
