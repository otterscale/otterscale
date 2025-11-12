package maas

import (
	"context"
	"strconv"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/network"
)

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

	ipRanges, err := client.IPRanges.Get()
	if err != nil {
		return nil, err
	}

	return r.toIPRanges(ipRanges), nil
}

func (r *ipRangeRepo) Create(_ context.Context, subnetID int, startIP, endIP, comment string) (*network.IPRange, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.IPRangeParams{
		Subnet:  strconv.Itoa(subnetID),
		StartIP: startIP,
		EndIP:   endIP,
		Comment: comment,
	}

	ipRange, err := client.IPRanges.Create(params)
	if err != nil {
		return nil, err
	}

	return r.toIPRange(ipRange), nil
}

func (r *ipRangeRepo) Update(ctx context.Context, id int, startIP, endIP, comment string) (*network.IPRange, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.IPRangeParams{
		StartIP: startIP,
		EndIP:   endIP,
		Comment: comment,
	}

	ipRange, err := client.IPRange.Update(id, params)
	if err != nil {
		return nil, err
	}

	return r.toIPRange(ipRange), nil
}

func (r *ipRangeRepo) Delete(_ context.Context, id int) error {
	client, err := r.maas.Client()
	if err != nil {
		return err
	}
	return client.IPRange.Delete(id)
}

func (r *ipRangeRepo) toIPRange(pr *entity.IPRange) *network.IPRange {
	return &network.IPRange{}
}

func (r *ipRangeRepo) toIPRanges(prs []entity.IPRange) []network.IPRange {
	ret := make([]network.IPRange, 0, len(prs))

	return ret
}
