package maas

import (
	"context"

	"github.com/openhdc/openhdc/internal/domain/model"
	"github.com/openhdc/openhdc/internal/domain/service"
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

func (r *vlan) Update(_ context.Context, fabricID, vid int, params *model.VLANParams) (*model.VLAN, error) {
	return r.maas.VLAN.Update(fabricID, vid, params)
}
