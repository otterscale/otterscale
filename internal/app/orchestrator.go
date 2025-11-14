package app

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/otterscale/otterscale/api/orchestrator/v1"
	"github.com/otterscale/otterscale/api/orchestrator/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core/application/chart"
	"github.com/otterscale/otterscale/internal/core/machine"
	"github.com/otterscale/otterscale/internal/core/orchestrator"
	"github.com/otterscale/otterscale/internal/core/orchestrator/extension"
	"github.com/otterscale/otterscale/internal/core/orchestrator/gpu"
	"github.com/otterscale/otterscale/internal/core/orchestrator/standalone"
)

type OrchestratorService struct {
	pbconnect.UnimplementedOrchestratorServiceHandler

	orchestrator *orchestrator.OrchestratorUseCase
	extension    *extension.ExtensionUseCase
	gpu          *gpu.GPUUseCase
	standalone   *standalone.StandaloneUseCase
}

func NewOrchestratorService(orchestrator *orchestrator.OrchestratorUseCase, extension *extension.ExtensionUseCase, gpu *gpu.GPUUseCase, standalone *standalone.StandaloneUseCase) *OrchestratorService {
	return &OrchestratorService{
		orchestrator: orchestrator,
		extension:    extension,
		gpu:          gpu,
		standalone:   standalone,
	}
}

var _ pbconnect.OrchestratorServiceHandler = (*OrchestratorService)(nil)

func (s *OrchestratorService) CreateNode(ctx context.Context, req *pb.CreateNodeRequest) (*emptypb.Empty, error) {
	if err := s.standalone.CreateNode(ctx, req.GetScope(), req.GetMachineId(), req.GetVirtualIps(), req.GetCalicoCidr(), req.GetOsdDevices()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *OrchestratorService) ListKubernetesNodeLabels(ctx context.Context, req *pb.ListKubernetesNodeLabelsRequest) (*pb.ListKubernetesNodeLabelsResponse, error) {
	labels, err := s.orchestrator.ListKubernetesNodeLabels(ctx, req.GetScope(), req.GetHostname(), req.GetAll())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListKubernetesNodeLabelsResponse{}
	resp.SetLabels(labels)
	return resp, nil
}

func (s *OrchestratorService) UpdateKubernetesNodeLabels(ctx context.Context, req *pb.UpdateKubernetesNodeLabelsRequest) (*pb.UpdateKubernetesNodeLabelsResponse, error) {
	labels, err := s.orchestrator.UpdateKubernetesNodeLabels(ctx, req.GetScope(), req.GetHostname(), req.GetLabels())
	if err != nil {
		return nil, err
	}

	resp := &pb.UpdateKubernetesNodeLabelsResponse{}
	resp.SetLabels(labels)
	return resp, nil
}

func (s *OrchestratorService) ListGPURelationsByMachine(ctx context.Context, req *pb.ListGPURelationsByMachineRequest) (*pb.ListGPURelationsByMachineResponse, error) {
	relations, err := s.gpu.ListGPURelationsByMachine(ctx, req.GetScope(), req.GetMachineId())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListGPURelationsByMachineResponse{}
	resp.SetGpuRelations(toProtoGPURelations(relations))
	return resp, nil
}

func (s *OrchestratorService) ListGPURelationsByModel(ctx context.Context, req *pb.ListGPURelationsByModelRequest) (*pb.ListGPURelationsByModelResponse, error) {
	relations, err := s.gpu.ListGPURelationsByModel(ctx, req.GetScope(), req.GetNamespace(), req.GetModelName())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListGPURelationsByModelResponse{}
	resp.SetGpuRelations(toProtoGPURelations(relations))
	return resp, nil
}

func (s *OrchestratorService) ListGeneralExtensions(ctx context.Context, req *pb.ListGeneralExtensionsRequest) (*pb.ListGeneralExtensionsResponse, error) {
	extensions, err := s.extension.ListGeneralExtensions(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListGeneralExtensionsResponse{}
	resp.SetExtensions(toProtoExtensions(extensions))
	return resp, nil
}

func (s *OrchestratorService) ListModelExtensions(ctx context.Context, req *pb.ListModelExtensionsRequest) (*pb.ListModelExtensionsResponse, error) {
	extensions, err := s.extension.ListModelExtensions(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListModelExtensionsResponse{}
	resp.SetExtensions(toProtoExtensions(extensions))
	return resp, nil
}

func (s *OrchestratorService) ListInstanceExtensions(ctx context.Context, req *pb.ListInstanceExtensionsRequest) (*pb.ListInstanceExtensionsResponse, error) {
	extensions, err := s.extension.ListInstanceExtensions(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListInstanceExtensionsResponse{}
	resp.SetExtensions(toProtoExtensions(extensions))
	return resp, nil
}

func (s *OrchestratorService) ListStorageExtensions(ctx context.Context, req *pb.ListStorageExtensionsRequest) (*pb.ListStorageExtensionsResponse, error) {
	extensions, err := s.extension.ListStorageExtensions(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListStorageExtensionsResponse{}
	resp.SetExtensions(toProtoExtensions(extensions))
	return resp, nil
}

func (s *OrchestratorService) InstallExtensions(ctx context.Context, req *pb.InstallExtensionsRequest) (*emptypb.Empty, error) {
	if err := s.extension.InstallExtensions(ctx, req.GetScope(), toChartRefMap(req.GetCharts())); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *OrchestratorService) UpgradeExtensions(ctx context.Context, req *pb.UpgradeExtensionsRequest) (*emptypb.Empty, error) {
	if err := s.extension.UpgradeExtensions(ctx, req.GetScope(), toChartRefMap(req.GetCharts())); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func toChartRefMap(cs []*pb.Extension_Chart) map[string]string {
	ret := map[string]string{}

	for _, c := range cs {
		ret[c.GetName()] = c.GetRef()
	}

	return ret
}

func toProtoGPURelations(rs *gpu.Relations) []*pb.GPURelation {
	ret := []*pb.GPURelation{}

	for i := range rs.Machines {
		ret = append(ret, toProtoGPURelationFromMachine(&rs.Machines[i]))
	}

	for i := range rs.GPUs {
		ret = append(ret, toProtoGPURelationFromGPU(&rs.GPUs[i]))
	}

	for i := range rs.Pods {
		ret = append(ret, toProtoGPURelationFromPod(&rs.Pods[i]))
	}

	return ret
}

func toProtoGPURelationFromMachine(m *machine.Machine) *pb.GPURelation {
	machine := &pb.GPURelation_Machine{}
	machine.SetId(m.SystemID)
	machine.SetHostname(m.Hostname)

	ret := &pb.GPURelation{}
	ret.SetMachine(machine)
	return ret
}

func toProtoGPURelationFromGPU(g *gpu.GPURelation) *pb.GPURelation {
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

func toProtoGPURelationFromPod(p *gpu.PodRelation) *pb.GPURelation {
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

func toProtoGPURelationPodDevices(pds []gpu.PodDevice) []*pb.GPURelation_Pod_Device {
	ret := []*pb.GPURelation_Pod_Device{}

	for i := range pds {
		ret = append(ret, toProtoGPURelationPodDevice(&pds[i]))
	}

	return ret
}

func toProtoGPURelationPodDevice(pd *gpu.PodDevice) *pb.GPURelation_Pod_Device {
	ret := &pb.GPURelation_Pod_Device{}
	ret.SetGpuId(pd.GPUID)
	ret.SetUsedCores(pd.UsedCores)
	ret.SetUsedMemoryBytes(pd.UsedMemoryBytes)
	return ret
}

func toProtoExtensions(ps []extension.Extension) []*pb.Extension {
	ret := []*pb.Extension{}

	for i := range ps {
		ret = append(ret, toProtoExtension(&ps[i]))
	}

	return ret
}

func toProtoExtension(p *extension.Extension) *pb.Extension {
	ret := &pb.Extension{}

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
			ret.SetCurrent(toProtoExtensionChart(current.Metadata, ""))
		}
	}

	latest := p.Latest

	if latest != nil && latest.Metadata != nil {
		ret.SetLatest(toProtoExtensionChart(latest.Metadata, p.LatestURL))
	}

	return ret
}

func toProtoExtensionChart(md *chart.Metadata, ref string) *pb.Extension_Chart {
	ret := &pb.Extension_Chart{}
	ret.SetName(md.Name)
	ret.SetVersion(md.Version)
	ret.SetAppVersion(md.AppVersion)
	ret.SetDescription(md.Description)
	ret.SetIcon(md.Icon)
	ret.SetRef(ref)
	return ret
}
