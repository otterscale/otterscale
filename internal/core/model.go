package core

import (
	"context"
	"strings"
	"fmt"
	"time"
	"sort"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

type ModelRepo interface {
	ListModelArtifactPVCs(ctx context.Context, config *rest.Config, namespace string) ([]corev1.PersistentVolumeClaim, error)
}

type ModelUseCase struct {
	action     ActionRepo
	facility   FacilityRepo
	kubeCore   KubeCoreRepo
	release    *ReleaseUseCase
	kubernetes *KubernetesUseCase
}

func NewModelUseCase(action ActionRepo, facility FacilityRepo, kubeCore KubeCoreRepo, release *ReleaseUseCase, kubernetes *KubernetesUseCase) *ModelUseCase {
	return &ModelUseCase{
		action:     action,
		facility:   facility,
		kubeCore:   kubeCore,
		release:    release,
		kubernetes: kubernetes,
	}
}

type ModelArtifact struct {
	Name      string
	Modelname string
	Size      int64
	CreatedAt time.Time
	Status    corev1.PersistentVolumeClaimPhase
}

type ModelGateway struct {
	Name          string
	Publicaddress string
}

type ModelScheduler struct {
	Name               string
	Modelartifactsname string
	HttpRouteCreate    bool
	BackendRefs        []BackendRef
	EppCreate          bool
	PrefillCreate      bool
}

type BackendRef struct {
	Name   string
	Weight int32
}

const (
	chartModelService = "https://github.com/llm-d-incubation/llm-d-modelservice/releases/download/llm-d-modelservice-v0.2.15/llm-d-modelservice-v0.2.15.tgz"
	chartInferencePool = "oci://registry.k8s.io/gateway-api-inference-extension/charts/inferencepool:v1.0.1"
	chartInfra  = "https://github.com/llm-d-incubation/llm-d-infra/releases/download/v1.3.3/llm-d-infra-v1.3.3.tgz"
	chartArtifact = "https://raw.githubusercontent.com/otterscale/otterscale-charts/refs/heads/main/docs/llm-d-artifact-0.1.0.tgz"

	keyLabelModelArtifact = "otterscale.com/modelartifact"
	kvConnector = "kv_connector"
)

// TODO: add back when llm infra is supported
// func (uc *ModelUseCase) CheckInfrastructureStatus(ctx context.Context, scope, facility string) (int32, error) {
// 	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
// 	if err != nil {
// 		return 0, err
// 	}
// 	rel, err := uc.release.Get(config, LLMd, LLMd)
// 	if err != nil {
// 		if errors.Is(err, driver.ErrReleaseNotFound) {
// 			return llmdInfraHealthNotInstalled, nil
// 		}
// 		return 0, err
// 	}
// 	switch {
// 	case rel.Info.Status.IsPending():
// 		return llmdInfraHealthPending, nil
// 	case rel.Info.Status == release.StatusDeployed:
// 		return llmdInfraHealthOK, nil
// 	case rel.Info.Status == release.StatusFailed:
// 		return llmdInfraHealthFailed, nil
// 	}
// 	return 0, nil
// }

func (uc *ModelUseCase)CreateModelArtifact(ctx context.Context, scope, facility, namespace, name, modelname string, size int64) (*ModelArtifact, error) {
	escapedKey := strings.ReplaceAll(keyLabelModelArtifact, ".", `\.`)

	ValuesMap := map[string]string{
		"commonLabels." + escapedKey: strings.NewReplacer(".", "-", "/", "-", "_", "-").Replace(modelname),
		"modelname": modelname,
		"pvc.name": name,
		"pvc.size": fmt.Sprintf("%d", size),
	}
	_, err := uc.release.CreateRelease(ctx, scope, facility, namespace, name, false, chartArtifact, "", ValuesMap)
	if err != nil {
		return nil, err
	}

	cfg, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}
	pvc, err := uc.kubeCore.GetPersistentVolumeClaim(ctx, cfg, namespace, name)
	if err != nil {
		return nil, err
	}
	
	ma := &ModelArtifact{
		Name:      pvc.Name,
		Modelname: modelname,
		Size:      size,
		CreatedAt: pvc.CreationTimestamp.Time,
		Status:    pvc.Status.Phase,
	}

	return ma, nil
}

func (uc *ModelUseCase) CreateModelScheduler(ctx context.Context, scope, facility, namespace, name, modelartifactsname, uri string, multinode bool, 
		nvidia, parentrefsname string, httproutecreate bool, backendRefs []BackendRef, backendrequest, request string, eppcreate, decodecreate bool, decodereplicas int32, 
		decodeargs []string, decodeenvs map[string]string, decoderesourceslimitsgpu, decoderesourceslimitsgpumem int32, prefillcreate bool, prefillreplicas int32, 
		prefillargs []string, prefillenvs map[string]string, prefillresourceslimitsgpu, prefillresourceslimitsgpumem int32, inferenceextensionreplicas int32, 
		extprocport int32, pluginsconfigfile string, targetportnumber int32) (*ModelScheduler, error) {

	ipName := fmt.Sprintf("%s-inferencepool", name)

	ipValuesMap := map[string]string{
		"inferenceExtension.replicas": fmt.Sprintf("%d", inferenceextensionreplicas),
		"inferenceExtension.extProcPort": fmt.Sprintf("%d", extprocport),
		"inferenceExtension.pluginsConfigFile": pluginsconfigfile,
		"inferencePool.targetPorts[0].number": fmt.Sprintf("%d",targetportnumber),
		"inferencePool.targetPortNumber":  fmt.Sprintf("%d",targetportnumber),
		"inferenceExtension.image.name": "llm-d-inference-scheduler",
		"inferenceExtension.image.hub": "ghcr.io/llm-d",
		"inferenceExtension.image.tag": "v0.3.2",
		"inferenceExtension.image.pullPolicy": "Always",
		"inferenceExtension.monitoring.interval": "10s",
		"inferenceExtension.monitoring.prometheus.enabled": "true",
		"inferencePool.modelServerType": "vllm",
		"inferencePool.apiVersion": "inference.networking.x-k8s.io/v1alpha2",
		"inferencePool.modelServers.matchLabels.llm-d\\.ai/inferenceServing": "true",
		"provider.name": "istio",
		"istio.destinationRule.trafficPolicy.tls.mode": "SIMPLE",
		"istio.destinationRule.trafficPolicy.tls.insecureSkipVerify": "true",
	}

	decodeEnvYAML := renderEnvMapYAML(decodeenvs, "      ")
	prefillEnvYAML := renderEnvMapYAML(prefillenvs, "      ")

	decodeArgsYAML := renderStringListYAML(append([]string{
		"--enforce-eager",
		"--kv-transfer-config",
		"{\"kv_connector\":\"NixlConnector\",\"kv_role\":\"kv_both\"}",
	}, decodeargs...), "        ")

	prefillArgsYAML := renderStringListYAML(append([]string{
		"--enforce-eager",
		"--kv-transfer-config",
		"{\"kv_connector\":\"NixlConnector\",\"kv_role\":\"kv_both\"}",
	}, prefillargs...), "        ")

	backendRefsYAML := ""
	if httproutecreate && len(backendRefs) > 0 {
		backendRefsYAML = renderBackendRefsYAML(backendRefs, "        ")
	}

	valuesYAML := fmt.Sprintf(
	`fullnameOverride: %q
modelArtifacts:
  name: %q
  uri: %q
multinode: %t
accelerator:
  type: "nvidia"
  resources:
    nvidia: %q
routing:
  servicePort: 8000
  proxy:
    image: "ghcr.io/llm-d/llm-d-routing-sidecar:v0.3.0"
    connector: "nixlv2"
    secure: false
  parentRefs:
  - name: %q
  inferencePool:
    create: false
  httpRoute:
    create: %t
    timeouts:
      backendRequest: %q
      request: %q
    rules:
      - matches:
          - path:
              type: PathPrefix
              value: "/"
%s
  epp:
    create: false

decode:
  create: %t
  replicas: %d
  monitoring:
    podmonitor:
      enabled: true
      portName: "metrics"
      path: "/metrics"
      interval: "30s"
  containers:
  - name: "vllm"
    image: "ghcr.io/llm-d/llm-d:v0.2.0"
    modelCommand: "vllmServe"
    args:
%s
    ports:
      - containerPort: 5557
        protocol: TCP
      - containerPort: 8200
        name: metrics
        protocol: TCP
    mountModelVolume: true
    volumeMounts:
    - name: "metrics-volume"
      mountPath: "/.config"
    - name: "torch-compile-cache"
      mountPath: "/.compile-cache"
    env:
%s
    resources:
      limits:
        otterscale.com/vgpu: "1"
        otterscale.com/vgpumem-percentage: "40"
      requests:
        otterscale.com/vgpu: "1"
        otterscale.com/vgpumem-percentage: "40"
  volumes:
  - name: "metrics-volume"
    emptyDir: {}
  - name: "torch-compile-cache"
    emptyDir: {}
prefill:
  create: %t
  replicas: %d
  containers:
  - name: "vllm"
    image: "ghcr.io/llm-d/llm-d:v0.2.0"
    modelCommand: "vllmServe"
    args:
%s
    env:
%s
`,
		// top
		name,
		modelartifactsname, uri,
		multinode,
		nvidia,
		// routing
		parentrefsname,
		httproutecreate, backendrequest, request,
		backendRefsYAML,
		// decode
		decodecreate, decodereplicas,
		decodeArgsYAML,
		decodeEnvYAML,
		// prefill
		prefillcreate, prefillreplicas,
		prefillArgsYAML,
		prefillEnvYAML,
	)
	_, err := uc.release.CreateRelease(ctx, scope, facility, namespace, name, false, chartModelService, valuesYAML, nil)
	if err != nil {
		return nil, err
	}

	if eppcreate {
		_, err := uc.release.CreateRelease(ctx, scope, facility, namespace, ipName, false, chartInferencePool, "", ipValuesMap)
		if err != nil {
			return nil, err
		}
	}

	ms := &ModelScheduler{
		Name:      name,
		Modelartifactsname: modelartifactsname,
		HttpRouteCreate:    httproutecreate,
		BackendRefs:        backendRefs,
		EppCreate:          eppcreate,
		PrefillCreate:      prefillcreate,
	}
	return ms, nil
}

func escapeDots(s string) string { return strings.ReplaceAll(s, ".", `\.`) }

func deriveMemPctKey(nvidia string) string {
	n := strings.TrimSpace(nvidia)
	if n == "" {
		return "vgpumem-percentage" 
	}
	if idx := strings.IndexRune(n, '/'); idx > 0 {
		return n[:idx] + "/vgpumem-percentage"
	}
	return n + "/vgpumem-percentage"
}

func parseEnvKV(list []string) [][2]string {
	out := make([][2]string, 0, len(list))
	for _, item := range list {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		parts := strings.SplitN(item, "=", 2)
		if len(parts) == 2 {
			out = append(out, [2]string{strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])})
		}
	}
	return out
}

func renderEnvMapYAML(m map[string]string, indent string) string {
	if len(m) == 0 {
		return indent + "[]"
	}
	keys := make([]string, 0, len(m))
	for k := range m { keys = append(keys, k) }
	sort.Strings(keys)

	var b strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&b, "%s- name: %q\n", indent, k)
		fmt.Fprintf(&b, "%s  value: %q\n", indent, m[k])
	}
	return b.String()
}

func renderStringListYAML(items []string, indent string) string {
	if len(items) == 0 {
		return indent + "[]"
	}
	var b strings.Builder
	for _, it := range items {
		fmt.Fprintf(&b, "%s- %q\n", indent, it)
	}
	return b.String()
}

func renderBackendRefsYAML(brs []BackendRef, indent string) string {
	if len(brs) == 0 {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%sbackendRefs:\n", indent)
	for _, br := range brs {
		fmt.Fprintf(&b, "%s- group: inference.networking.x-k8s.io\n", indent)
		fmt.Fprintf(&b, "%s  kind: InferencePool\n", indent)
		fmt.Fprintf(&b, "%s  name: %q\n", indent, br.Name)
		fmt.Fprintf(&b, "%s  port: 8000\n", indent)
		if br.Weight > 0 {
			fmt.Fprintf(&b, "%s  weight: %d\n", indent, br.Weight)
		}
	}
	return b.String()
}

func (uc *ModelUseCase) CreateModelGateway(ctx context.Context, scope, facility, namespace, name string, cpu, memory int32, servicetype string ) (*ModelGateway, error) {
	ValuesMap := map[string]string{
		"nameOverride": name,
		"gateway.gatewayParameters.resources.limits.cpu": fmt.Sprintf("%d", cpu),
		"gateway.gatewayParameters.resources.limits.memory": fmt.Sprintf("%d", memory),
		"gateway.service.type": servicetype,
	}
	release, err := uc.release.CreateRelease(ctx, scope, facility, namespace, name, false, chartInfra, "", ValuesMap)
	if err != nil {
		return nil, err
	}
	print(release)
	pa, err:= uc.kubernetes.GetPublicAddress(ctx, scope, facility)
	if err != nil {
		return nil, err
	}

	mg := &ModelGateway{
		Name:          name,
		Publicaddress: pa,
	}
	return mg, nil
}

func (uc *ModelUseCase) ListModelArtifactPVCs(ctx context.Context, scope, facility, namespace string,) ([]corev1.PersistentVolumeClaim, error) {
	const LabelModelArtifact = "otterscale.com/modelartifact"

	cfg, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	pvcs, err := uc.kubeCore.ListPersistentVolumeClaims(ctx, cfg, namespace)
	if err != nil {
		return nil, err
	}

	ma := make([]corev1.PersistentVolumeClaim, 0, len(pvcs))
	for _, pvc := range pvcs {
		if pvc.Labels != nil {
			if _, ok := pvc.Labels[LabelModelArtifact]; ok {
				ma = append(ma, pvc)
			}
		}
	}
	return ma, nil
}