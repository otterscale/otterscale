package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/machine/tag"
)

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

	tags, err := client.Tags.Get()
	if err != nil {
		return nil, err
	}

	return r.toTags(tags), nil
}

func (r *tagRepo) Get(_ context.Context, name string) (*tag.Tag, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	tag, err := client.Tag.Get(name)
	if err != nil {
		return nil, err
	}

	return r.toTag(tag), nil
}

func (r *tagRepo) Create(_ context.Context, name, comment string) (*tag.Tag, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	tag, err := client.Tags.Create(&entity.TagParams{
		Name:    name,
		Comment: comment,
	})
	if err != nil {
		return nil, err
	}

	return r.toTag(tag), nil
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

func (r *tagRepo) toTag(t *entity.Tag) *tag.Tag {
	return &tag.Tag{}
}

func (r *tagRepo) toTags(ts []entity.Tag) []tag.Tag {
	ret := make([]tag.Tag, 0, len(ts))

	return ret
}
