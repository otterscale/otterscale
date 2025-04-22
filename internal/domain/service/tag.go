package service

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/openhdc/openhdc/internal/domain/model"
)

func (s *NexusService) ListTags(ctx context.Context) ([]model.Tag, error) {
	return s.tag.List(ctx)
}

func (s *NexusService) GetTag(ctx context.Context, name string) (*model.Tag, error) {
	return s.tag.Get(ctx, name)
}

func (s *NexusService) CreateTag(ctx context.Context, name, comment string) (*model.Tag, error) {
	return s.tag.Create(ctx, name, comment)
}

func (s *NexusService) DeleteTag(ctx context.Context, name string) error {
	return s.tag.Delete(ctx, name)
}

func (s *NexusService) AddMachineTags(ctx context.Context, id string, tags []string) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, tag := range tags {
		eg.Go(func() error {
			return s.tag.AddMachines(ctx, tag, []string{id})
		})
	}
	return eg.Wait()
}

func (s *NexusService) RemoveMachineTags(ctx context.Context, id string, tags []string) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, tag := range tags {
		eg.Go(func() error {
			return s.tag.RemoveMachines(ctx, tag, []string{id})
		})
	}
	return eg.Wait()
}
