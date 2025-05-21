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
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.BootResources.Get(nil)
}

func (r *bootResource) Import(_ context.Context) error {
	client, err := r.maas.client()
	if err != nil {
		return err
	}
	return client.BootResources.Import()
}

func (r *bootResource) IsImporting(_ context.Context) (bool, error) {
	client, err := r.maas.client()
	if err != nil {
		return false, err
	}
	return client.BootResources.IsImporting()
}
