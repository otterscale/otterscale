package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/network"
)

type vlanRepo struct {
	maas *MAAS
}

func NewVLANRepo(maas *MAAS) network.VLANRepo {
	return &vlanRepo{
		maas: maas,
	}
}

var _ network.VLANRepo = (*vlanRepo)(nil)

func (r *vlanRepo) Update(_ context.Context, fabricID, vid int, name string, mtu int, description string, dhcpOn bool) (*network.VLAN, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.VLANParams{
		VID:         vid,
		Name:        name,
		MTU:         mtu,
		Description: description,
		DHCPOn:      dhcpOn,
	}

	return client.VLAN.Update(fabricID, vid, params)
}
