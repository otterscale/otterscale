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
	return r.maas.Subnets.Get()
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

func (r *subnet) GetIPAddresses(_ context.Context, id int) ([]api.IPAddress, error) {
	return r.maas.Subnet.GetIPAddresses(id)
}

func (r *subnet) GetReservedIPRanges(_ context.Context, id int) ([]api.ReservedIPRange, error) {
	return r.maas.Subnet.GetReservedIPRanges(id)
}

func (r *subnet) GetUnreservedIPRanges(_ context.Context, id int) ([]api.IPRange, error) {
	return r.maas.Subnet.GetUnreservedIPRanges(id)
}

func (r *subnet) GetStatistics(_ context.Context, id int) (*api.Statistics, error) {
	return r.maas.Subnet.GetStatistics(id)
}
