package app

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/openhdc/openhdc/api/nexus/v1"
	"github.com/openhdc/openhdc/internal/domain/model"
)

func (a *NexusApp) ListStorageClasses(ctx context.Context, req *connect.Request[pb.ListStorageClassesRequest]) (*connect.Response[pb.ListStorageClassesResponse], error) {
	scs, err := a.svc.ListStorageClasses(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	res := &pb.ListStorageClassesResponse{}
	res.SetStorageClasses(toProtoStorageClasses(scs))
	return connect.NewResponse(res), nil
}

func toProtoStorageClasses(scs []model.StorageClass) []*pb.StorageClass {
	ret := []*pb.StorageClass{}
	for i := range scs {
		ret = append(ret, toProtoStorageClass(&scs[i]))
	}
	return ret
}

func toProtoStorageClass(sc *model.StorageClass) *pb.StorageClass {
	reclaimPolicy := ""
	if v := sc.ReclaimPolicy; v != nil {
		reclaimPolicy = string(*v)
	}
	volumeBindingMode := ""
	if v := sc.VolumeBindingMode; v != nil {
		volumeBindingMode = string(*v)
	}
	ret := &pb.StorageClass{}
	ret.SetName(sc.Name)
	ret.SetProvisioner(sc.Provisioner)
	ret.SetReclaimPolicy(reclaimPolicy)
	ret.SetVolumeBindingMode(volumeBindingMode)
	ret.SetParameters(sc.Parameters)
	ret.SetCreatedAt(timestamppb.New(sc.CreationTimestamp.Time))
	return ret
}
