package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type ipRange struct {
	maas *MAAS
}

func NewIPRange(maas *MAAS) service.MAASIPRange {
	return &ipRange{
		maas: maas,
	}
}

var _ service.MAASIPRange = (*ipRange)(nil)

func (r *ipRange) List(_ context.Context) ([]entity.IPRange, error) {
	return r.maas.IPRanges.Get()
}

func (r *ipRange) Create(_ context.Context, params *entity.IPRangeParams) (*entity.IPRange, error) {
	return r.maas.IPRanges.Create(params)
}

func (r *ipRange) Update(_ context.Context, id int, params *entity.IPRangeParams) (*entity.IPRange, error) {
	return r.maas.IPRange.Update(id, params)
}

func (r *ipRange) Delete(_ context.Context, id int) error {
	return r.maas.IPRange.Delete(id)
}
