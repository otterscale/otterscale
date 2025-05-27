package core

import (
	"context"

	"github.com/canonical/gomaasclient/entity"
)

type Tag = entity.Tag

type TagRepo interface {
	List(ctx context.Context) ([]Tag, error)
	Get(ctx context.Context, name string) (*Tag, error)
	Create(ctx context.Context, name, comment string) (*Tag, error)
	Delete(ctx context.Context, name string) error
	AddMachines(ctx context.Context, name string, machineIDs []string) error
	RemoveMachines(ctx context.Context, name string, machineIDs []string) error
}

type TagUseCase struct {
	tag TagRepo
}

func NewTagUseCase(tag TagRepo) *TagUseCase {
	return &TagUseCase{
		tag: tag,
	}
}

func (uc *TagUseCase) ListTags(ctx context.Context) ([]Tag, error) {
	return uc.tag.List(ctx)
}

func (uc *TagUseCase) GetTag(ctx context.Context, name string) (*Tag, error) {
	return uc.tag.Get(ctx, name)
}

func (uc *TagUseCase) CreateTag(ctx context.Context, name, comment string) (*Tag, error) {
	return uc.tag.Create(ctx, name, comment)
}

func (uc *TagUseCase) DeleteTag(ctx context.Context, name string) error {
	return uc.tag.Delete(ctx, name)
}
