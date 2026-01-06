package openfeature

import (
	"context"
	"fmt"
	"sync"

	"connectrpc.com/connect"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/otterscale/otterscale/api"
)

const domain = "otterscale"

type Interceptor struct {
	client     *openfeature.Client
	featureMap *sync.Map
}

func NewInterceptor() (*Interceptor, error) {
	provider, err := flagd.NewProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to create flagd provider: %w", err)
	}
	if err := openfeature.SetProviderAndWait(provider); err != nil {
		return nil, fmt.Errorf("failed to set openfeature provider: %w", err)
	}

	client := openfeature.NewClient(domain)
	featureMap := &sync.Map{}

	return &Interceptor{
		client:     client,
		featureMap: featureMap,
	}, nil
}

func (i *Interceptor) Evaluate(ctx context.Context, flagName string) (bool, error) {
	return i.client.BooleanValue(ctx, flagName, false, openfeature.EvaluationContext{})
}

func (i *Interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		featureName, err := i.resolveFeatureName(req.Spec())
		if err != nil {
			return nil, err
		}

		if featureName == "" {
			return next(ctx, req)
		}

		enabled, err := i.Evaluate(ctx, featureName)
		if err != nil {
			return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("failed to evaluate feature %q: %w", featureName, err))
		}

		if !enabled {
			return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("feature %q is disabled", featureName))
		}

		return next(ctx, req)
	}
}

func (i *Interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i *Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		featureName, err := i.resolveFeatureName(conn.Spec())
		if err != nil {
			return err
		}

		if featureName == "" {
			return next(ctx, conn)
		}

		enabled, err := i.Evaluate(ctx, featureName)
		if err != nil {
			return connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("failed to evaluate feature %q: %w", featureName, err))
		}

		if !enabled {
			return connect.NewError(connect.CodeUnimplemented, fmt.Errorf("feature %q is disabled", featureName))
		}

		return next(ctx, conn)
	}
}

func (i *Interceptor) resolveFeatureName(spec connect.Spec) (string, error) {
	if v, ok := i.featureMap.Load(spec.Procedure); ok {
		return v.(string), nil
	}

	featureName, err := extractFeatureFromDescriptor(spec)
	if err != nil {
		return "", err
	}

	i.featureMap.Store(spec.Procedure, featureName)
	return featureName, nil
}

func extractFeatureFromDescriptor(spec connect.Spec) (string, error) {
	desc, ok := spec.Schema.(protoreflect.MethodDescriptor)
	if !ok {
		return "", connect.NewError(connect.CodeInternal, fmt.Errorf("unable to get method descriptor"))
	}

	opts, ok := desc.Options().(*descriptorpb.MethodOptions)
	if !ok {
		return "", connect.NewError(connect.CodeInternal, fmt.Errorf("unable to get method options"))
	}

	if !proto.HasExtension(opts, api.E_Feature) {
		return "", nil
	}

	feature, ok := proto.GetExtension(opts, api.E_Feature).(*api.Feature)
	if !ok {
		return "", connect.NewError(connect.CodeInternal, fmt.Errorf("unable to get feature extension"))
	}

	return feature.GetName(), nil
}
