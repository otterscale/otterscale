package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/openhdc/otterscale/internal/core"
)

type bootSourceSelection struct {
	maas *MAAS
}

func NewBootSourceSelection(maas *MAAS) core.BootSourceSelectionRepo {
	return &bootSourceSelection{
		maas: maas,
	}
}

var _ core.BootSourceSelectionRepo = (*bootSourceSelection)(nil)

func (r *bootSourceSelection) List(_ context.Context, id int) ([]entity.BootSourceSelection, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.BootSourceSelections.Get(id)
}

func (r *bootSourceSelection) Create(_ context.Context, bootSourceID int, params *entity.BootSourceSelectionParams) (*entity.BootSourceSelection, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.BootSourceSelections.Create(bootSourceID, params)
}
