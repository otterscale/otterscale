package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core"
)

type vlan struct {
	maas *MAAS
}

func NewVLAN(maas *MAAS) core.VLANRepo {
	return &vlan{
		maas: maas,
	}
}

var _ core.VLANRepo = (*vlan)(nil)

func (r *vlan) Update(_ context.Context, fabricID, vid int, params *entity.VLANParams) (*core.VLAN, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.VLAN.Update(fabricID, vid, params)
}
