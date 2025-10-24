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

func (s *OrchestratorService) ListEssentials(ctx context.Context, req *pb.ListEssentialsRequest) (*pb.ListEssentialsResponse, error) {
	essentials, err := s.uc.ListEssentials(ctx, core.EssentialType(req.GetType()), req.GetScope())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListEssentialsResponse{}
	resp.SetEssentials(toProtoEssentials(essentials))
	return resp, nil
}

func (s *OrchestratorService) CreateNode(ctx context.Context, req *pb.CreateNodeRequest) (*emptypb.Empty, error) {
	if err := s.uc.CreateNode(ctx, req.GetScope(), req.GetMachineId(), req.GetPrefixName(), req.GetVirtualIps(), req.GetCalicoCidr(), req.GetOsdDevices()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *OrchestratorService) ListKubernetesNodeLabels(ctx context.Context, req *pb.ListKubernetesNodeLabelsRequest) (*pb.ListKubernetesNodeLabelsResponse, error) {
	labels, err := s.uc.ListKubernetesNodeLabels(ctx, req.GetScope(), req.GetFacility(), req.GetHostname(), req.GetAll())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListKubernetesNodeLabelsResponse{}
	resp.SetLabels(labels)
	return resp, nil
}

func (s *OrchestratorService) UpdateKubernetesNodeLabels(ctx context.Context, req *pb.UpdateKubernetesNodeLabelsRequest) (*pb.UpdateKubernetesNodeLabelsResponse, error) {
	labels, err := s.uc.UpdateKubernetesNodeLabels(ctx, req.GetScope(), req.GetFacility(), req.GetHostname(), req.GetLabels())
	if err != nil {
		return nil, err
	}
	resp := &pb.UpdateKubernetesNodeLabelsResponse{}
	resp.SetLabels(labels)
	return resp, nil
}

func (s *OrchestratorService) ListGPURelationsByMachine(ctx context.Context, req *pb.ListGPURelationsByMachineRequest) (*pb.ListGPURelationsByMachineResponse, error) {
	gpuRelations, err := s.uc.ListGPURelationsByMachine(ctx, req.GetScope(), req.GetFacility(), req.GetMachineId())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListGPURelationsByMachineResponse{}
	resp.SetGpuRelations(toProtoGPURelations(gpuRelations))
	return resp, nil
}

func (s *OrchestratorService) ListGPURelationsByModel(ctx context.Context, req *pb.ListGPURelationsByModelRequest) (*pb.ListGPURelationsByModelResponse, error) {
	gpuRelations, err := s.uc.ListGPURelationsByModel(ctx, req.GetScope(), req.GetFacility(), req.GetNamespace(), req.GetModelName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListGPURelationsByModelResponse{}
	resp.SetGpuRelations(toProtoGPURelations(gpuRelations))
	return resp, nil
}

func (s *OrchestratorService) ListGeneralPlugins(ctx context.Context, req *pb.ListGeneralPluginsRequest) (*pb.ListGeneralPluginsResponse, error) {
	plugins, err := s.uc.ListGeneralPlugins(ctx, req.GetScope(), req.GetFacility())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListGeneralPluginsResponse{}
	resp.SetPlugins(toProtoPlugins(plugins))
	return resp, nil
}

func (s *OrchestratorService) ListModelPlugins(ctx context.Context, req *pb.ListModelPluginsRequest) (*pb.ListModelPluginsResponse, error) {
	plugins, err := s.uc.ListModelPlugins(ctx, req.GetScope(), req.GetFacility())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListModelPluginsResponse{}
	resp.SetPlugins(toProtoPlugins(plugins))
	return resp, nil
}

func (s *OrchestratorService) ListInstancePlugins(ctx context.Context, req *pb.ListInstancePluginsRequest) (*pb.ListInstancePluginsResponse, error) {
	plugins, err := s.uc.ListInstancePlugins(ctx, req.GetScope(), req.GetFacility())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListInstancePluginsResponse{}
	resp.SetPlugins(toProtoPlugins(plugins))
	return resp, nil
}

func (s *OrchestratorService) ListStoragePlugins(ctx context.Context, req *pb.ListStoragePluginsRequest) (*pb.ListStoragePluginsResponse, error) {
	plugins, err := s.uc.ListStoragePlugins(ctx, req.GetScope(), req.GetFacility())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListStoragePluginsResponse{}
	resp.SetPlugins(toProtoPlugins(plugins))
	return resp, nil
}

func (s *OrchestratorService) InstallPlugins(ctx context.Context, req *pb.InstallPluginsRequest) (*emptypb.Empty, error) {
	if err := s.uc.InstallPlugins(ctx, req.GetScope(), req.GetFacility(), toChartRefMap(req.GetCharts())); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *OrchestratorService) UpgradePlugins(ctx context.Context, req *pb.UpgradePluginsRequest) (*emptypb.Empty, error) {
	if err := s.uc.UpgradePlugins(ctx, req.GetScope(), req.GetFacility(), toChartRefMap(req.GetCharts())); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func toChartRefMap(cs []*pb.Plugin_Chart) map[string]string {
	ret := map[string]string{}
	for _, c := range cs {
		ret[c.GetName()] = c.GetRef()
	}
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
	ret.SetScope(e.Scope)
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

func toProtoPlugins(ps []core.Plugin) []*pb.Plugin {
	ret := []*pb.Plugin{}
	for i := range ps {
		ret = append(ret, toProtoPlugin(&ps[i]))
	}
	return ret
}

func toProtoPlugin(p *core.Plugin) *pb.Plugin {
	ret := &pb.Plugin{}
	release := p.Release
	if release != nil {
		ret.SetName(release.Name)
		ret.SetNamespace(release.Namespace)
		info := release.Info
		if info != nil {
			ret.SetStatus(info.Status.String())
			ret.SetDescription(info.Description)
			ret.SetFirstDeployedAt(timestamppb.New(info.FirstDeployed.Time))
			ret.SetLastDeployedAt(timestamppb.New(info.LastDeployed.Time))
			if !info.Deleted.IsZero() {
				ret.SetDeletedAt(timestamppb.New(info.Deleted.Time))
			}
		}
		current := release.Chart
		if current != nil && current.Metadata != nil {
			ret.SetCurrent(toProtoPluginChart(current.Metadata))
		}
	}
	latest := p.Latest
	if latest != nil && latest.Metadata != nil {
		ret.SetLatest(toProtoPluginChart(latest.Metadata))
	}
	return ret
}

func toProtoPluginChart(md *core.ChartMetadata) *pb.Plugin_Chart {
	ret := &pb.Plugin_Chart{}
	ret.SetName(md.Name)
	ret.SetVersion(md.Version)
	ret.SetAppVersion(md.AppVersion)
	ret.SetDescription(md.Description)
	ret.SetIcon(md.Icon)
	return ret
}
