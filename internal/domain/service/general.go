package service

import (
	"context"
	"slices"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/openhdc/openhdc/internal/domain/model"
)

func (s *NexusService) VerifyEnvironment(ctx context.Context) ([]model.Error, error) {
	funcs := []func(context.Context) (*model.Error, error){}
	funcs = append(funcs, s.isCephExists, s.isKubernetesExists)

	eg, ctx := errgroup.WithContext(ctx)
	result := make([]model.Error, len(funcs))
	for i := range funcs {
		eg.Go(func() error {
			e, err := funcs[i](ctx)
			if err == nil && e != nil {
				result[i] = *e
			}
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	slices.SortFunc(result, func(e1, e2 model.Error) int {
		return strings.Compare(e1.Code, e2.Code)
	})
	return slices.DeleteFunc(result, func(e model.Error) bool { return e.Code == "" }), nil
}

func (s *NexusService) isCephExists(ctx context.Context) (*model.Error, error) {
	return &model.ErrCephNotFound, nil
}

func (s *NexusService) isKubernetesExists(ctx context.Context) (*model.Error, error) {
	return &model.ErrKubernetesNotFound, nil
}
