package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/network"
)

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

	frabrics, err := client.Fabrics.Get()
	if err != nil {
		return nil, err
	}

	return r.toFabrics(frabrics), nil
}

func (r *fabricRepo) Get(_ context.Context, id int) (*network.Fabric, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	fabric, err := client.Fabric.Get(id)
	if err != nil {
		return nil, err
	}

	return r.toFabric(fabric), nil
}

func (r *fabricRepo) Create(_ context.Context) (*network.Fabric, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.FabricParams{}

	fabric, err := client.Fabrics.Create(params)
	if err != nil {
		return nil, err
	}

	return r.toFabric(fabric), nil
}

func (r *fabricRepo) Update(_ context.Context, id int, name string) (*network.Fabric, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.FabricParams{
		Name: name,
	}

	fabric, err := client.Fabric.Update(id, params)
	if err != nil {
		return nil, err
	}

	return r.toFabric(fabric), nil
}

func (r *fabricRepo) Delete(_ context.Context, id int) error {
	client, err := r.maas.Client()
	if err != nil {
		return err
	}
	return client.Fabric.Delete(id)
}

func (r *fabricRepo) toFabric(pr *entity.Fabric) *network.Fabric {
	return &network.Fabric{}
}

func (r *fabricRepo) toFabrics(prs []entity.Fabric) []network.Fabric {
	ret := make([]network.Fabric, 0, len(prs))

	return ret
}
