package core

import (
	"context"

	"golang.org/x/sync/errgroup"

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

func (uc *TagUseCase) AddMachineTags(ctx context.Context, id string, tags []string) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, tag := range tags {
		eg.Go(func() error {
			return uc.tag.AddMachines(ctx, tag, []string{id})
		})
	}
	return eg.Wait()
}

func (uc *TagUseCase) RemoveMachineTags(ctx context.Context, id string, tags []string) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, tag := range tags {
		eg.Go(func() error {
			return uc.tag.RemoveMachines(ctx, tag, []string{id})
		})
	}
	return eg.Wait()
}
