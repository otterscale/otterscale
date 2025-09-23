package app

import (
	"context"

	"connectrpc.com/connect"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	corev1 "k8s.io/api/core/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"

	pb "github.com/otterscale/otterscale/api/virtual_machine/v1"
	"github.com/otterscale/otterscale/api/virtual_machine/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type VirtualMachineService struct {
	pbconnect.UnimplementedVirtualMachineServiceHandler

	uc *core.VirtualMachineUseCase
}

func NewVirtualMachineService(uc *core.VirtualMachineUseCase) *VirtualMachineService {
	return &VirtualMachineService{uc: uc}
}

var _ pbconnect.VirtualMachineServiceHandler = (*VirtualMachineService)(nil)

func (s *VirtualMachineService) ListDataVolumes(ctx context.Context, req *connect.Request[pb.ListDataVolumesRequest]) (*connect.Response[pb.ListDataVolumesResponse], error) {
	its, err := s.uc.ListDataVolumes(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListDataVolumesResponse{}
	resp.SetDataVolumes(toProtoDataVolumes(its))
	return connect.NewResponse(resp), nil
}

func (s *VirtualMachineService) GetDataVolume(ctx context.Context, req *connect.Request[pb.GetDataVolumeRequest]) (*connect.Response[pb.DataVolume], error) {
	it, err := s.uc.GetDataVolume(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoDataVolume(it)
	return connect.NewResponse(resp), nil
}

func (s *VirtualMachineService) CreateDataVolume(ctx context.Context, req *connect.Request[pb.CreateDataVolumeRequest]) (*connect.Response[pb.DataVolume], error) {
	src := req.Msg.GetSource()
	it, err := s.uc.CreateDataVolume(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), core.SourceType(src.GetType()), src.GetData(), req.Msg.GetSizeBytes(), req.Msg.GetBootImage())
	if err != nil {
		return nil, err
	}
	resp := toProtoDataVolume(it)
	return connect.NewResponse(resp), nil
}

func (s *VirtualMachineService) DeleteDataVolume(ctx context.Context, req *connect.Request[pb.DeleteDataVolumeRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeleteDataVolume(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *VirtualMachineService) ExtendDataVolume(ctx context.Context, req *connect.Request[pb.ExtendDataVolumeRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.ExtendDataVolume(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetNewSizeBytes()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *VirtualMachineService) ListClusterWideInstanceTypes(ctx context.Context, req *connect.Request[pb.ListClusterWideInstanceTypesRequest]) (*connect.Response[pb.ListClusterWideInstanceTypesResponse], error) {
	its, err := s.uc.ListClusterWideInstanceTypes(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListClusterWideInstanceTypesResponse{}
	resp.SetInstanceTypes(toProtoClusterInstanceTypes(its))
	return connect.NewResponse(resp), nil
}

func (s *VirtualMachineService) ListInstanceTypes(ctx context.Context, req *connect.Request[pb.ListInstanceTypesRequest]) (*connect.Response[pb.ListInstanceTypesResponse], error) {
	its, err := s.uc.ListInstanceTypes(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListInstanceTypesResponse{}
	resp.SetInstanceTypes(toProtoInstanceTypes(its))
	return connect.NewResponse(resp), nil
}

func (s *VirtualMachineService) GetInstanceType(ctx context.Context, req *connect.Request[pb.GetInstanceTypeRequest]) (*connect.Response[pb.InstanceType], error) {
	it, err := s.uc.GetInstanceType(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoInstanceType(it)
	return connect.NewResponse(resp), nil
}

func (s *VirtualMachineService) CreateInstanceType(ctx context.Context, req *connect.Request[pb.CreateInstanceTypeRequest]) (*connect.Response[pb.InstanceType], error) {
	it, err := s.uc.CreateInstanceType(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetCpuCores(), req.Msg.GetMemoryBytes())
	if err != nil {
		return nil, err
	}
	resp := toProtoInstanceType(it)
	return connect.NewResponse(resp), nil
}

func (s *VirtualMachineService) DeleteInstanceType(ctx context.Context, req *connect.Request[pb.DeleteInstanceTypeRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeleteInstanceType(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func getDataVolumeSize(spec cdiv1beta1.DataVolumeSpec) int64 {
	if spec.PVC != nil {
		if size := extractStorageSize(spec.PVC.Resources.Requests); size > 0 {
			return size
		}
	}
	if spec.Storage != nil {
		if size := extractStorageSize(spec.Storage.Resources.Requests); size > 0 {
			return size
		}
	}
	return 0
}

func extractStorageSize(requests corev1.ResourceList) int64 {
	if requests == nil {
		return 0
	}
	if s, ok := requests[corev1.ResourceStorage]; ok {
		if v, ok := s.AsInt64(); ok {
			return v
		}
	}
	return 0
}

func toProtoDataVolumeSource(s *cdiv1beta1.DataVolumeSource) *pb.DataVolume_Source {
	ret := &pb.DataVolume_Source{}
	switch {
	case s.Blank != nil:
		ret.SetType(pb.DataVolume_Source_BLANK_IMAGE)
		ret.SetData("")
	case s.HTTP != nil:
		ret.SetType(pb.DataVolume_Source_HTTP_URL)
		ret.SetData(s.HTTP.URL)
	case s.PVC != nil:
		ret.SetType(pb.DataVolume_Source_EXISTING_PERSISTENT_VOLUME_CLAIM)
		ret.SetData(s.PVC.Name)
	}
	return ret
}

func toProtoDataVolumeCondition(c *cdiv1beta1.DataVolumeCondition) *pb.DataVolume_Condition {
	ret := &pb.DataVolume_Condition{}
	ret.SetType(string(c.Type))
	ret.SetStatus(string(c.Status))
	ret.SetReason((c.Reason))
	ret.SetMessage((c.Message))
	if !c.LastHeartbeatTime.IsZero() {
		ret.SetHeartbeatAt(timestamppb.New(c.LastHeartbeatTime.Time))
	}
	if !c.LastTransitionTime.IsZero() {
		ret.SetTransitionedAt(timestamppb.New(c.LastTransitionTime.Time))
	}
	return ret
}

func toProtoDataVolumes(its []core.DataVolumePVC) []*pb.DataVolume {
	ret := []*pb.DataVolume{}
	for i := range its {
		ret = append(ret, toProtoDataVolume(&its[i]))
	}
	return ret
}

func toProtoDataVolume(it *core.DataVolumePVC) *pb.DataVolume {
	ret := &pb.DataVolume{}
	ret.SetName(it.Name)
	ret.SetNamespace(it.Namespace)
	ret.SetSource(toProtoDataVolumeSource(it.Spec.Source))
	ret.SetBootImage(it.Labels[core.DataVolumeBootImageLabel] == "true")
	ret.SetPhase(string(it.Status.Phase))
	ret.SetProgress(string(it.Status.Progress))
	ret.SetSizeBytes(getDataVolumeSize(it.Spec))
	if it.Storage != nil {
		ret.SetPersistentVolumeClaim(toProtoPersistentVolumeClaim(it.Storage))
	}
	if len(it.Status.Conditions) > 0 {
		index := len(it.Status.Conditions) - 1
		ret.SetLastCondition(toProtoDataVolumeCondition(&it.Status.Conditions[index]))
	}
	return ret
}

func toProtoClusterInstanceTypes(its []core.VirtualMachineClusterInstanceType) []*pb.InstanceType {
	ret := []*pb.InstanceType{}
	for i := range its {
		ret = append(ret, toProtoClusterInstanceType(&its[i]))
	}
	return ret
}

func toProtoClusterInstanceType(it *core.VirtualMachineClusterInstanceType) *pb.InstanceType {
	ret := &pb.InstanceType{}
	ret.SetName(it.Name)
	ret.SetNamespace(it.Namespace)
	ret.SetCpuCores(it.Spec.CPU.Guest)
	ret.SetMemoryBytes(it.Spec.Memory.Guest.Value())
	return ret
}

func toProtoInstanceTypes(its []core.VirtualMachineInstanceType) []*pb.InstanceType {
	ret := []*pb.InstanceType{}
	for i := range its {
		ret = append(ret, toProtoInstanceType(&its[i]))
	}
	return ret
}

func toProtoInstanceType(it *core.VirtualMachineInstanceType) *pb.InstanceType {
	ret := &pb.InstanceType{}
	ret.SetName(it.Name)
	ret.SetNamespace(it.Namespace)
	ret.SetCpuCores(it.Spec.CPU.Guest)
	ret.SetMemoryBytes(it.Spec.Memory.Guest.Value())
	return ret
}
