package maas

import (
	"context"

	entity "github.com/otterscale/otterscale/internal/core/_entity"
	"github.com/otterscale/otterscale/internal/core/configuration"
)

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

	resources, err := client.BootResources.Get(nil)
	if err != nil {
		return nil, err
	}

	return r.toBootResources(resources), nil
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

func (r *bootResourceRepo) toBootResources(brs []entity.BootResource) []configuration.BootResource {
	ret := []configuration.BootResource{}

	return ret
}
