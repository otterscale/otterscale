package service

import (
	"context"

	"github.com/openhdc/openhdc/internal/domain/model"
)

const defaultScopeName = "default"

func (s *NexusService) ListScopes(ctx context.Context) ([]model.Scope, error) {
	return s.scope.List(ctx)
}

func (s *NexusService) CreateScope(ctx context.Context, name string) (*model.Scope, error) {
	sshKey, err := s.sshKey.Default(ctx)
	if err != nil {
		return nil, err
	}
	mi, err := s.scope.Create(ctx, name)
	if err != nil {
		return nil, err
	}
	results, err := s.keyManager.Add(ctx, mi.UUID, sshKey.Key)
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		if result.Error != nil {
			return nil, result.Error
		}
	}
	return &model.Scope{
		Name:            mi.Name,
		UUID:            mi.UUID,
		Type:            mi.Type,
		ControllerUUID:  mi.ControllerUUID,
		IsController:    mi.IsController,
		ProviderType:    mi.ProviderType,
		Cloud:           mi.Cloud,
		CloudRegion:     mi.CloudRegion,
		CloudCredential: mi.CloudCredential,
		Owner:           mi.Owner,
		Life:            mi.Life,
		Status:          mi.Status,
		AgentVersion:    mi.AgentVersion,
	}, nil
}

func (s *NexusService) CreateDefaultScope(ctx context.Context) (*model.Scope, error) {
	return s.CreateScope(ctx, defaultScopeName)
}
