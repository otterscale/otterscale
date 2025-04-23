package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openhdc/openhdc/internal/domain/model"
)

func (s *NexusService) ListStorageClasses(ctx context.Context, uuid, facility string) ([]model.StorageClass, error) {
	return s.storage.ListStorageClasses(ctx, uuid, facility)
}

func (s *NexusService) CreateStorageClass(ctx context.Context, kubernetes, ceph *model.FacilityInfo) (*model.StorageClass, error) {
	if kubernetes.ScopeUUID != ceph.ScopeUUID {
		return nil, status.Error(codes.Unimplemented, "cross-model integration between facilities is not yet supported")
	}
	// check pool list
	// deploy ceph-csi
	// get storage class
	return nil, status.Error(codes.Unimplemented, "cross-model integration between facilities is not yet supported")
}
