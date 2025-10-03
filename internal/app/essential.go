package app

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/otterscale/otterscale/api/essential/v1"
	"github.com/otterscale/otterscale/api/essential/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type EssentialService struct {
	pbconnect.UnimplementedEssentialServiceHandler

	uc *core.EssentialUseCase
}

func NewEssentialService(uc *core.EssentialUseCase) *EssentialService {
	return &EssentialService{uc: uc}
}

var _ pbconnect.EssentialServiceHandler = (*EssentialService)(nil)

func (s *EssentialService) IsMachineDeployed(ctx context.Context, req *pb.IsMachineDeployedRequest) (*pb.IsMachineDeployedResponse, error) {
	message, deployed, err := s.uc.IsMachineDeployed(ctx, req.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	resp := &pb.IsMachineDeployedResponse{}
	resp.SetMessage(message)
	resp.SetDeployed(deployed)
	return resp, nil
}

func (s *EssentialService) ListStatuses(ctx context.Context, req *pb.ListStatusesRequest) (*pb.ListStatusesResponse, error) {
	statuses, err := s.uc.ListStatuses(ctx, req.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListStatusesResponse{}
	resp.SetStatuses(toProtoStatuses(statuses))
	return resp, nil
}

func (s *EssentialService) ListEssentials(ctx context.Context, req *pb.ListEssentialsRequest) (*pb.ListEssentialsResponse, error) {
	essentials, err := s.uc.ListEssentials(ctx, int32(req.GetType()), req.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListEssentialsResponse{}
	resp.SetEssentials(toProtoEssentials(essentials))
	return resp, nil
}

func (s *EssentialService) CreateSingleNode(ctx context.Context, req *pb.CreateSingleNodeRequest) (*emptypb.Empty, error) {
	if err := s.uc.CreateSingleNode(ctx,
		req.GetScopeUuid(),
		req.GetMachineId(),
		req.GetPrefixName(),
		req.GetVirtualIps(),
		req.GetCalicoCidr(),
		req.GetOsdDevices(),
	); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *EssentialService) ListKubernetesNodeLabels(ctx context.Context, req *pb.ListKubernetesNodeLabelsRequest) (*pb.ListKubernetesNodeLabelsResponse, error) {
	labels, err := s.uc.ListKubernetesNodeLabels(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetHostname(), req.GetAll())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListKubernetesNodeLabelsResponse{}
	resp.SetLabels(labels)
	return resp, nil
}

func (s *EssentialService) UpdateKubernetesNodeLabels(ctx context.Context, req *pb.UpdateKubernetesNodeLabelsRequest) (*pb.UpdateKubernetesNodeLabelsResponse, error) {
	labels, err := s.uc.UpdateKubernetesNodeLabels(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetHostname(), req.GetLabels())
	if err != nil {
		return nil, err
	}
	resp := &pb.UpdateKubernetesNodeLabelsResponse{}
	resp.SetLabels(labels)
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

func (s *EssentialService) ListGPURelationsByMachine(ctx context.Context, req *pb.ListGPURelationsByMachineRequest) (*pb.ListGPURelationsByMachineResponse, error) {
	gpuRelations, err := s.uc.ListGPURelationsByMachine(ctx,
		req.GetScopeUuid(),
		req.GetFacilityName(),
		req.GetMachineId(),
	)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListGPURelationsByMachineResponse{}
	resp.SetGpuRelations(convertGPURelationsToPB(gpuRelations))

	return resp, nil
}

func (s *EssentialService) ListGPURelationsByModel(ctx context.Context, req *pb.ListGPURelationsByModelRequest) (*pb.ListGPURelationsByModelResponse, error) {
	gpuRelations, err := s.uc.ListGPURelationsByModel(ctx,
		req.GetScopeUuid(),
		req.GetFacilityName(),
		req.GetNamespace(),
		req.GetModelName(),
	)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListGPURelationsByModelResponse{}
	resp.SetGpuRelations(convertGPURelationsToPB(gpuRelations))

	return resp, nil
}

// convertGPURelationsToPB converts domain GPURelation to protobuf GPURelation
func convertGPURelationsToPB(domainRelations []core.GPURelation) []*pb.GPURelation {
	pbRelations := make([]*pb.GPURelation, 0, len(domainRelations))

	for i := range domainRelations {
		pbRelation := &pb.GPURelation{}

		switch {
		case domainRelations[i].Machine != nil:
			machineEntity := &pb.GPURelation_Machine{}
			machineEntity.SetId(domainRelations[i].Machine.ID)
			machineEntity.SetHostname(domainRelations[i].Machine.Hostname)
			pbRelation.SetMachine(machineEntity)

		case domainRelations[i].GPU != nil:
			gpuEntity := &pb.GPURelation_GPU{}
			gpuEntity.SetId(domainRelations[i].GPU.ID)
			gpuEntity.SetVendor(domainRelations[i].GPU.Vendor)
			gpuEntity.SetProduct(domainRelations[i].GPU.Product)
			gpuEntity.SetMachineId(domainRelations[i].GPU.MachineID)

			// Convert vGPUs
			if len(domainRelations[i].GPU.VGPUs) > 0 {
				pbVGPUs := make([]*pb.GPURelation_GPUVGPU, 0, len(domainRelations[i].GPU.VGPUs))
				for j := range domainRelations[i].GPU.VGPUs {
					vgpu := &domainRelations[i].GPU.VGPUs[j]
					pbVGPU := &pb.GPURelation_GPUVGPU{}
					pbVGPU.SetPodName(vgpu.PodName)
					pbVGPU.SetBindingPhase(vgpu.BindingPhase)
					pbVGPU.SetVramBytes(vgpu.VramBytes)
					pbVGPU.SetVcoresPercent(vgpu.VcoresPercent)
					if !vgpu.BoundAt.IsZero() {
						pbVGPU.SetBoundAt(timestamppb.New(vgpu.BoundAt))
					}
					pbVGPUs = append(pbVGPUs, pbVGPU)
				}
				gpuEntity.SetVgpus(pbVGPUs)
			}
			pbRelation.SetGpu(gpuEntity)

		case domainRelations[i].Pod != nil:
			podEntity := &pb.GPURelation_Pod{}
			podEntity.SetName(domainRelations[i].Pod.Name)
			podEntity.SetNamespace(domainRelations[i].Pod.Namespace)
			podEntity.SetModelName(domainRelations[i].Pod.ModelName)
			podEntity.SetGpuIds(domainRelations[i].Pod.GPUIDs)
			pbRelation.SetPod(podEntity)
		}

		pbRelations = append(pbRelations, pbRelation)
	}

	return pbRelations
}
