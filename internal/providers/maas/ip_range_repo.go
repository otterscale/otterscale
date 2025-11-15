package maas

import (
	"context"
	"strconv"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/network"
)

// Note: MAAS API do not support context.
type ipRangeRepo struct {
	maas *MAAS
}

func NewIPRangeRepo(maas *MAAS) network.IPRangeRepo {
	return &ipRangeRepo{
		maas: maas,
	}
}

var _ network.IPRangeRepo = (*ipRangeRepo)(nil)

func (r *ipRangeRepo) List(_ context.Context) ([]network.IPRange, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	return client.IPRanges.Get()
}

func (r *ipRangeRepo) Create(_ context.Context, subnetID int, startIP, endIP, comment string) (*network.IPRange, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.IPRangeParams{
		Type:    "reserved",
		Subnet:  strconv.Itoa(subnetID),
		StartIP: startIP,
		EndIP:   endIP,
		Comment: comment,
	}

	return client.IPRanges.Create(params)
}

func (r *ipRangeRepo) Update(_ context.Context, id int, startIP, endIP, comment string) (*network.IPRange, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.IPRangeParams{
		StartIP: startIP,
		EndIP:   endIP,
		Comment: comment,
	}

	return client.IPRange.Update(id, params)
}

func (r *ipRangeRepo) Delete(_ context.Context, id int) error {
	client, err := r.maas.Client()
	if err != nil {
		return err
	}

	return client.IPRange.Delete(id)
}
