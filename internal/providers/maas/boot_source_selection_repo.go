package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/configuration"
)

type bootSourceSelectionRepo struct {
	maas *MAAS
}

func NewBootSourceSelectionRepo(maas *MAAS) configuration.BootSourceSelectionRepo {
	return &bootSourceSelectionRepo{
		maas: maas,
	}
}

var _ configuration.BootSourceSelectionRepo = (*bootSourceSelectionRepo)(nil)

func (r *bootSourceSelectionRepo) List(_ context.Context, id int) ([]configuration.BootSourceSelection, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	selections, err := client.BootSourceSelections.Get(id)
	if err != nil {
		return nil, err
	}

	return r.toBootSourceSelections(selections), nil
}

func (r *bootSourceSelectionRepo) Create(_ context.Context, bootSourceID int, distroSeries string, architectures []string) (*configuration.BootSourceSelection, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.BootSourceSelectionParams{
		OS:        "ubuntu",
		Release:   distroSeries,
		Arches:    architectures,
		Subarches: []string{"*"},
		Labels:    []string{"*"},
	}

	selection, err := client.BootSourceSelections.Create(bootSourceID, params)
	if err != nil {
		return nil, err
	}

	return r.toBootSourceSelection(selection), nil
}

func (r *bootSourceSelectionRepo) Update(_ context.Context, bootSourceID, id int, distroSeries string, architectures []string) (*configuration.BootSourceSelection, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.BootSourceSelectionParams{
		OS:        "ubuntu",
		Release:   distroSeries,
		Arches:    architectures,
		Subarches: []string{"*"},
		Labels:    []string{"*"},
	}

	selection, err := client.BootSourceSelection.Update(bootSourceID, id, params)
	if err != nil {
		return nil, err
	}

	return r.toBootSourceSelection(selection), nil

}

func (r *bootSourceSelectionRepo) toBootSourceSelection(bss *entity.BootSourceSelection) *configuration.BootSourceSelection {
	return &configuration.BootSourceSelection{}
}

func (r *bootSourceSelectionRepo) toBootSourceSelections(bsss []entity.BootSourceSelection) []configuration.BootSourceSelection {
	ret := []configuration.BootSourceSelection{}

	return ret
}
