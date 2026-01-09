package app

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/otterscale/otterscale/api/orchestrator/v1"
	"github.com/otterscale/otterscale/api/orchestrator/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core/machine"
	"github.com/otterscale/otterscale/internal/core/orchestrator"
	"github.com/otterscale/otterscale/internal/core/orchestrator/extension"
	"github.com/otterscale/otterscale/internal/core/orchestrator/gpu"
	"github.com/otterscale/otterscale/internal/core/orchestrator/standalone"
)

type OrchestratorService struct {
	pbconnect.UnimplementedOrchestratorServiceHandler

	orchestrator *orchestrator.UseCase
	extension    *extension.UseCase
	gpu          *gpu.UseCase
	standalone   *standalone.UseCase
}

func NewOrchestratorService(orchestrator *orchestrator.UseCase, extension *extension.UseCase, gpu *gpu.UseCase, standalone *standalone.UseCase) *OrchestratorService {
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

func (s *OrchestratorService) ListExtensions(ctx context.Context, req *pb.ListExtensionsRequest) (*pb.ListExtensionsResponse, error) {
	extensions, err := s.extension.ListExtensions(ctx, req.GetScope(), toExtensionType(req.GetType()))
	if err != nil {
		return nil, err
	}

	resp := &pb.ListExtensionsResponse{}
	resp.SetExtensions(toProtoExtensions(extensions))
	return resp, nil
}

func (s *OrchestratorService) InstallOrUpgradeExtensions(ctx context.Context, req *pb.InstallOrUpgradeExtensionsRequest) (*emptypb.Empty, error) {
	if err := s.extension.InstallOrUpgradeExtensions(ctx, req.GetScope(), toManifests(req.GetManifests()), toExtensionArguments(req.GetArguments())); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func toManifests(ms []*pb.Extension_Manifest) []extension.Manifest {
	ret := []extension.Manifest{}

	for _, m := range ms {
		ret = append(ret, extension.Manifest{
			ID:      m.GetId(),
			Version: m.GetVersion(),
		})
	}

	return ret
}

func toExtensionArguments(args map[string]*pb.InstallOrUpgradeExtensionsRequest_ExtensionArguments) map[string]map[string]string {
	ret := make(map[string]map[string]string)

	for extensionID, extArgs := range args {
		if extArgs != nil && len(extArgs.GetValues()) > 0 {
			ret[extensionID] = extArgs.GetValues()
		}
	}

	return ret
}

func toExtensionType(t pb.Extension_Type) extension.Type {
	switch t {
	case pb.Extension_TYPE_METRICS:
		return extension.TypeMetrics

	case pb.Extension_TYPE_SERVICE_MESH:
		return extension.TypeServiceMesh

	case pb.Extension_TYPE_REGISTRY:
		return extension.TypeRegistry

	case pb.Extension_TYPE_MODEL:
		return extension.TypeModel

	case pb.Extension_TYPE_INSTANCE:
		return extension.TypeInstance

	case pb.Extension_TYPE_STORAGE:
		return extension.TypeStorage

	default:
		return extension.TypeUnspecified
	}
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
	ret.SetDisplayName(p.DisplayName)
	ret.SetDescription(p.Description)
	ret.SetIcon(p.Icon)
	ret.SetStatus(p.Status)
	if p.DeployedAt != nil {
		ret.SetDeployedAt(timestamppb.New(*p.DeployedAt))
	}
	if p.Current != nil {
		ret.SetCurrent(toProtoExtensionManifest(p.Current))
	}
	if p.Latest != nil {
		ret.SetLatest(toProtoExtensionManifest(p.Latest))
	}
	return ret
}

func toProtoExtensionManifest(m *extension.Manifest) *pb.Extension_Manifest {
	ret := &pb.Extension_Manifest{}
	ret.SetId(m.ID)
	ret.SetVersion(m.Version)
	return ret
}
