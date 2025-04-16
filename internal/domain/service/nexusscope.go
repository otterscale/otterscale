package service

import (
	"context"

	"github.com/openhdc/openhdc/internal/domain/model"
)

func (s *NexusService) ListScopes(ctx context.Context) ([]*model.Scope, error) {
	return s.scope.List(ctx)
}

func (s *NexusService) CreateScope(ctx context.Context, name string) (*model.Scope, error) {
	mi, err := s.scope.Create(ctx, name)
	if err != nil {
		return nil, err
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
