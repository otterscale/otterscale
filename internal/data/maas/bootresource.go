package maas

import (
	"context"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type bootResource struct {
	maas *MAAS
}

func NewBootResource(maas *MAAS) service.MAASBootResource {
	return &bootResource{
		maas: maas,
	}
}

var _ service.MAASBootResource = (*bootResource)(nil)

func (r *bootResource) Import(_ context.Context) error {
	return r.maas.BootResources.Import()
}
