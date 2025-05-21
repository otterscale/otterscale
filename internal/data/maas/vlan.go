package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type vlan struct {
	maas *MAAS
}

func NewVLAN(maas *MAAS) service.MAASVLAN {
	return &vlan{
		maas: maas,
	}
}

var _ service.MAASVLAN = (*vlan)(nil)

func (r *vlan) Update(_ context.Context, fabricID, vid int, params *entity.VLANParams) (*entity.VLAN, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.VLAN.Update(fabricID, vid, params)
}
