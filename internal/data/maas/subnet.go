package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type subnet struct {
	maas *MAAS
}

func NewSubnet(maas *MAAS) service.MAASSubnet {
	return &subnet{
		maas: maas,
	}
}

var _ service.MAASSubnet = (*subnet)(nil)

func (r *subnet) List(ctx context.Context) ([]*entity.Subnet, error) {
	fs, err := r.maas.Subnets.Get()
	if err != nil {
		return nil, err
	}

	ret := make([]*entity.Subnet, len(fs))
	for i := range fs {
		ret[i] = &fs[i]
	}
	return ret, nil
}

func (r *subnet) Get(_ context.Context, id int) (*entity.Subnet, error) {
	return r.maas.Subnet.Get(id)
}

func (r *subnet) Create(_ context.Context, params *entity.SubnetParams) (*entity.Subnet, error) {
	return r.maas.Subnets.Create(params)
}

func (r *subnet) Update(_ context.Context, id int, params *entity.SubnetParams) (*entity.Subnet, error) {
	return r.maas.Subnet.Update(id, params)
}

func (r *subnet) Delete(_ context.Context, id int) error {
	return r.maas.Subnet.Delete(id)
}
