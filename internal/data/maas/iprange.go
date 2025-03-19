package maas

import (
	"context"

	"github.com/openhdc/openhdc/internal/domain/model"
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

func (r *ipRange) List(ctx context.Context) ([]*model.IPRange, error) {
	rs, err := r.maas.IPRanges.Get()
	if err != nil {
		return nil, err
	}

	ret := make([]*model.IPRange, len(rs))
	for i := range rs {
		ret[i] = &rs[i]
	}
	return ret, nil
}

func (r *ipRange) Create(_ context.Context, params *model.IPRangeParams) (*model.IPRange, error) {
	return r.maas.IPRanges.Create(params)
}

func (r *ipRange) Update(_ context.Context, id int, params *model.IPRangeParams) (*model.IPRange, error) {
	return r.maas.IPRange.Update(id, params)
}
