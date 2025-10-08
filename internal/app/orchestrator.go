package app

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/otterscale/otterscale/api/orchestrator/v1"
	"github.com/otterscale/otterscale/api/orchestrator/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type OrchestratorService struct {
	pbconnect.UnimplementedOrchestratorServiceHandler

	uc *core.OrchestratorUseCase
}

func NewOrchestratorService(uc *core.OrchestratorUseCase) *OrchestratorService {
	return &OrchestratorService{uc: uc}
}

var _ pbconnect.OrchestratorServiceHandler = (*OrchestratorService)(nil)

func (s *OrchestratorService) IsMachineDeployed(ctx context.Context, req *pb.IsMachineDeployedRequest) (*pb.IsMachineDeployedResponse, error) {
	message, deployed, err := s.uc.IsMachineDeployed(ctx, req.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	resp := &pb.IsMachineDeployedResponse{}
	resp.SetMessage(message)
	resp.SetDeployed(deployed)
	return resp, nil
}

func (s *OrchestratorService) ListStatuses(ctx context.Context, req *pb.ListStatusesRequest) (*pb.ListStatusesResponse, error) {
	statuses, err := s.uc.ListStatuses(ctx, req.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListStatusesResponse{}
	resp.SetStatuses(toProtoStatuses(statuses))
	return resp, nil
}

func (s *OrchestratorService) ListEssentials(ctx context.Context, req *pb.ListEssentialsRequest) (*pb.ListEssentialsResponse, error) {
	essentials, err := s.uc.ListEssentials(ctx, int32(req.GetType()), req.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListEssentialsResponse{}
	resp.SetEssentials(toProtoEssentials(essentials))
	return resp, nil
}

func (s *OrchestratorService) CreateSingleNode(ctx context.Context, req *pb.CreateSingleNodeRequest) (*emptypb.Empty, error) {
	if err := s.uc.CreateSingleNode(ctx, req.GetScopeUuid(), req.GetMachineId(), req.GetPrefixName(), req.GetVirtualIps(), req.GetCalicoCidr(), req.GetOsdDevices()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *OrchestratorService) ListKubernetesNodeLabels(ctx context.Context, req *pb.ListKubernetesNodeLabelsRequest) (*pb.ListKubernetesNodeLabelsResponse, error) {
	labels, err := s.uc.ListKubernetesNodeLabels(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetHostname(), req.GetAll())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListKubernetesNodeLabelsResponse{}
	resp.SetLabels(labels)
	return resp, nil
}

func (s *OrchestratorService) UpdateKubernetesNodeLabels(ctx context.Context, req *pb.UpdateKubernetesNodeLabelsRequest) (*pb.UpdateKubernetesNodeLabelsResponse, error) {
	labels, err := s.uc.UpdateKubernetesNodeLabels(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetHostname(), req.GetLabels())
	if err != nil {
		return nil, err
	}
	resp := &pb.UpdateKubernetesNodeLabelsResponse{}
	resp.SetLabels(labels)
	return resp, nil
}

func (s *OrchestratorService) ListGPURelationsByMachine(ctx context.Context, req *pb.ListGPURelationsByMachineRequest) (*pb.ListGPURelationsByMachineResponse, error) {
	gpuRelations, err := s.uc.ListGPURelationsByMachine(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetMachineId())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListGPURelationsByMachineResponse{}
	resp.SetGpuRelations(toProtoGPURelations(gpuRelations))
	return resp, nil
}

func (s *OrchestratorService) ListGPURelationsByModel(ctx context.Context, req *pb.ListGPURelationsByModelRequest) (*pb.ListGPURelationsByModelResponse, error) {
	gpuRelations, err := s.uc.ListGPURelationsByModel(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetModelName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListGPURelationsByModelResponse{}
	resp.SetGpuRelations(toProtoGPURelations(gpuRelations))
	return resp, nil
}

func toProtoStatuses(ess []core.EssentialStatus) []*pb.Status {
	ret := []*pb.Status{}
	for i := range ess {
		ret = append(ret, toProtoStatus(&ess[i]))
	}
	return ret
}

func toProtoStatus(es *core.EssentialStatus) *pb.Status {
	ret := &pb.Status{}
	ret.SetLevel(pb.Status_Level(es.Level))
	ret.SetMessage(es.Message)
	ret.SetDetails(es.Details)
	return ret
}

func toProtoEssentials(es []core.Essential) []*pb.Essential {
	ret := []*pb.Essential{}
	for i := range es {
		ret = append(ret, toProtoEssential(&es[i]))
	}
	return ret
}

func toProtoEssential(e *core.Essential) *pb.Essential {
	ret := &pb.Essential{}
	ret.SetType(pb.Essential_Type(e.Type))
	ret.SetName(e.Name)
	ret.SetScopeUuid(e.ScopeUUID)
	ret.SetScopeName(e.ScopeName)
	ret.SetUnits(toProtoEssentialUnits(e.Units))
	return ret
}

func toProtoEssentialUnits(eus []core.EssentialUnit) []*pb.Essential_Unit {
	ret := []*pb.Essential_Unit{}
	for i := range eus {
		ret = append(ret, toProtoEssentialUnit(&eus[i]))
	}
	return ret
}

func toProtoEssentialUnit(eu *core.EssentialUnit) *pb.Essential_Unit {
	ret := &pb.Essential_Unit{}
	ret.SetName(eu.Name)
	ret.SetDirective(eu.Directive)
	return ret
}

func toProtoGPURelations(rs *core.GPURelations) []*pb.GPURelation {
	ret := []*pb.GPURelation{}
	for _, machine := range rs.Machines {
		ret = append(ret, toProtoGPURelationFromMachine(&machine))
	}
	for _, gpu := range rs.GPUs {
		ret = append(ret, toProtoGPURelationFromGPU(&gpu))
	}
	for _, pod := range rs.Pods {
		ret = append(ret, toProtoGPURelationFromPod(&pod))
	}
	return ret
}

func toProtoGPURelationFromMachine(m *core.Machine) *pb.GPURelation {
	machine := &pb.GPURelation_Machine{}
	machine.SetId(m.SystemID)
	machine.SetHostname(m.Hostname)

	ret := &pb.GPURelation{}
	ret.SetMachine(machine)
	return ret
}

func toProtoGPURelationFromGPU(g *core.GPURelationsGPU) *pb.GPURelation {
	gpu := &pb.GPURelation_GPU{}
	gpu.SetId(g.ID)
	gpu.SetIndex(g.Index)
	gpu.SetCount(g.Count)
	gpu.SetCores(g.Cores)
	gpu.SetMemoryBytes(g.MemoryBytes)
	gpu.SetType(g.Type)
	gpu.SetHealth(g.Health)
	gpu.SetMachineId(g.MachineID)

	ret := &pb.GPURelation{}
	ret.SetGpu(gpu)
	return ret
}

func toProtoGPURelationFromPod(p *core.GPURelationsPod) *pb.GPURelation {
	pod := &pb.GPURelation_Pod{}
	pod.SetName(p.Name)
	pod.SetNamespace(p.Namespace)
	pod.SetModelName(p.ModelName)
	pod.SetBindingPhase(p.BindingPhase)
	if !p.BoundAt.IsZero() {
		pod.SetBoundAt(timestamppb.New(p.BoundAt))
	}
	pod.SetDevices(toProtoGPURelationPodDevices(p.PodDevices))

	ret := &pb.GPURelation{}
	ret.SetPod(pod)
	return ret
}

func toProtoGPURelationPodDevices(pds []core.GPURelationPodDevice) []*pb.GPURelation_Pod_Device {
	ret := []*pb.GPURelation_Pod_Device{}
	for i := range pds {
		ret = append(ret, toProtoGPURelationPodDevice(&pds[i]))
	}
	return ret
}

func toProtoGPURelationPodDevice(pd *core.GPURelationPodDevice) *pb.GPURelation_Pod_Device {
	ret := &pb.GPURelation_Pod_Device{}
	ret.SetGpuId(pd.GPUID)
	ret.SetUsedCores(pd.UsedCores)
	ret.SetUsedMemoryBytes(pd.UsedMemoryBytes)
	return ret
}
