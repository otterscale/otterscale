package model

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "sigs.k8s.io/gateway-api-inference-extension/api/v1"
	gav1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/otterscale/otterscale/internal/core/application/chart"
	"github.com/otterscale/otterscale/internal/core/application/persistent"
	"github.com/otterscale/otterscale/internal/core/application/release"
	"github.com/otterscale/otterscale/internal/core/application/service"
	"github.com/otterscale/otterscale/internal/core/application/workload"
	"github.com/otterscale/otterscale/internal/core/versions"
)

const ModelNameAnnotation = "otterscale.com/model.name"

// Release represents a Helm Release resource.
type Model = release.Release

type Resource struct {
	VGPU       uint32
	VGPUMemory uint32
}

type UseCase struct {
	inferencePool InferencePoolRepo

	chart                 chart.ChartRepo
	deployment            workload.DeploymentRepo
	gateway               service.GatewayRepo
	httpRoute             service.HTTPRouteRepo
	persistentVolumeClaim persistent.PersistentVolumeClaimRepo
	release               release.ReleaseRepo
	service               service.ServiceRepo
}

func NewUseCase(inferencePool InferencePoolRepo, chart chart.ChartRepo, deployment workload.DeploymentRepo, gateway service.GatewayRepo, httpRoute service.HTTPRouteRepo, persistentVolumeClaim persistent.PersistentVolumeClaimRepo, release release.ReleaseRepo, service service.ServiceRepo) *UseCase {
	return &UseCase{
		inferencePool:         inferencePool,
		chart:                 chart,
		deployment:            deployment,
		gateway:               gateway,
		httpRoute:             httpRoute,
		persistentVolumeClaim: persistentVolumeClaim,
		release:               release,
		service:               service,
	}
}

func (uc *UseCase) ListModels(ctx context.Context, scope, namespace string) (models []Model, uri string, err error) {
	selector := release.TypeLabel + "=" + "model"

	models, err = uc.release.List(ctx, scope, namespace, selector)
	if err != nil {
		return nil, "", err
	}

	uri, err = uc.gatewayURL(ctx, scope, "llm-d", "llm-d-infra-inference-gateway-istio")
	if err != nil {
		return nil, "", err
	}

	return models, uri, nil
}

func (uc *UseCase) CreateModel(ctx context.Context, scope, namespace, name, modelName string, sizeBytes uint64, limits, requests *Resource) (*Model, error) {
	gatewayName := "llm-d-infra-inference-gateway" // from llm-d-infra helm chart (.Values.nameOverride)
	inferencePoolName := "inferencepool-" + shortID(modelName)
	inferencePoolPort := int32(8000) //nolint:mnd // default port for inference pool
	httpRouteName := "httproute-" + shortID(gatewayName+inferencePoolName)

	// check gateway exists
	if _, err := uc.gateway.Get(ctx, scope, namespace, gatewayName); err != nil {
		return nil, err
	}

	// reconcile inference pool
	if err := uc.reconcileInferencePool(ctx, scope, namespace, inferencePoolName, modelName, inferencePoolPort); err != nil {
		return nil, err
	}

	// reconcile http route
	if err := uc.reconcileHTTPRoute(ctx, scope, namespace, httpRouteName, gatewayName, inferencePoolName, inferencePoolPort); err != nil {
		return nil, err
	}

	// deploy model service
	return uc.installModelService(ctx, scope, namespace, name, modelName, sizeBytes, limits, requests)
}

func (uc *UseCase) UpdateModel(ctx context.Context, scope, namespace, name string, requests, limits *Resource) (*Model, error) {
	rel, err := uc.release.Get(ctx, scope, namespace, name)
	if err != nil {
		return nil, err
	}

	vname, ok := rel.Config["modelArtifacts.name"]
	if !ok {
		return nil, fmt.Errorf("modelArtifacts.name not found in release config")
	}

	vsize, ok := rel.Config["modelArtifacts.size"]
	if !ok {
		return nil, fmt.Errorf("modelArtifacts.size not found in release config")
	}

	modelName, ok := vname.(string)
	if !ok {
		return nil, fmt.Errorf("modelArtifacts.name is not a string")
	}

	sizeBytes, ok := vsize.(int)
	if !ok {
		return nil, fmt.Errorf("modelArtifacts.size is not a int64")
	}

	return uc.upgradeModelService(ctx, scope, namespace, name, modelName, uint64(sizeBytes), limits, requests) //nolint:gosec // ignore
}

func (uc *UseCase) DeleteModel(ctx context.Context, scope, namespace, name string) error {
	_, err := uc.release.Uninstall(ctx, scope, namespace, name, false)
	return err
}

func (uc *UseCase) gatewayURL(ctx context.Context, scope, namespace, serviceName string) (string, error) {
	url, err := uc.service.URL(scope)
	if err != nil {
		return "", err
	}

	service, err := uc.service.Get(ctx, scope, namespace, serviceName)
	if err != nil {
		return "", err
	}

	var port int32
	for _, sp := range service.Spec.Ports {
		if sp.Name == "default" {
			port = sp.Port
			break
		}
	}

	if port == 0 {
		return "", fmt.Errorf("default port not found for llm-d-infra-inference-gateway-istio service")
	}

	return fmt.Sprintf("http://%s:%d", url.Hostname(), port), nil
}

func (uc *UseCase) reconcileHTTPRoute(ctx context.Context, scope, namespace, name, gatewayName, inferencePoolName string, inferencePoolPort int32) error {
	_, err := uc.httpRoute.Get(ctx, scope, namespace, name)

	if k8serrors.IsNotFound(err) {
		httpRoute := uc.buildHTTPRoute(namespace, name, gatewayName, inferencePoolName, inferencePoolPort)

		_, err := uc.httpRoute.Create(ctx, scope, namespace, httpRoute)
		return err
	}

	return err
}

func (uc *UseCase) buildHTTPRoute(namespace, name, gatewayName, inferencePoolName string, inferencePoolPort int32) *service.HTTPRoute {
	// parent reference
	parentGroup := gav1.Group(gav1.GroupName)
	parentKind := gav1.Kind("Gateway")

	// matches
	pathMatchType := gav1.PathMatchPathPrefix

	// backend references
	backendGroup := gav1.Group(v1.GroupName)
	backendKind := gav1.Kind("InferencePool")
	weight := int32(1)

	// timeouts
	timeout := gav1.Duration("0s")

	return &service.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: gav1.HTTPRouteSpec{
			CommonRouteSpec: gav1.CommonRouteSpec{
				ParentRefs: []gav1.ParentReference{
					{
						Group: &parentGroup,
						Kind:  &parentKind,
						Name:  gav1.ObjectName(gatewayName), // from llm-d-infra helm chart (.Values.nameOverride)
					},
				},
			},
			Rules: []gav1.HTTPRouteRule{
				{
					Matches: []gav1.HTTPRouteMatch{
						{
							Path: &gav1.HTTPPathMatch{
								Type: &pathMatchType,
							},
						},
					},
					BackendRefs: []gav1.HTTPBackendRef{
						{
							BackendRef: gav1.BackendRef{
								BackendObjectReference: gav1.BackendObjectReference{
									Group: &backendGroup,
									Kind:  &backendKind,
									Name:  gav1.ObjectName(inferencePoolName),
									Port:  &inferencePoolPort,
								},
								Weight: &weight,
							},
						},
					},
					Timeouts: &gav1.HTTPRouteTimeouts{
						Request:        &timeout,
						BackendRequest: &timeout,
					},
				},
			},
		},
	}
}

func (uc *UseCase) reconcileInferencePool(ctx context.Context, scope, namespace, name, modelName string, port int32) error {
	_, err := uc.inferencePool.Get(ctx, scope, namespace, name)

	if k8serrors.IsNotFound(err) {
		return uc.installInferencePool(ctx, scope, namespace, name, modelName, port)
	}

	return err
}

func (uc *UseCase) installInferencePool(ctx context.Context, scope, namespace, name, modelName string, port int32) error {
	// chart ref
	chartRef := fmt.Sprintf("oci://registry.k8s.io/gateway-api-inference-extension/charts/inferencepool:v%s", versions.GatewayAPIInferenceExtension)

	// labels
	labels := map[string]string{
		release.TypeLabel: "inference-pool",
	}

	// annotations
	annotations := map[string]string{
		ModelNameAnnotation: modelName,
	}

	// values
	valuesYAML := fmt.Sprintf(inferencePoolValuesYAML, name, port)

	_, err := uc.release.Install(ctx, scope, namespace, name, false, chartRef, labels, labels, annotations, valuesYAML, nil)
	if err != nil {
		return err
	}

	// FIXME: workaround to remove tracing arg until supported in llm-d-inference-scheduler
	time.Sleep(2 * time.Second)

	deployment, err := uc.deployment.Get(ctx, scope, namespace, name+"-epp")
	if err != nil {
		return err
	}

	containers := deployment.Spec.Template.Spec.Containers
	if len(containers) == 0 {
		return fmt.Errorf("no containers found in deployment %s", deployment.Name)
	}

	deployment.Spec.Template.Spec.Containers[0].Args = slices.DeleteFunc(containers[0].Args, func(arg string) bool {
		return arg == "--tracing=false"
	})

	_, err = uc.deployment.Update(ctx, scope, namespace, deployment)
	return err
}

func (uc *UseCase) installModelService(ctx context.Context, scope, namespace, name, modelName string, sizeBytes uint64, limits, requests *Resource) (*Model, error) {
	// chart ref
	chartRef := fmt.Sprintf("https://github.com/llm-d-incubation/llm-d-modelservice/releases/download/llm-d-modelservice-v%[1]s/llm-d-modelservice-v%[1]s.tgz", versions.LLMDModelService)

	// labels
	labels := map[string]string{
		release.TypeLabel: "model",
	}

	// annotations
	annotations := map[string]string{
		ModelNameAnnotation: modelName,
	}

	// values
	strSizeBytes := strconv.Itoa(int(sizeBytes)) //nolint:gosec // ignore
	valuesYAML := fmt.Sprintf(modelServiceValuesYAML, modelName, formatLabel(modelName), strSizeBytes, limits.VGPU, limits.VGPUMemory, requests.VGPU, requests.VGPUMemory)

	return uc.release.Install(ctx, scope, namespace, name, false, chartRef, labels, labels, annotations, valuesYAML, nil)
}

func (uc *UseCase) upgradeModelService(ctx context.Context, scope, namespace, name, modelName string, sizeBytes uint64, limits, requests *Resource) (*Model, error) {
	// chart ref
	chartRef := fmt.Sprintf("https://github.com/llm-d-incubation/llm-d-modelservice/releases/download/llm-d-modelservice-v%[1]s/llm-d-modelservice-v%[1]s.tgz", versions.LLMDModelService)

	// values
	strSizeBytes := strconv.Itoa(int(sizeBytes)) //nolint:gosec // ignore
	valuesYAML := fmt.Sprintf(modelServiceValuesYAML, modelName, formatLabel(modelName), strSizeBytes, limits.VGPU, limits.VGPUMemory, requests.VGPU, requests.VGPUMemory)

	return uc.release.Upgrade(ctx, scope, namespace, name, false, chartRef, valuesYAML, nil, false)
}

func shortID(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:4])
}

func formatLabel(input string) string {
	regex := regexp.MustCompile("[^a-zA-Z0-9]+")
	output := regex.ReplaceAllString(input, "-")
	trimRegex := regexp.MustCompile("^-|-$")
	finalOutput := trimRegex.ReplaceAllString(output, "")
	return strings.ToLower(finalOutput)
}
