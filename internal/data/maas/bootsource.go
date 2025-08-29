package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core"
)

type bootSource struct {
	maas *MAAS
}

func NewBootSource(maas *MAAS) core.BootSourceRepo {
	return &bootSource{
		maas: maas,
	}
}

var _ core.BootSourceRepo = (*bootSource)(nil)

func (r *bootSource) List(_ context.Context) ([]entity.BootSource, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.BootSources.Get()
}
