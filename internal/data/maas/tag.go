package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core"
)

type tag struct {
	maas *MAAS
}

func NewTag(maas *MAAS) core.TagRepo {
	return &tag{
		maas: maas,
	}
}

var _ core.TagRepo = (*tag)(nil)

func (r *tag) List(_ context.Context) ([]core.Tag, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Tags.Get()
}

func (r *tag) Get(_ context.Context, name string) (*core.Tag, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Tag.Get(name)
}

func (r *tag) Create(_ context.Context, name, comment string) (*core.Tag, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Tags.Create(&entity.TagParams{
		Name:    name,
		Comment: comment,
	})
}

func (r *tag) Delete(_ context.Context, name string) error {
	client, err := r.maas.client()
	if err != nil {
		return err
	}
	return client.Tag.Delete(name)
}

func (r *tag) AddMachines(_ context.Context, name string, machineIDs []string) error {
	client, err := r.maas.client()
	if err != nil {
		return err
	}
	return client.Tag.AddMachines(name, machineIDs)
}

func (r *tag) RemoveMachines(_ context.Context, name string, machineIDs []string) error {
	client, err := r.maas.client()
	if err != nil {
		return err
	}
	return client.Tag.RemoveMachines(name, machineIDs)
}
