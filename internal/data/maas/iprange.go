package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core"
)

type ipRange struct {
	maas *MAAS
}

func NewIPRange(maas *MAAS) core.IPRangeRepo {
	return &ipRange{
		maas: maas,
	}
}

var _ core.IPRangeRepo = (*ipRange)(nil)

func (r *ipRange) List(_ context.Context) ([]core.IPRange, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.IPRanges.Get()
}

func (r *ipRange) Create(_ context.Context, params *entity.IPRangeParams) (*core.IPRange, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.IPRanges.Create(params)
}

func (r *ipRange) Update(_ context.Context, id int, params *entity.IPRangeParams) (*core.IPRange, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.IPRange.Update(id, params)
}

func (r *ipRange) Delete(_ context.Context, id int) error {
	client, err := r.maas.client()
	if err != nil {
		return err
	}
	return client.IPRange.Delete(id)
}
