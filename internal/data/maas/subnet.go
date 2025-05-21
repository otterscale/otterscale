package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"
	api "github.com/canonical/gomaasclient/entity/subnet"

	"github.com/openhdc/otterscale/internal/domain/service"
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

func (r *subnet) List(_ context.Context) ([]entity.Subnet, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Subnets.Get()
}

func (r *subnet) Get(_ context.Context, id int) (*entity.Subnet, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Subnet.Get(id)
}

func (r *subnet) Create(_ context.Context, params *entity.SubnetParams) (*entity.Subnet, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Subnets.Create(params)
}

func (r *subnet) Update(_ context.Context, id int, params *entity.SubnetParams) (*entity.Subnet, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Subnet.Update(id, params)
}

func (r *subnet) Delete(_ context.Context, id int) error {
	client, err := r.maas.client()
	if err != nil {
		return err
	}
	return client.Subnet.Delete(id)
}

func (r *subnet) GetIPAddresses(_ context.Context, id int) ([]api.IPAddress, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Subnet.GetIPAddresses(id)
}

func (r *subnet) GetReservedIPRanges(_ context.Context, id int) ([]api.ReservedIPRange, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Subnet.GetReservedIPRanges(id)
}

func (r *subnet) GetUnreservedIPRanges(_ context.Context, id int) ([]api.IPRange, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Subnet.GetUnreservedIPRanges(id)
}

func (r *subnet) GetStatistics(_ context.Context, id int) (*api.Statistics, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Subnet.GetStatistics(id)
}
