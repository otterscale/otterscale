package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openhdc/openhdc/internal/domain/model"
)

func (s *NexusService) ListStorageClasses(ctx context.Context, uuid, facility string) ([]model.StorageClass, error) {
	if err := s.setKubernetesClient(ctx, uuid, facility); err != nil {
		return nil, err
	}
	return s.storage.ListStorageClasses(ctx, uuid, facility)
}

func (s *NexusService) CreateStorageClass(ctx context.Context, kubernetes, ceph *model.FacilityInfo, prefix string) (*model.StorageClass, error) {
	if kubernetes.ScopeUUID != ceph.ScopeUUID {
		return nil, status.Error(codes.Unimplemented, "cross-model integration between facilities is not yet supported")
	}

	configs, err := getCephCSIConfigs(prefix)
	if err != nil {
		return nil, err
	}
	if _, err := s.createGeneralFacility(ctx, kubernetes.ScopeUUID, "", prefix, charmNameCephCSI, cephCSIFacilityList, configs); err != nil {
		return nil, err
	}
	if err := s.createGeneralRelations(ctx, kubernetes.ScopeUUID, toCephCSIEndpointList(kubernetes, ceph, prefix)); err != nil {
		return nil, err
	}

	if err := s.setKubernetesClient(ctx, kubernetes.ScopeUUID, kubernetes.FacilityName); err != nil {
		return nil, err
	}
	return s.storage.GetStorageClass(ctx, kubernetes.ScopeUUID, kubernetes.FacilityName, defaultStorage)
}
