package service

import (
	"context"

	"github.com/openhdc/otterscale/internal/domain/model"
)

func (s *NexusService) ListStorageClasses(ctx context.Context, uuid, facility string) ([]model.StorageClass, error) {
	if err := s.setKubernetesClient(ctx, uuid, facility); err != nil {
		return nil, err
	}
	return s.storage.ListStorageClasses(ctx, uuid, facility)
}
