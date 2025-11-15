package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/configuration"
)

// Note: MAAS API do not support context.
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

	return client.BootSourceSelections.Get(id)
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

	return client.BootSourceSelections.Create(bootSourceID, params)
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

	return client.BootSourceSelection.Update(bootSourceID, id, params)
}
