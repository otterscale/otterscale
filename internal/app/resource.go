package app

import (
	"context"
	"encoding/json"
	"fmt"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/kube-openapi/pkg/validation/spec"

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

func (s *ResourceService) Discovery(_ context.Context, req *pb.DiscoveryRequest) (*pb.DiscoveryResponse, error) {
	apiResources, err := s.resource.ListAPIResources(req.GetCluster())
	if err != nil {
		return nil, err
	}

	pbAPIResources, err := s.toProtoAPIResources(apiResources)
	if err != nil {
		return nil, err
	}

	resp := &pb.DiscoveryResponse{}
	resp.SetApiResources(pbAPIResources)
	return resp, nil
}

func (s *ResourceService) Schema(_ context.Context, req *pb.SchemaRequest) (*pb.SchemaResponse, error) {
	schema, err := s.resource.GetSchema(req.GetCluster(), req.GetGroup(), req.GetVersion(), req.GetKind())
	if err != nil {
		return nil, err
	}

	jsonSchema, err := s.toProtoStructFromJSONSchema(schema.JSONSchema)
	if err != nil {
		return nil, err
	}

	resp := &pb.SchemaResponse{}
	resp.SetJsonSchema(jsonSchema)
	// todo: ui schema
	// resp.SetUiSchema(uiSchema)
	return resp, nil
}

func (s *ResourceService) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	cgvr, err := s.resource.Validate(req.GetCluster(), req.GetGroup(), req.GetVersion(), req.GetResource())
	if err != nil {
		return nil, err
	}

	resources, err := s.resource.ListResources(ctx, cgvr, req.GetNamespace(), req.GetLabelSelector(), req.GetFieldSelector(), req.GetLimit(), req.GetContinue())
	if err != nil {
		return nil, err
	}

	pbResources, err := s.toProtoResources(resources.Items)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListResponse{}
	resp.SetResourceVersion(resources.GetResourceVersion())
	resp.SetContinue(resources.GetContinue())
	resp.SetRemainingItemCount(deref(resources.GetRemainingItemCount(), 0))
	resp.SetItems(pbResources)
	return resp, nil
}

func (s *ResourceService) Get(ctx context.Context, req *pb.GetRequest) (*pb.Resource, error) {
	cgvr, err := s.resource.Validate(req.GetCluster(), req.GetGroup(), req.GetVersion(), req.GetResource())
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
	cgvr, err := s.resource.Validate(req.GetCluster(), req.GetGroup(), req.GetVersion(), req.GetResource())
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
	cgvr, err := s.resource.Validate(req.GetCluster(), req.GetGroup(), req.GetVersion(), req.GetResource())
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
	cgvr, err := s.resource.Validate(req.GetCluster(), req.GetGroup(), req.GetVersion(), req.GetResource())
	if err != nil {
		return nil, err
	}

	if err := s.resource.DeleteResource(ctx, cgvr, req.GetNamespace(), req.GetName(), req.GetGracePeriodSeconds()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *ResourceService) Watch(ctx context.Context, req *pb.WatchRequest, stream *connect.ServerStream[pb.WatchEvent]) error {
	cgvr, err := s.resource.Validate(req.GetCluster(), req.GetGroup(), req.GetVersion(), req.GetResource())
	if err != nil {
		return err
	}

	watcher, err := s.resource.WatchResource(ctx, cgvr, req.GetNamespace(), req.GetLabelSelector(), req.GetFieldSelector(), req.GetResourceVersion())
	if err != nil {
		return err
	}
	defer watcher.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case event, ok := <-watcher.ResultChan():
			if !ok {
				return connect.NewError(connect.CodeUnavailable, fmt.Errorf("watch closed"))
			}

			msg, ok := s.processEvent(event)
			if !ok {
				continue
			}

			if err := stream.Send(msg); err != nil {
				return err
			}
		}
	}
}

func (s *ResourceService) processEvent(event watch.Event) (*pb.WatchEvent, bool) {
	switch event.Type {
	case watch.Added, watch.Modified, watch.Deleted:
		obj, ok := event.Object.(*unstructured.Unstructured)
		if !ok {
			return nil, false
		}

		resource, err := s.toProtoResource(obj.Object)
		if err != nil {
			return nil, false
		}

		ret := &pb.WatchEvent{}
		ret.SetType(s.toProtoWatchEventType(event.Type))
		ret.SetResource(resource)
		return ret, true

	case watch.Bookmark:
		metadata, _ := meta.Accessor(event.Object)

		ret := &pb.WatchEvent{}
		ret.SetType(pb.WatchEvent_TYPE_BOOKMARK)
		ret.SetResourceVersion(metadata.GetResourceVersion())
		return ret, true

	case watch.Error:
		ret := &pb.WatchEvent{}
		ret.SetType(pb.WatchEvent_TYPE_ERROR)
		return ret, true

	default:
		return nil, false
	}
}

func (s *ResourceService) toProtoAPIResources(list []*metav1.APIResourceList) ([]*pb.APIResource, error) {
	ret := []*pb.APIResource{}

	for i := range list {
		gv, err := schema.ParseGroupVersion(list[i].GroupVersion)
		if err != nil {
			return nil, err
		}

		for j := range list[i].APIResources {
			ret = append(ret, s.toProtoAPIResource(gv, &list[i].APIResources[j]))
		}
	}

	return ret, nil
}

func (s *ResourceService) toProtoAPIResource(gv schema.GroupVersion, r *metav1.APIResource) *pb.APIResource {
	ret := &pb.APIResource{}
	ret.SetGroup(gv.Group)
	ret.SetVersion(gv.Version)
	ret.SetResource(r.Name)
	ret.SetKind(r.Kind)
	ret.SetNamespaced(r.Namespaced)
	ret.SetVerbs(r.Verbs)
	ret.SetShortNames(r.ShortNames)
	return ret
}

func (s *ResourceService) toProtoStructFromJSONSchema(js *spec.Schema) (*structpb.Struct, error) {
	jsBytes, err := json.Marshal(js)
	if err != nil {
		return nil, err
	}

	ret := &structpb.Struct{}
	if err := protojson.Unmarshal(jsBytes, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (s *ResourceService) toProtoResources(list []unstructured.Unstructured) ([]*pb.Resource, error) {
	ret := []*pb.Resource{}

	for i := range list {
		res, err := s.toProtoResource(list[i].Object)
		if err != nil {
			return nil, err
		}

		ret = append(ret, res)
	}

	return ret, nil
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

func (s *ResourceService) toProtoWatchEventType(t watch.EventType) pb.WatchEvent_Type {
	switch t {
	case watch.Added:
		return pb.WatchEvent_TYPE_ADDED

	case watch.Modified:
		return pb.WatchEvent_TYPE_MODIFIED

	case watch.Deleted:
		return pb.WatchEvent_TYPE_DELETED

	case watch.Bookmark:
		return pb.WatchEvent_TYPE_BOOKMARK

	case watch.Error:
		return pb.WatchEvent_TYPE_ERROR

	default:
		return pb.WatchEvent_TYPE_UNSPECIFIED
	}
}

func deref[T any](ptr *T, def T) T {
	if ptr != nil {
		return *ptr
	}
	return def
}
