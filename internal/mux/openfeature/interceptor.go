package openfeature

import (
	"context"
	"fmt"
	"sync"

	"connectrpc.com/connect"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/otterscale/otterscale/api"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

const domain = "otterscale"

type Interceptor struct {
	connect.UnaryInterceptorFunc
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

	interceptor := &Interceptor{
		client:     client,
		featureMap: featureMap,
	}

	interceptor.UnaryInterceptorFunc = connect.UnaryInterceptorFunc(interceptor.intercept)
	return interceptor, nil
}

func (i *Interceptor) intercept(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		featureName, err := i.getFeatureName(req)
		if err != nil {
			return nil, err
		}

		enabled, err := i.client.BooleanValue(ctx, featureName, false, openfeature.EvaluationContext{})
		if err != nil {
			return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("failed to evaluate feature %q: %w", featureName, err))
		}

		if !enabled {
			return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("feature %q is disabled", featureName))
		}

		return next(ctx, req)
	}
}

func (i *Interceptor) getFeatureName(req connect.AnyRequest) (string, error) {
	procedure := req.Spec().Procedure

	if v, ok := i.featureMap.Load(procedure); ok {
		return v.(string), nil
	}

	featureName, err := extractFeatureFromDescriptor(req)
	if err != nil {
		return "", err
	}

	i.featureMap.Store(procedure, featureName)
	return featureName, nil
}

func extractFeatureFromDescriptor(req connect.AnyRequest) (string, error) {
	desc, ok := req.Spec().Schema.(protoreflect.MethodDescriptor)
	if !ok {
		return "", connect.NewError(connect.CodeInternal, fmt.Errorf("unable to get method descriptor"))
	}

	opts, ok := desc.Options().(*descriptorpb.MethodOptions)
	if !ok {
		return "", connect.NewError(connect.CodeInternal, fmt.Errorf("unable to get method options"))
	}

	feature, ok := proto.GetExtension(opts, api.E_Feature).(*api.Feature)
	if !ok {
		return "", connect.NewError(connect.CodeInternal, fmt.Errorf("unable to get feature extension"))
	}

	return feature.GetName(), nil
}
