package app

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/otterscale/otterscale/api/model/v1"
	pbconnect "github.com/otterscale/otterscale/api/model/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core/model"
)

type ModelService struct {
	pbconnect.UnimplementedModelServiceHandler

	model *model.UseCase
}

func NewModelService(model *model.UseCase) *ModelService {
	return &ModelService{
		model: model,
	}
}

var _ pbconnect.ModelServiceHandler = (*ModelService)(nil)

func (s *ModelService) ListModels(ctx context.Context, req *pb.ListModelsRequest) (*pb.ListModelsResponse, error) {
	models, uri, err := s.model.ListModels(ctx, req.GetScope(), req.GetNamespace())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListModelsResponse{}
	resp.SetModels(toProtoModels(models))
	resp.SetServiceUri(uri)
	return resp, nil
}

func (s *ModelService) CreateModel(ctx context.Context, req *pb.CreateModelRequest) (*pb.Model, error) {
	var prefill, decode *model.Resource

	if r := req.GetPrefill(); r != nil {
		prefill = toModelResource(r)
	}

	if r := req.GetDecode(); r != nil {
		decode = toModelResource(r)
	}

	model, err := s.model.CreateModel(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetModelName(), req.GetSizeBytes(), prefill, decode)
	if err != nil {
		return nil, err
	}

	resp := toProtoModel(model)
	return resp, nil
}

func (s *ModelService) UpdateModel(ctx context.Context, req *pb.UpdateModelRequest) (*pb.Model, error) {
	var prefill, decode *model.Resource

	if r := req.GetPrefill(); r != nil {
		prefill = toModelResource(r)
	}

	if r := req.GetDecode(); r != nil {
		decode = toModelResource(r)
	}

	model, err := s.model.UpdateModel(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), prefill, decode)
	if err != nil {
		return nil, err
	}

	resp := toProtoModel(model)
	return resp, nil
}

func (s *ModelService) DeleteModel(ctx context.Context, req *pb.DeleteModelRequest) (*emptypb.Empty, error) {
	if err := s.model.DeleteModel(ctx, req.GetScope(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ModelService) ListModelArtifacts(ctx context.Context, req *pb.ListModelArtifactsRequest) (*pb.ListModelArtifactsResponse, error) {
	artifacts, err := s.model.ListModelArtifacts(ctx, req.GetScope(), req.GetNamespace())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListModelArtifactsResponse{}
	resp.SetModelArtifacts(toProtoModelArtifacts(artifacts))
	return resp, nil
}

func (s *ModelService) CreateModelArtifact(ctx context.Context, req *pb.CreateModelArtifactRequest) (*pb.ModelArtifact, error) {
	artifact, err := s.model.CreateModelArtifact(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetModelName(), req.GetSize())
	if err != nil {
		return nil, err
	}

	resp := toProtoModelArtifact(artifact)
	return resp, nil
}

func (s *ModelService) DeleteModelArtifact(ctx context.Context, req *pb.DeleteModelArtifactRequest) (*emptypb.Empty, error) {
	if err := s.model.DeleteModelArtifact(ctx, req.GetScope(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func toModelResource(r *pb.Model_Resource) *model.Resource {
	return &model.Resource{
		VGPU:       r.GetVgpu(),
		VGPUMemory: r.GetVgpumemPercentage(),
	}
}

func toProtoModelResource(r *model.Resource) *pb.Model_Resource {
	ret := &pb.Model_Resource{}
	ret.SetVgpu(r.VGPU)
	ret.SetVgpumemPercentage(r.VGPUMemory)
	return ret
}

func toProtoModels(ms []model.Model) []*pb.Model {
	ret := []*pb.Model{}

	for i := range ms {
		ret = append(ret, toProtoModel(&ms[i]))
	}

	return ret
}

func toProtoModel(m *model.Model) *pb.Model {
	ret := &pb.Model{}
	ret.SetId(m.ID)

	release := m.Release
	if release != nil {
		ret.SetName(release.Name)
		ret.SetNamespace(release.Namespace)

		info := release.Info
		if info != nil {
			ret.SetStatus(string(info.Status))
			ret.SetDescription(info.Description)
			ret.SetFirstDeployedAt(timestamppb.New(info.FirstDeployed.Time))
			ret.SetLastDeployedAt(timestamppb.New(info.LastDeployed.Time))
		}

		chart := release.Chart
		if chart != nil && chart.Metadata != nil {
			ret.SetChartVersion(chart.Metadata.Version)
			ret.SetAppVersion(chart.Metadata.AppVersion)
		}
	}

	prefill := m.Prefill
	if prefill != nil {
		ret.SetPrefill(toProtoModelResource(prefill))
	}

	decode := m.Decode
	if decode != nil {
		ret.SetDecode(toProtoModelResource(decode))
	}

	ret.SetPods(toProtoPods(m.Pods))

	return ret
}

func toProtoModelArtifacts(mas []model.Artifact) []*pb.ModelArtifact {
	ret := []*pb.ModelArtifact{}

	for i := range mas {
		ret = append(ret, toProtoModelArtifact(&mas[i]))
	}

	return ret
}

func toProtoModelArtifact(ma *model.Artifact) *pb.ModelArtifact {
	ret := &pb.ModelArtifact{}
	ret.SetName(ma.Name)
	ret.SetNamespace(ma.Namespace)
	ret.SetModelName(ma.Modelname)
	ret.SetPhase(string(ma.Phase))
	ret.SetSize(ma.Size)
	ret.SetVolumeName(ma.VolumeName)
	ret.SetCreatedAt(timestamppb.New(ma.CreatedAt))
	return ret
}
