package app

import (
	"context"

	corev1 "k8s.io/api/core/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/otterscale/otterscale/internal/core"

	pb "github.com/otterscale/otterscale/api/model/v1"
	pbconnect "github.com/otterscale/otterscale/api/model/v1/pbconnect"
)

type ModelService struct {
	pbconnect.UnimplementedModelServiceHandler

	model *core.ModelUseCase
}

func NewModelService(model *core.ModelUseCase) *ModelService {
	return &ModelService{
		model: model,
	}
}

var _ pbconnect.ModelServiceHandler = (*ModelService)(nil)

func (s *ModelService) CreateModelArtifact(ctx context.Context, req *pb.CreateModelArtifactRequest) (*pb.ModelArtifact, error){
	ma, err := s.model.CreateModelArtifact(ctx, req.GetScope(), req.GetFacility(), req.GetNamespace(), req.GetName(), req.GetModelname(), req.GetSize())
	if err != nil {
		return nil, err
	}
	resp := &pb.ModelArtifact{}
	resp.SetName(ma.Name)
	resp.SetModelname(ma.Modelname)
	resp.SetSize(ma.Size)
	return resp, nil
}

func (s *ModelService) CreateModelScheduler(ctx context.Context, req *pb.CreateModelSchedulerRequest) (*pb.ModelScheduler, error) {
	var brs []core.BackendRef
	for _, r := range req.GetBackendRefs() {
    	brs = append(brs, core.BackendRef{
        	Name:   r.GetName(),
        	Weight: r.GetWeight(),
    	})
	}
	decEnvMap := make(map[string]string, len(req.GetDecodeEnvs()))
	for _, e := range req.GetDecodeEnvs() {
		if n := e.GetName(); n != "" {
			decEnvMap[n] = e.GetValue()
		}
	}
    preEnvMap := make(map[string]string, len(req.GetPrefillEnvs()))
	for _, e := range req.GetPrefillEnvs() {
		if n := e.GetName(); n != "" {
			preEnvMap[n] = e.GetValue()
		}
	}
	ms, err := s.model.CreateModelScheduler(ctx, req.GetScope(), req.GetFacility(), req.GetNamespace(), req.GetName(), req.GetModelArtifactsName(), 
		req.GetUri(), req.GetMultinode(), req.GetNvidia(), req.GetParentRefsName(), req.GetHttpRouteCreate(), brs, req.GetBackendRequest(), 
		req.GetRequest(), req.GetEppCreate(), req.GetDecodeCreate(), req.GetDecodeReplicas(), req.GetDecodeArgs(), decEnvMap, req.GetDecodeResourcesLimitsGpu(),
		req.GetDecodeResourcesLimitsGpumem(), req.GetPrefillCreate(), req.GetPrefillReplicas(), req.GetPrefillArgs(), preEnvMap, req.GetPrefillResourcesLimitsGpu(), 
		req.GetPrefillResourcesLimitsGpumem(), req.GetInferenceExtensionReplicas(), req.GetExtProcPort(), req.GetPluginsConfigFile(), req.GetTargetPortNumber())
	if err != nil {
		return nil, err
	}
	resp := &pb.ModelScheduler{}
	resp.SetName(ms.Name)
	resp.SetModelArtifactsName(ms.Modelartifactsname)
	resp.SetHttpRouteCreate(ms.HttpRouteCreate)
	resp.SetEppCreate(ms.EppCreate)
	resp.SetPrefillCreate(ms.PrefillCreate)
	if len(ms.BackendRefs) > 0 {
		out := make([]*pb.ModelScheduler_BackendRef, 0, len(ms.BackendRefs))
		for _, br := range ms.BackendRefs {
	    	r := &pb.ModelScheduler_BackendRef{}
    		r.SetName(br.Name)
	  		r.SetWeight(br.Weight)
   			out = append(out, r)
		}
	}
	return resp, nil
}

func (s *ModelService) CreateModelGateway(ctx context.Context, req *pb.CreateModelGatewayRequest) (*pb.ModelGateway, error) {
	mg, err := s.model.CreateModelGateway(ctx, req.GetScope(), req.GetFacility(), req.GetNamespace(), req.GetName(), req.GetCpu(), req.GetMemory(), req.GetType())
	if err != nil {
		return nil, err
	}
	resp := &pb.ModelGateway{}
	resp.SetName(mg.Name)
	resp.SetPublicAddress(mg.Publicaddress)
	return resp, nil
}

func (s *ModelService) ListModelArtifacts(ctx context.Context, req *pb.ListModelArtifactsRequest) (*pb.ListModelArtifactsResponse, error) {
	pvcs, err := s.model.ListModelArtifactPVCs(ctx, req.GetScope(), req.GetFacility(), req.GetNamespace())
	if err != nil {
		return nil, err
	}
	modelartifacts := make([]*pb.ListModelArtifactsResponse_PersistentVolumeClaim, 0, len(pvcs))
	for i := range pvcs {
		modelartifacts = append(modelartifacts, toProtoModelArtifactPVC(&pvcs[i]))
	}
	resp := &pb.ListModelArtifactsResponse{}
	resp.SetPersistentVolumeClaims(modelartifacts)
	return resp, nil
}

func toProtoModelArtifactPVC(pvc *corev1.PersistentVolumeClaim) *pb.ListModelArtifactsResponse_PersistentVolumeClaim {
	modelartifact := &pb.ListModelArtifactsResponse_PersistentVolumeClaim{}
	modelartifact.SetName(pvc.Name)
	modelartifact.SetStatus(string(pvc.Status.Phase))
	modelartifact.SetAccessModes(accessModesToStrings(pvc.Spec.AccessModes))
	modelartifact.SetCapacity(pvc.Spec.Resources.Requests.Storage().String())
	modelartifact.SetCreatedAt(timestamppb.New(pvc.CreationTimestamp.Time))
	return modelartifact
}