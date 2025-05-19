package maas

import (
	"context"

	"github.com/openhdc/otterscale/internal/domain/service"

	"github.com/canonical/gomaasclient/entity"
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

func (r *bootResource) List(_ context.Context) ([]entity.BootResource, error) {
	return r.maas.BootResources.Get(&entity.BootResourcesReadParams{})
}

func (r *bootResource) Import(_ context.Context) error {
	return r.maas.BootResources.Import()
}

func (r *bootResource) IsImporting(_ context.Context) (bool, error) {
	return r.maas.BootResources.IsImporting()
}
