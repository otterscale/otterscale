package maas

import (
	"context"

	"github.com/openhdc/openhdc/internal/domain/service"

	"github.com/canonical/gomaasclient/entity"
)

const maasIOBootSourceID = 1

type bootSourceSelection struct {
	maas *MAAS
}

func NewBootSourceSelection(maas *MAAS) service.MAASBootSourceSelection {
	return &bootSourceSelection{
		maas: maas,
	}
}

var _ service.MAASBootSourceSelection = (*bootSourceSelection)(nil)

func (r *bootSourceSelection) List(_ context.Context, id int) ([]entity.BootSourceSelection, error) {
	return r.maas.BootSourceSelections.Get(id)
}

func (r *bootSourceSelection) CreateFromMAASIO(_ context.Context, distroSeries string, architectures []string) (*entity.BootSourceSelection, error) {
	return r.maas.BootSourceSelections.Create(maasIOBootSourceID, &entity.BootSourceSelectionParams{
		OS:        "ubuntu",
		Release:   distroSeries,
		Arches:    architectures,
		Subarches: []string{"*"},
		Labels:    []string{"*"},
	})
}
