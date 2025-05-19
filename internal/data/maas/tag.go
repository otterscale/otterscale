package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type tag struct {
	maas *MAAS
}

func NewTag(maas *MAAS) service.MAASTag {
	return &tag{
		maas: maas,
	}
}

var _ service.MAASTag = (*tag)(nil)

func (r *tag) List(_ context.Context) ([]entity.Tag, error) {
	return r.maas.Tags.Get()
}

func (r *tag) Get(_ context.Context, name string) (*entity.Tag, error) {
	return r.maas.Tag.Get(name)
}

func (r *tag) Create(_ context.Context, name, comment string) (*entity.Tag, error) {
	return r.maas.Tags.Create(&entity.TagParams{
		Name:    name,
		Comment: comment,
	})
}

func (r *tag) Delete(_ context.Context, name string) error {
	return r.maas.Tag.Delete(name)
}

func (r *tag) AddMachines(_ context.Context, name string, machineIDs []string) error {
	return r.maas.Tag.AddMachines(name, machineIDs)
}

func (r *tag) RemoveMachines(_ context.Context, name string, machineIDs []string) error {
	return r.maas.Tag.RemoveMachines(name, machineIDs)
}
