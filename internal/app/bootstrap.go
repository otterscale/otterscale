package app

import (
	"github.com/otterscale/otterscale/api/bootstrap/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type BootstrapService struct {
	pbconnect.UnimplementedBootstrapServiceHandler

	uc *core.BootstrapUseCase
}

func NewBootstrapService(uc *core.BootstrapUseCase) *BootstrapService {
	return &BootstrapService{uc: uc}
}

var _ pbconnect.BootstrapHandler = (*Bootstrap)(nil)
