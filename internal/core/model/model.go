package model

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "sigs.k8s.io/gateway-api-inference-extension/api/v1"
	gav1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/otterscale/otterscale/internal/core/application/cluster"
	"github.com/otterscale/otterscale/internal/core/application/persistent"
	"github.com/otterscale/otterscale/internal/core/application/release"
	"github.com/otterscale/otterscale/internal/core/application/service"
	"github.com/otterscale/otterscale/internal/core/application/workload"
	"github.com/otterscale/otterscale/internal/core/versions"
)

const (
	ModelNameAnnotation = "otterscale.com/model.name"
	ModelPVCAnnotation  = "otterscale.com/model.pvc"
)

const (
	vgpuResource              = "otterscale.com/vgpu"
	vgpuMemPercentageResource = "otterscale.com/vgpumem-percentage"
)

type Model struct {
	ID             string
	Release        *release.Release
	Mode           Mode
	Prefill        *Prefill
	Decode         *Decode
	MaxModelLength uint32
	Pods           []workload.Pod
	FromPVC        bool
	PVCName        string
}

type Prefill struct {
	Replica    uint32
	Tensor     uint32
	VGPUMemory uint32
}

type Decode struct {
	Replica    uint32
	Tensor     uint32
	VGPUMemory uint32
}

type UseCase struct {
	inferencePool InferencePoolRepo

	deployment            workload.DeploymentRepo
	gateway               service.GatewayRepo
	httpRoute             service.HTTPRouteRepo
	persistentVolumeClaim persistent.PersistentVolumeClaimRepo
	pod                   workload.PodRepo
	node                  cluster.NodeRepo
	release               release.ReleaseRepo
	service               service.ServiceRepo
}

func NewUseCase(inferencePool InferencePoolRepo, deployment workload.DeploymentRepo, gateway service.GatewayRepo, httpRoute service.HTTPRouteRepo, persistentVolumeClaim persistent.PersistentVolumeClaimRepo, pod workload.PodRepo, node cluster.NodeRepo, release release.ReleaseRepo, service service.ServiceRepo) *UseCase {
	return &UseCase{
		inferencePool:         inferencePool,
		deployment:            deployment,
		gateway:               gateway,
		httpRoute:             httpRoute,
		persistentVolumeClaim: persistentVolumeClaim,
		pod:                   pod,
		node:                  node,
		release:               release,
		service:               service,
	}
}

func (uc *UseCase) ListModels(ctx context.Context, scope, namespace string) (models []Model, uri string, err error) {
	selector := release.TypeLabel + "=" + "model"

	releases, err := uc.release.List(ctx, scope, namespace, selector)
	if err != nil {
		return nil, "", err
	}

	for i := range releases {
		modelName, ok := extractModelName(releases[i].Config)
		if !ok {
			slog.Error("model name not found in release config", "release", releases[i].Name)
			continue
		}

		mode := ModeIntelligentInferenceScheduling
		prefill := extractPrefill(releases[i].Config)
		decode := extractDecode(releases[i].Config)
		maxModelLength := extractMaxModelLength(releases[i].Config)

		if prefill != nil {
			mode = ModePrefillDecodeDisaggregation
		}

		selector := "llm-d.ai/model" + "=" + releases[i].Name

		pods, err := uc.pod.List(ctx, scope, namespace, selector)
		if err != nil {
			return nil, "", err
		}

		pvcName, fromPVC := releases[i].Labels[ModelPVCAnnotation]

		models = append(models, Model{
			ID:             modelName,
			Release:        &releases[i],
			Mode:           mode,
			Prefill:        prefill,
			Decode:         decode,
			MaxModelLength: maxModelLength,
			Pods:           pods,
			FromPVC:        fromPVC,
			PVCName:        pvcName,
		})
	}

	uri, err = uc.gatewayURL(ctx, scope, "llm-d", "llm-d-infra-inference-gateway-istio")
	if err != nil {
		return nil, "", err
	}

	return models, uri, nil
}

func inferencePoolName(name string) string {
	return "inferencepool-" + shortID(name)
}

func httpRouteName(name string) string {
	return "httproute-" + shortID(name)
}

func (uc *UseCase) CreateModel(ctx context.Context, scope, namespace, name, modelName string, fromPVC bool, pvcName string, sizeBytes uint64, mode Mode, prefill *Prefill, decode *Decode, maxModelLength uint32) (*Model, error) {
	gatewayName := "llm-d-infra-inference-gateway" // from llm-d-infra helm chart (.Values.nameOverride)
	inferencePoolName := inferencePoolName(name)
	inferencePoolPort := int32(8000) //nolint:mnd // default port for inference pool
	httpRouteName := httpRouteName(name)

	// check gateway exists
	if _, err := uc.gateway.Get(ctx, scope, namespace, gatewayName); err != nil {
		return nil, err
	}

	// reconcile inference pool
	if err := uc.reconcileInferencePool(ctx, scope, namespace, inferencePoolName, modelName, name, mode, inferencePoolPort); err != nil {
		return nil, err
	}

	// reconcile http route
	if err := uc.reconcileHTTPRoute(ctx, scope, namespace, httpRouteName, name, gatewayName, inferencePoolName, inferencePoolPort); err != nil {
		return nil, err
	}

	// deploy model service
	release, err := uc.installModelService(ctx, scope, namespace, name, modelName, fromPVC, pvcName, sizeBytes, mode, prefill, decode, maxModelLength)
	if err != nil {
		return nil, err
	}

	newMode := ModeIntelligentInferenceScheduling
	newPrefill := extractPrefill(release.Config)
	newDecode := extractDecode(release.Config)
	newMaxModelLength := extractMaxModelLength(release.Config)

	if newPrefill != nil {
		newMode = ModePrefillDecodeDisaggregation
	}

	newPVCName, newFromPVC := release.Labels[ModelPVCAnnotation]

	return &Model{
		ID:             modelName,
		Release:        release,
		Mode:           newMode,
		Prefill:        newPrefill,
		Decode:         newDecode,
		MaxModelLength: newMaxModelLength,
		FromPVC:        newFromPVC,
		PVCName:        newPVCName,
	}, nil
}

func (uc *UseCase) UpdateModel(ctx context.Context, scope, namespace, name string, mode Mode, prefill *Prefill, decode *Decode, maxModelLength uint32) (*Model, error) {
	rel, err := uc.release.Get(ctx, scope, namespace, name)
	if err != nil {
		return nil, err
	}

	v, ok := rel.Config["modelArtifacts"]
	if !ok {
		return nil, fmt.Errorf("modelArtifacts not found in release config")
	}

	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("modelArtifacts is not a map")
	}

	modelName, ok := m["name"].(string)
	if !ok {
		return nil, fmt.Errorf("modelArtifacts.name is not a string")
	}

	strSizeBytes, ok := m["size"].(string)
	if !ok {
		return nil, fmt.Errorf("modelArtifacts.size is not a string")
	}

	strSizeBytes = strings.Trim(strSizeBytes, `"`)
	sizeBytes, err := strconv.ParseUint(strSizeBytes, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse modelArtifacts.size: %w", err)
	}

	release, err := uc.upgradeModelService(ctx, scope, namespace, name, modelName, sizeBytes, mode, prefill, decode, maxModelLength)
	if err != nil {
		return nil, err
	}

	newMode := ModeIntelligentInferenceScheduling
	newPrefill := extractPrefill(release.Config)
	newDecode := extractDecode(release.Config)
	newMaxModelLength := extractMaxModelLength(release.Config)

	if newPrefill != nil {
		newMode = ModePrefillDecodeDisaggregation
	}

	pvcName, fromPVC := release.Labels[ModelPVCAnnotation]

	return &Model{
		ID:             modelName,
		Release:        release,
		Mode:           newMode,
		Prefill:        newPrefill,
		Decode:         newDecode,
		MaxModelLength: newMaxModelLength,
		FromPVC:        fromPVC,
		PVCName:        pvcName,
	}, nil
}

func (uc *UseCase) DeleteModel(ctx context.Context, scope, namespace, name string) error {
	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		httpRouteName := httpRouteName(name)
		return uc.httpRoute.Delete(egctx, scope, namespace, httpRouteName)
	})

	eg.Go(func() error {
		inferencePoolName := inferencePoolName(name)
		_, err := uc.release.Uninstall(egctx, scope, namespace, inferencePoolName, false)
		return err
	})

	eg.Go(func() error {
		_, err := uc.release.Uninstall(egctx, scope, namespace, name, false)
		return err
	})

	return eg.Wait()
}

func (uc *UseCase) gatewayURL(ctx context.Context, scope, namespace, serviceName string) (string, error) {
	internalIP, err := uc.node.InternalIP(ctx, scope)
	if err != nil {
		return "", err
	}

	service, err := uc.service.Get(ctx, scope, namespace, serviceName)
	if err != nil {
		return "", err
	}

	var nodePort int32
	for _, sp := range service.Spec.Ports {
		if sp.Name == "default" {
			nodePort = sp.NodePort
			break
		}
	}

	if nodePort == 0 {
		return "", fmt.Errorf("default node port not found for llm-d-infra-inference-gateway-istio service")
	}

	return fmt.Sprintf("http://%s:%d", internalIP, nodePort), nil
}

func (uc *UseCase) reconcileHTTPRoute(ctx context.Context, scope, namespace, name, releaseName, gatewayName, inferencePoolName string, inferencePoolPort int32) error {
	_, err := uc.httpRoute.Get(ctx, scope, namespace, name)

	if k8serrors.IsNotFound(err) {
		httpRoute := uc.buildHTTPRoute(namespace, name, releaseName, gatewayName, inferencePoolName, inferencePoolPort)

		_, err := uc.httpRoute.Create(ctx, scope, namespace, httpRoute)
		return err
	}

	return err
}

func (uc *UseCase) buildHTTPRoute(namespace, name, releaseName, gatewayName, inferencePoolName string, inferencePoolPort int32) *service.HTTPRoute {
	// parent reference
	parentGroup := gav1.Group(gav1.GroupName)
	parentKind := gav1.Kind("Gateway")

	// matches
	pathMatchType := gav1.PathMatchPathPrefix
	headerMatchType := gav1.HeaderMatchExact

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
							Headers: []gav1.HTTPHeaderMatch{
								{
									Type:  &headerMatchType,
									Name:  "OtterScale-Model-Name",
									Value: releaseName,
								},
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

func (uc *UseCase) reconcileInferencePool(ctx context.Context, scope, namespace, name, modelName, releaseName string, mode Mode, port int32) error {
	_, err := uc.inferencePool.Get(ctx, scope, namespace, name)

	if k8serrors.IsNotFound(err) {
		return uc.installInferencePool(ctx, scope, namespace, name, modelName, releaseName, mode, port)
	}

	return err
}

func (uc *UseCase) installInferencePool(ctx context.Context, scope, namespace, name, modelName, releaseName string, mode Mode, port int32) error {
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
	valuesMap := convertGAIEValuesMap(mode, name, releaseName, port)

	_, err := uc.release.Install(ctx, scope, namespace, name, false, chartRef, labels, labels, annotations, "", valuesMap)
	return err
}

func (uc *UseCase) installModelService(ctx context.Context, scope, namespace, name, modelName string, fromPVC bool, pvcName string, sizeBytes uint64, mode Mode, prefill *Prefill, decode *Decode, maxModelLength uint32) (*release.Release, error) {
	// chart ref
	chartRef := fmt.Sprintf("https://github.com/llm-d-incubation/llm-d-modelservice/releases/download/llm-d-modelservice-v%[1]s/llm-d-modelservice-v%[1]s.tgz", versions.LLMDModelService)

	// labels
	labels := map[string]string{
		release.TypeLabel: "model",
	}

	if fromPVC {
		labels[ModelPVCAnnotation] = pvcName
	}

	// annotations
	annotations := map[string]string{
		ModelNameAnnotation: modelName,
	}

	// values
	valuesMap := convertModelServiceValuesMap(mode, name, modelName, fromPVC, pvcName, sizeBytes, prefill, decode, maxModelLength)

	return uc.release.Install(ctx, scope, namespace, name, false, chartRef, labels, labels, annotations, "", valuesMap)
}

func (uc *UseCase) upgradeModelService(ctx context.Context, scope, namespace, name, modelName string, sizeBytes uint64, mode Mode, prefill *Prefill, decode *Decode, maxModelLength uint32) (*release.Release, error) {
	// chart ref
	chartRef := fmt.Sprintf("https://github.com/llm-d-incubation/llm-d-modelservice/releases/download/llm-d-modelservice-v%[1]s/llm-d-modelservice-v%[1]s.tgz", versions.LLMDModelService)

	rel, err := uc.release.Get(ctx, scope, namespace, name)
	if err != nil {
		return nil, err
	}

	pvcName, fromPVC := rel.Labels[ModelPVCAnnotation]

	// values
	valuesMap := convertModelServiceValuesMap(mode, name, modelName, fromPVC, pvcName, sizeBytes, prefill, decode, maxModelLength)

	return uc.release.Upgrade(ctx, scope, namespace, name, false, chartRef, "", valuesMap, false)
}

func extractModelName(config map[string]any) (string, bool) {
	modelArtifacts, ok := config["modelArtifacts"].(map[string]any)
	if !ok {
		return "", false
	}

	name, ok := modelArtifacts["name"].(string)
	return name, ok
}

func extractMaxModelLength(config map[string]any) uint32 {
	decode, ok := config["decode"].(map[string]any)
	if !ok {
		return 0
	}

	containers, ok := decode["containers"].([]any)
	if !ok || len(containers) == 0 {
		return 0
	}

	container, ok := containers[0].(map[string]any)
	if !ok {
		return 0
	}

	args, ok := container["args"].([]any)
	if !ok {
		return 0
	}

	return parseMaxModelLenFromArgs(args)
}

func parseMaxModelLenFromArgs(args []any) uint32 {
	for i, arg := range args {
		argStr, ok := arg.(string)
		if !ok || argStr != "--max-model-len" || i+1 >= len(args) {
			continue
		}

		valStr, ok := args[i+1].(string)
		if !ok {
			return 0
		}

		valStr = strings.TrimPrefix(valStr, "[str]")
		val, err := strconv.ParseUint(valStr, 10, 32)
		if err != nil {
			return 0
		}

		return uint32(val)
	}

	return 0
}

func extractPrefill(config map[string]any) *Prefill {
	prefillConfig, ok := config["prefill"].(map[string]any)
	if !ok {
		return nil
	}

	replica := extractReplica(prefillConfig)
	tensor := extractTensor(prefillConfig)
	vgpuMemory := extractVGPUMemory(prefillConfig)

	if replica == 0 || vgpuMemory == 0 {
		return nil
	}

	return &Prefill{
		Replica:    replica,
		Tensor:     tensor,
		VGPUMemory: vgpuMemory,
	}
}

func extractDecode(config map[string]any) *Decode {
	decodeConfig, ok := config["decode"].(map[string]any)
	if !ok {
		return nil
	}

	replica := extractReplica(decodeConfig)
	tensor := extractTensor(decodeConfig)
	vgpuMemory := extractVGPUMemory(decodeConfig)

	if replica == 0 || tensor == 0 || vgpuMemory == 0 {
		return nil
	}

	return &Decode{
		Replica:    replica,
		Tensor:     tensor,
		VGPUMemory: vgpuMemory,
	}
}

func extractReplica(config map[string]any) uint32 {
	return parseResourceValue(config, "replicas")
}

func extractTensor(config map[string]any) uint32 {
	parallelism, ok := config["parallelism"].(map[string]any)
	if !ok {
		return 1
	}

	return parseResourceValue(parallelism, "tensor")
}

func extractVGPUMemory(config map[string]any) uint32 {
	containers, ok := config["containers"].([]any)
	if !ok || len(containers) == 0 {
		return 0
	}

	container, ok := containers[0].(map[string]any)
	if !ok {
		return 0
	}

	resources, ok := container["resources"].(map[string]any)
	if !ok {
		return 0
	}

	requests, ok := resources["requests"].(map[string]any)
	if !ok {
		return 0
	}

	return parseResourceValue(requests, vgpuMemPercentageResource)
}

func parseResourceValue(config map[string]any, key string) uint32 {
	value, ok := config[key]
	if !ok {
		return 0
	}

	switch v := value.(type) {
	case float64:
		return uint32(v)

	case int64:
		return uint32(v) //nolint:gosec // safe conversion

	default:
		return 0
	}
}

func shortID(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:4])
}
