package maas

import (
	"context"

	"github.com/openhdc/otterscale/internal/domain/service"

	"github.com/canonical/gomaasclient/entity"
)

type bootSource struct {
	maas *MAAS
}

func NewBootSource(maas *MAAS) service.MAASBootSource {
	return &bootSource{
		maas: maas,
	}
}

var _ service.MAASBootSource = (*bootSource)(nil)

func (r *bootSource) List(_ context.Context) ([]entity.BootSource, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.BootSources.Get()
}
