package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/openhdc/openhdc/internal/domain/service"
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

func (r *fabric) List(ctx context.Context) ([]*entity.Fabric, error) {
	fs, err := r.maas.Fabrics.Get()
	if err != nil {
		return nil, err
	}

	ret := make([]*entity.Fabric, len(fs))
	for i := range fs {
		ret[i] = &fs[i]
	}
	return ret, nil
}

func (r *fabric) Get(_ context.Context, id int) (*entity.Fabric, error) {
	return r.maas.Fabric.Get(id)
}

func (r *fabric) Create(_ context.Context, params *entity.FabricParams) (*entity.Fabric, error) {
	return r.maas.Fabrics.Create(params)
}

func (r *fabric) Update(_ context.Context, id int, params *entity.FabricParams) (*entity.Fabric, error) {
	return r.maas.Fabric.Update(id, params)
}

func (r *fabric) Delete(_ context.Context, id int) error {
	return r.maas.Fabric.Delete(id)
}
