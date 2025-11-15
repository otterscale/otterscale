package maas

import (
	"context"

	"github.com/otterscale/otterscale/internal/core/configuration"
)

// Note: MAAS API do not support context.
type bootResourceRepo struct {
	maas *MAAS
}

func NewBootResourceRepo(maas *MAAS) configuration.BootResourceRepo {
	return &bootResourceRepo{
		maas: maas,
	}
}

var _ configuration.BootResourceRepo = (*bootResourceRepo)(nil)

func (r *bootResourceRepo) List(_ context.Context) ([]configuration.BootResource, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	return client.BootResources.Get(nil)
}

func (r *bootResourceRepo) Import(_ context.Context) error {
	client, err := r.maas.Client()
	if err != nil {
		return err
	}

	return client.BootResources.Import()
}

func (r *bootResourceRepo) IsImporting(_ context.Context) (bool, error) {
	client, err := r.maas.Client()
	if err != nil {
		return false, err
	}

	return client.BootResources.IsImporting()
}
