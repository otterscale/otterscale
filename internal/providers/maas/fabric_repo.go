package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/network"
)

// Note: MAAS API do not support context.
type fabricRepo struct {
	maas *MAAS
}

func NewFabricRepo(maas *MAAS) network.FabricRepo {
	return &fabricRepo{
		maas: maas,
	}
}

var _ network.FabricRepo = (*fabricRepo)(nil)

func (r *fabricRepo) List(_ context.Context) ([]network.Fabric, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	return client.Fabrics.Get()
}

func (r *fabricRepo) Get(_ context.Context, id int) (*network.Fabric, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	return client.Fabric.Get(id)
}

func (r *fabricRepo) Create(_ context.Context) (*network.Fabric, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.FabricParams{}

	return client.Fabrics.Create(params)
}

func (r *fabricRepo) Update(_ context.Context, id int, name string) (*network.Fabric, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.FabricParams{
		Name: name,
	}

	return client.Fabric.Update(id, params)
}

func (r *fabricRepo) Delete(_ context.Context, id int) error {
	client, err := r.maas.Client()
	if err != nil {
		return err
	}

	return client.Fabric.Delete(id)
}
