package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/openhdc/otterscale/internal/core"
)

type subnet struct {
	maas *MAAS
}

func NewSubnet(maas *MAAS) core.SubnetRepo {
	return &subnet{
		maas: maas,
	}
}

var _ core.SubnetRepo = (*subnet)(nil)

func (r *subnet) List(_ context.Context) ([]core.Subnet, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Subnets.Get()
}

func (r *subnet) Get(_ context.Context, id int) (*core.Subnet, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Subnet.Get(id)
}

func (r *subnet) Create(_ context.Context, params *entity.SubnetParams) (*core.Subnet, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Subnets.Create(params)
}

func (r *subnet) Update(_ context.Context, id int, params *entity.SubnetParams) (*core.Subnet, error) {
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

func (r *subnet) GetIPAddresses(_ context.Context, id int) ([]core.IPAddress, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Subnet.GetIPAddresses(id)
}

func (r *subnet) GetStatistics(_ context.Context, id int) (*core.NetworkStatistics, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Subnet.GetStatistics(id)
}
