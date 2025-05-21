package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type fabric struct {
	maas *MAAS
}

func NewFabric(maas *MAAS) service.MAASFabric {
	return &fabric{
		maas: maas,
	}
}

var _ service.MAASFabric = (*fabric)(nil)

func (r *fabric) List(_ context.Context) ([]entity.Fabric, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Fabrics.Get()
}

func (r *fabric) Get(_ context.Context, id int) (*entity.Fabric, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Fabric.Get(id)
}

func (r *fabric) Create(_ context.Context, params *entity.FabricParams) (*entity.Fabric, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Fabrics.Create(params)
}

func (r *fabric) Update(_ context.Context, id int, params *entity.FabricParams) (*entity.Fabric, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Fabric.Update(id, params)
}

func (r *fabric) Delete(_ context.Context, id int) error {
	client, err := r.maas.client()
	if err != nil {
		return err
	}
	return client.Fabric.Delete(id)
}
