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

func (s *NexusService) CreateStorageClass(ctx context.Context, kubernetes, ceph *model.FacilityInfo, prefix string) (*model.StorageClass, error) {
	if kubernetes.ScopeUUID != ceph.ScopeUUID {
		return nil, status.Error(codes.Unimplemented, "cross-model integration between facilities is not yet supported")
	}

	leader, err := s.facility.GetLeader(ctx, kubernetes.ScopeUUID, kubernetes.FacilityName)
	if err != nil {
		return nil, err
	}
	unit, err := s.facility.GetUnitInfo(ctx, kubernetes.ScopeUUID, leader)
	if err != nil {
		return nil, err
	}
	jujuToMAASMachineMap, err := s.JujuToMAASMachineMap(ctx, kubernetes.ScopeUUID)
	if err != nil {
		return nil, err
	}
	machineID, ok := jujuToMAASMachineMap[unit.Machine]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "kubernetes %q machine %q not found", kubernetes.FacilityName, unit.Machine)
	}

	configs := map[string]string{}
	fi, err := s.createGeneralFacility(ctx, kubernetes.ScopeUUID, machineID, prefix, charmNameCephCSI, cephCSIFacilityList, configs)
	if err != nil {
		return nil, err
	}
	if err := s.createGeneralRelations(ctx, kubernetes.ScopeUUID, toCephCSIEndpointList(kubernetes, ceph, prefix)); err != nil {
		return nil, err
	}

	return s.storage.GetStorageClass(ctx, kubernetes.ScopeUUID, kubernetes.FacilityName, fi.FacilityName)
}
