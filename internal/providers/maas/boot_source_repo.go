package maas

import (
	"context"

	"github.com/otterscale/otterscale/internal/core/configuration"
)

type bootSourceRepo struct {
	maas *MAAS
}

func NewBootSourceRepo(maas *MAAS) configuration.BootSourceRepo {
	return &bootSourceRepo{
		maas: maas,
	}
}

var _ configuration.BootSourceRepo = (*bootSourceRepo)(nil)

func (r *bootSourceRepo) List(_ context.Context) ([]configuration.BootSource, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	return client.BootSources.Get()
}
