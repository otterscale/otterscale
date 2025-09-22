package app

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/otterscale/otterscale/api/virtual_machine/v1"
	"github.com/otterscale/otterscale/api/virtual_machine/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
	"google.golang.org/protobuf/types/known/emptypb"
)

type VirtualMachineService struct {
	pbconnect.UnimplementedVirtualMachineServiceHandler

	uc *core.VirtualMachineUseCase
}

func NewVirtualMachineService(uc *core.VirtualMachineUseCase) *VirtualMachineService {
	return &VirtualMachineService{uc: uc}
}

var _ pbconnect.VirtualMachineServiceHandler = (*VirtualMachineService)(nil)

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
