package app

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pb "github.com/otterscale/otterscale/api/resource/v1"
	"github.com/otterscale/otterscale/api/resource/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core/resource"
)

type ResourceService struct {
	pbconnect.UnimplementedResourceServiceHandler

	resource *resource.UseCase
}

func NewResourceService(resource *resource.UseCase) *ResourceService {
	return &ResourceService{
		resource: resource,
	}
}

var _ pbconnect.ResourceServiceHandler = (*ResourceService)(nil)

func (s *ResourceService) Discovery(ctx context.Context, req *pb.DiscoveryRequest) (*pb.DiscoveryResponse, error) {
	apiResources, err := s.resource.ListAPIResources(ctx, req.GetCluster())
	if err != nil {
		return nil, err
	}

	resp := &pb.DiscoveryResponse{}
	resp.SetApiResources(s.toProtoAPIResourcesFromList(apiResources))
	return resp, nil
}

func (s *ResourceService) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	cgvr, err := s.resource.Validate(ctx, req.GetCluster(), req.GetGroup(), req.GetVersion(), req.GetResource())
	if err != nil {
		return nil, err
	}

	resources, err := s.resource.ListResources(ctx, cgvr, req.GetNamespace(), req.GetLabelSelector(), req.GetFieldSelector(), req.GetLimit(), req.GetContinue())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListResponse{}
	// resp.SetResourceVersion("")
	// resp.SetContinue("")
	// resp.SetRemainingItemCount(0)
	// resp.SetItems(s.toProtoResources(resources))
	return resp, nil
}

func (s *ResourceService) Get(ctx context.Context, req *pb.GetRequest) (*pb.Resource, error) {
	cgvr, err := s.resource.Validate(ctx, req.GetCluster(), req.GetGroup(), req.GetVersion(), req.GetResource())
	if err != nil {
		return nil, err
	}

	resource, err := s.resource.GetResource(ctx, cgvr, req.GetNamespace(), req.GetName())
	if err != nil {
		return nil, err
	}

	return s.toProtoResource(resource.Object)
}

func (s *ResourceService) Create(ctx context.Context, req *pb.CreateRequest) (*pb.Resource, error) {
	cgvr, err := s.resource.Validate(ctx, req.GetCluster(), req.GetGroup(), req.GetVersion(), req.GetResource())
	if err != nil {
		return nil, err
	}

	resource, err := s.resource.CreateResource(ctx, cgvr, req.GetNamespace(), req.GetManifest())
	if err != nil {
		return nil, err
	}

	return s.toProtoResource(resource.Object)
}

func (s *ResourceService) Apply(ctx context.Context, req *pb.ApplyRequest) (*pb.Resource, error) {
	cgvr, err := s.resource.Validate(ctx, req.GetCluster(), req.GetGroup(), req.GetVersion(), req.GetResource())
	if err != nil {
		return nil, err
	}

	resource, err := s.resource.ApplyResource(ctx, cgvr, req.GetNamespace(), req.GetName(), req.GetManifest(), req.GetForce(), req.GetFieldManager())
	if err != nil {
		return nil, err
	}

	return s.toProtoResource(resource.Object)
}

func (s *ResourceService) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	cgvr, err := s.resource.Validate(ctx, req.GetCluster(), req.GetGroup(), req.GetVersion(), req.GetResource())
	if err != nil {
		return nil, err
	}

	if err := s.resource.DeleteResource(ctx, cgvr, req.GetNamespace(), req.GetName(), req.GetGracePeriodSeconds()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// func (s *ResourceService) Watch(ctx context.Context, req *pb.WatchRequest) (*pb.WatchResponse, error) {
// }

func (s *ResourceService) toProtoAPIResourcesFromList(list []*metav1.APIResourceList) []*pb.APIResource {
	ret := []*pb.APIResource{}

	for i := range list {
		ret = append(ret, s.toProtoAPIResources(list[i].APIResources)...)
	}

	return ret
}

func (s *ResourceService) toProtoAPIResources(rs []metav1.APIResource) []*pb.APIResource {
	ret := []*pb.APIResource{}

	for i := range rs {
		ret = append(ret, s.toProtoAPIResource(rs[i]))
	}

	return ret
}

func (s *ResourceService) toProtoAPIResource(r metav1.APIResource) *pb.APIResource {
	ret := &pb.APIResource{}
	ret.SetGroup(r.Group)
	ret.SetVersion(r.Version)
	ret.SetResource(r.Name)
	ret.SetKind(r.Kind)
	ret.SetNamespaced(r.Namespaced)
	ret.SetVerbs(r.Verbs)
	ret.SetShortNames(r.ShortNames)
	return ret
}

func (s *ResourceService) toProtoResource(obj map[string]any) (*pb.Resource, error) {
	object, err := structpb.NewStruct(obj)
	if err != nil {
		return nil, err
	}

	ret := &pb.Resource{}
	ret.SetObject(object)
	return ret, nil
}
