package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kube-openapi/pkg/validation/spec"

	pb "github.com/otterscale/api/resource/v1"
	"github.com/otterscale/otterscale/internal/core"
)

// ResourceService implements the Resource gRPC service. It proxies
// Kubernetes CRUD and watch operations through the tunnel, translating
// between protobuf and unstructured Kubernetes objects.
type ResourceService struct {
	pb.UnimplementedResourceServiceHandler

	resource *core.ResourceUseCase
}

// NewResourceService returns a ResourceService backed by the given
// use-case.
func NewResourceService(resource *core.ResourceUseCase) *ResourceService {
	return &ResourceService{
		resource: resource,
	}
}

var _ pb.ResourceServiceHandler = (*ResourceService)(nil)

// ---------------------------------------------------------------------------
// Discovery / Schema
// ---------------------------------------------------------------------------

// Discovery returns the full list of API resources available on the
// target cluster.
func (s *ResourceService) Discovery(ctx context.Context, req *pb.DiscoveryRequest) (*pb.DiscoveryResponse, error) {
	apiResources, err := s.resource.ServerResources(ctx, req.GetCluster())
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}

	pbAPIResources, err := toProtoAPIResources(apiResources)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := &pb.DiscoveryResponse{}
	resp.SetApiResources(pbAPIResources)
	return resp, nil
}

// Schema returns the OpenAPI schema for the given GVK, serialized as
// a protobuf Struct.
func (s *ResourceService) Schema(ctx context.Context, req *pb.SchemaRequest) (*structpb.Struct, error) {
	resolved, err := s.resource.ResolveSchema(
		ctx,
		req.GetCluster(),
		req.GetGroup(),
		req.GetVersion(),
		req.GetKind(),
	)
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}
	result, err := toProtoStructFromJSONSchema(resolved)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// CRUD
// ---------------------------------------------------------------------------

// List returns a paged list of resources matching the request filters.
func (s *ResourceService) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	resources, err := s.resource.ListResources(
		ctx,
		&core.ResourceIdentifier{
			Cluster:   req.GetCluster(),
			Group:     req.GetGroup(),
			Version:   req.GetVersion(),
			Resource:  req.GetResource(),
			Namespace: req.GetNamespace(),
		},
		core.ListOptions{
			LabelSelector: req.GetLabelSelector(),
			FieldSelector: req.GetFieldSelector(),
			Limit:         req.GetLimit(),
			Continue:      req.GetContinue(),
		},
	)
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}

	// Strip noisy metadata (managedFields, last-applied-configuration)
	// before serializing to protobuf. This is a presentation concern
	// that belongs in the handler layer, not the domain use-case.
	for i := range resources.Items {
		cleanObject(resources.Items[i].Object)
	}

	pbResources, err := toProtoResources(resources.Items)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := &pb.ListResponse{}
	resp.SetResourceVersion(resources.GetResourceVersion())
	resp.SetContinue(resources.GetContinue())
	resp.SetRemainingItemCount(deref(resources.GetRemainingItemCount(), 0))
	resp.SetItems(pbResources)
	return resp, nil
}

// Get returns a single resource by name.
func (s *ResourceService) Get(ctx context.Context, req *pb.GetRequest) (*pb.Resource, error) {
	resource, err := s.resource.GetResource(
		ctx,
		&core.ResourceIdentifier{
			Cluster:   req.GetCluster(),
			Group:     req.GetGroup(),
			Version:   req.GetVersion(),
			Resource:  req.GetResource(),
			Namespace: req.GetNamespace(),
			Name:      req.GetName(),
		},
	)
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}
	result, err := toProtoResource(resource.Object)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return result, nil
}

// Create creates a new resource from the YAML manifest in the request.
func (s *ResourceService) Create(ctx context.Context, req *pb.CreateRequest) (*pb.Resource, error) {
	resource, err := s.resource.CreateResource(
		ctx,
		&core.ResourceIdentifier{
			Cluster:   req.GetCluster(),
			Group:     req.GetGroup(),
			Version:   req.GetVersion(),
			Resource:  req.GetResource(),
			Namespace: req.GetNamespace(),
		},
		req.GetManifest(),
	)
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}
	result, err := toProtoResource(resource.Object)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return result, nil
}

// Apply performs a server-side apply for the given resource.
func (s *ResourceService) Apply(ctx context.Context, req *pb.ApplyRequest) (*pb.Resource, error) {
	resource, err := s.resource.ApplyResource(
		ctx,
		&core.ResourceIdentifier{
			Cluster:   req.GetCluster(),
			Group:     req.GetGroup(),
			Version:   req.GetVersion(),
			Resource:  req.GetResource(),
			Namespace: req.GetNamespace(),
			Name:      req.GetName(),
		},
		req.GetManifest(),
		core.ApplyOptions{
			Force:        req.GetForce(),
			FieldManager: req.GetFieldManager(),
		},
	)
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}
	result, err := toProtoResource(resource.Object)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return result, nil
}

// Delete removes the named resource. An optional grace period may be
// specified in the request.
func (s *ResourceService) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	var opts core.DeleteOptions
	if req.HasGracePeriodSeconds() {
		v := req.GetGracePeriodSeconds()
		opts.GracePeriodSeconds = &v
	}

	if err := s.resource.DeleteResource(
		ctx,
		&core.ResourceIdentifier{
			Cluster:   req.GetCluster(),
			Group:     req.GetGroup(),
			Version:   req.GetVersion(),
			Resource:  req.GetResource(),
			Namespace: req.GetNamespace(),
			Name:      req.GetName(),
		},
		opts,
	); err != nil {
		return nil, domainErrorToConnectError(err)
	}
	return &emptypb.Empty{}, nil
}

// ---------------------------------------------------------------------------
// Describe
// ---------------------------------------------------------------------------

// Describe returns a resource together with its related Kubernetes
// events, equivalent to `kubectl describe`.
func (s *ResourceService) Describe(ctx context.Context, req *pb.DescribeRequest) (*pb.DescribeResponse, error) {
	obj, events, err := s.resource.DescribeResource(
		ctx,
		&core.ResourceIdentifier{
			Cluster:   req.GetCluster(),
			Group:     req.GetGroup(),
			Version:   req.GetVersion(),
			Resource:  req.GetResource(),
			Namespace: req.GetNamespace(),
			Name:      req.GetName(),
		},
	)
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}

	pbResource, err := toProtoResource(obj.Object)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbEvents, err := toProtoResources(events.Items)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := &pb.DescribeResponse{}
	resp.SetResource(pbResource)
	resp.SetEvents(pbEvents)
	return resp, nil
}

// ---------------------------------------------------------------------------
// Watch
// ---------------------------------------------------------------------------

// Watch opens a server-streaming RPC that forwards Kubernetes watch
// events to the client. The stream ends when the client cancels the
// context or the upstream watcher closes.
func (s *ResourceService) Watch(ctx context.Context, req *pb.WatchRequest, stream *connect.ServerStream[pb.WatchEvent]) error {
	watcher, err := s.resource.WatchResource(
		ctx,
		&core.ResourceIdentifier{
			Cluster:   req.GetCluster(),
			Group:     req.GetGroup(),
			Version:   req.GetVersion(),
			Resource:  req.GetResource(),
			Namespace: req.GetNamespace(),
		},
		core.WatchOptions{
			LabelSelector:   req.GetLabelSelector(),
			FieldSelector:   req.GetFieldSelector(),
			ResourceVersion: req.GetResourceVersion(),
		},
	)
	if err != nil {
		return domainErrorToConnectError(err)
	}
	defer watcher.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case event, ok := <-watcher.ResultChan():
			if !ok {
				return connect.NewError(connect.CodeUnavailable, errors.New("watch closed"))
			}

			msg, err := processEvent(event)
			if err != nil {
				slog.Warn("watch: skipping event", "error", err)
				continue
			}

			if err := stream.Send(msg); err != nil {
				return err
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// processEvent converts a domain core.WatchEvent into a protobuf
// WatchEvent. Returns an error if the event should be skipped; the
// caller is responsible for logging. This keeps the function free of
// side effects so it can be unit-tested without producing log output.
func processEvent(event core.WatchEvent) (*pb.WatchEvent, error) {
	switch event.Type {
	case core.WatchEventAdded, core.WatchEventModified, core.WatchEventDeleted:
		if event.Object == nil {
			return nil, fmt.Errorf("nil object for event type %s", event.Type)
		}

		resource, err := toProtoResource(event.Object)
		if err != nil {
			return nil, fmt.Errorf("convert resource for event type %s: %w", event.Type, err)
		}

		ret := &pb.WatchEvent{}
		ret.SetType(toProtoWatchEventType(event.Type))
		ret.SetResource(resource)
		return ret, nil

	case core.WatchEventBookmark:
		ret := &pb.WatchEvent{}
		ret.SetType(pb.WatchEvent_TYPE_BOOKMARK)
		// Extract resourceVersion from the bookmark object.
		if event.Object != nil {
			if metadata, ok := event.Object["metadata"].(map[string]any); ok {
				if rv, ok := metadata["resourceVersion"].(string); ok {
					ret.SetResourceVersion(rv)
				}
			}
		}
		return ret, nil

	case core.WatchEventError:
		ret := &pb.WatchEvent{}
		ret.SetType(pb.WatchEvent_TYPE_ERROR)
		if event.Object != nil {
			if resource, err := toProtoResource(event.Object); err == nil {
				ret.SetResource(resource)
			}
		}
		return ret, nil

	default:
		return nil, fmt.Errorf("unknown event type: %s", event.Type)
	}
}

// toProtoAPIResources flattens the Kubernetes APIResourceList slice
// into a single []*pb.APIResource list, embedding the parsed
// group/version into each entry.
func toProtoAPIResources(list []*metav1.APIResourceList) ([]*pb.APIResource, error) {
	// Estimate total capacity to avoid repeated allocations.
	total := 0
	for i := range list {
		total += len(list[i].APIResources)
	}
	ret := make([]*pb.APIResource, 0, total)

	for i := range list {
		gv, err := schema.ParseGroupVersion(list[i].GroupVersion)
		if err != nil {
			return nil, err
		}

		for j := range list[i].APIResources {
			ret = append(ret, toProtoAPIResource(gv, &list[i].APIResources[j]))
		}
	}

	return ret, nil
}

// toProtoAPIResource converts a single Kubernetes APIResource into its
// protobuf representation.
func toProtoAPIResource(gv schema.GroupVersion, r *metav1.APIResource) *pb.APIResource {
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

// toProtoStructFromJSONSchema serializes an OpenAPI spec.Schema to
// JSON and re-parses it into a protobuf Struct so it can be returned
// as a generic structured response.
func toProtoStructFromJSONSchema(js *spec.Schema) (*structpb.Struct, error) {
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

// toProtoResources converts a slice of Unstructured objects into
// protobuf Resource messages.
func toProtoResources(list []unstructured.Unstructured) ([]*pb.Resource, error) {
	ret := make([]*pb.Resource, 0, len(list))

	for i := range list {
		res, err := toProtoResource(list[i].Object)
		if err != nil {
			return nil, err
		}

		ret = append(ret, res)
	}

	return ret, nil
}

// toProtoResource wraps a raw Kubernetes object map in a protobuf
// Resource message.
func toProtoResource(obj map[string]any) (*pb.Resource, error) {
	object, err := structpb.NewStruct(obj)
	if err != nil {
		return nil, err
	}

	ret := &pb.Resource{}
	ret.SetObject(object)
	return ret, nil
}

// toProtoWatchEventType maps a domain WatchEventType to the protobuf
// WatchEvent_Type enum.
func toProtoWatchEventType(t core.WatchEventType) pb.WatchEvent_Type {
	switch t {
	case core.WatchEventAdded:
		return pb.WatchEvent_TYPE_ADDED
	case core.WatchEventModified:
		return pb.WatchEvent_TYPE_MODIFIED
	case core.WatchEventDeleted:
		return pb.WatchEvent_TYPE_DELETED
	case core.WatchEventBookmark:
		return pb.WatchEvent_TYPE_BOOKMARK
	case core.WatchEventError:
		return pb.WatchEvent_TYPE_ERROR
	default:
		return pb.WatchEvent_TYPE_UNSPECIFIED
	}
}
