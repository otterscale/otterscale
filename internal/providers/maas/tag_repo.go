package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/machine/tag"
)

// Note: MAAS API do not support context.
type tagRepo struct {
	maas *MAAS
}

func NewTagRepo(maas *MAAS) tag.TagRepo {
	return &tagRepo{
		maas: maas,
	}
}

var _ tag.TagRepo = (*tagRepo)(nil)

func (r *tagRepo) List(_ context.Context) ([]tag.Tag, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	return client.Tags.Get()
}

func (r *tagRepo) Get(_ context.Context, name string) (*tag.Tag, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	return client.Tag.Get(name)
}

func (r *tagRepo) Create(_ context.Context, name, comment string) (*tag.Tag, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	return client.Tags.Create(&entity.TagParams{
		Name:    name,
		Comment: comment,
	})
}

func (r *tagRepo) Delete(_ context.Context, name string) error {
	client, err := r.maas.Client()
	if err != nil {
		return err
	}

	return client.Tag.Delete(name)
}

func (r *tagRepo) AddMachines(_ context.Context, name string, machineIDs []string) error {
	client, err := r.maas.Client()
	if err != nil {
		return err
	}

	return client.Tag.AddMachines(name, machineIDs)
}

func (r *tagRepo) RemoveMachines(_ context.Context, name string, machineIDs []string) error {
	client, err := r.maas.Client()
	if err != nil {
		return err
	}

	return client.Tag.RemoveMachines(name, machineIDs)
}
