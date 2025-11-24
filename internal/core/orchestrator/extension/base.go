package extension

import (
	"fmt"

	"github.com/otterscale/otterscale/internal/core/application/release"
	"github.com/otterscale/otterscale/internal/core/versions"
)

type chartManifest struct {
	Ref         string
	Version     string
	Labels      map[string]string
	Annotations map[string]string
	ValuesMap   map[string]string
	PostFunc    func(scope string) error
}

type crdManifest struct {
	Ref                  string
	Version              string
	AnnotationVersionKey string
}

type base struct {
	Name        string
	Namespace   string
	DisplayName string
	Description string
	Logo        string
	Charts      []chartManifest
	CRD         *crdManifest
}

var (
	general = []base{
		{
			Name:        "gateway-api",
			DisplayName: "Gateway API",
			Description: "Gateway API is an official Kubernetes project focused on L4 and L7 routing in Kubernetes.",
			Logo:        "https://github.com/kubernetes-sigs.png",
			CRD: &crdManifest{
				Ref:                  fmt.Sprintf("https://github.com/kubernetes-sigs/gateway-api.git/config/crd?ref=v%s", versions.GatewayAPI),
				Version:              versions.GatewayAPI,
				AnnotationVersionKey: "gateway.networking.k8s.io/bundle-version",
			},
		},
		{
			Name:        "gateway-api-inference-extension",
			DisplayName: "Gateway API Inference Extension",
			Description: "Gateway API Inference Extension is an official Kubernetes project that optimizes self-hosting Generative Models on Kubernetes.",
			Logo:        "https://github.com/kubernetes-sigs.png",
			CRD: &crdManifest{
				Ref:                  fmt.Sprintf("https://github.com/kubernetes-sigs/gateway-api-inference-extension.git/config/crd?ref=v%s", versions.GatewayAPIInferenceExtension),
				Version:              versions.GatewayAPIInferenceExtension,
				AnnotationVersionKey: "inference.networking.k8s.io/bundle-version",
			},
		},
		{
			Name:        "base",
			Namespace:   "istio-system",
			DisplayName: "Istio",
			Description: "Istio is an open source service mesh that layers transparently onto existing distributed applications.",
			Logo:        "https://github.com/istio.png",
			Charts: []chartManifest{
				{
					Ref:     fmt.Sprintf("https://istio-release.storage.googleapis.com/charts/base-%s.tgz", versions.Istio),
					Version: versions.Istio,
					Labels: map[string]string{
						release.TypeLabel: "extension",
					},
				},
				{
					Ref:     fmt.Sprintf("https://istio-release.storage.googleapis.com/charts/istiod-%s.tgz", versions.Istio),
					Version: versions.Istio,
					Labels: map[string]string{
						release.TypeLabel: "extension",
					},
					ValuesMap: map[string]string{
						"env.ENABLE_GATEWAY_API_INFERENCE_EXTENSION": "true",
					},
				},
			},
		},
		{
			Name:        "kube-prometheus-stack",
			Namespace:   "monitoring",
			DisplayName: "Prometheus Stack",
			Description: "Installs the kube-prometheus stack for easy, end-to-end Kubernetes cluster monitoring using the Prometheus Operator.",
			Logo:        "https://github.com/prometheus-community.png",
			Charts: []chartManifest{
				{
					Ref:     fmt.Sprintf("https://github.com/prometheus-community/helm-charts/releases/download/kube-prometheus-stack-%[1]s/kube-prometheus-stack-%[1]s.tgz", versions.KubePrometheusStack),
					Version: versions.KubePrometheusStack,
					Labels: map[string]string{
						release.TypeLabel: "extension",
					},
				},
			},
		},
	}

	registry = []base{
		{
			Name:        "kubevirt-infra", // helm chart name
			Namespace:   "distribution",
			DisplayName: "Registry",
			Description: "",
			Logo:        "https://github.com/distribution.png",
			Charts: []chartManifest{
				{
					Ref:     "",
					Version: "",
					Labels: map[string]string{
						release.TypeLabel: "extension",
					},
					PostFunc: func(_ string) error {
						// Run Juju Config on Scope
						return nil
					},
				},
			},
		},
	}

	model = []base{
		{
			Name:        "gpu-operator",
			Namespace:   "gpu-operator",
			DisplayName: "GPU Operator",
			Description: "GPU Operator creates, configures, and manages GPUs in Kubernetes.",
			Logo:        "https://github.com/otterscale.png",
			Charts: []chartManifest{
				{
					Ref:     fmt.Sprintf("https://github.com/otterscale/charts/releases/download/gpu-operator-%[1]s/gpu-operator-%[1]s.tgz", versions.GPUOperator),
					Version: versions.GPUOperator,
					Labels: map[string]string{
						release.TypeLabel: "extension",
					},
				},
			},
		},
		{
			Name:        "llm-d-infra",
			Namespace:   "llm-d",
			DisplayName: "llm-d",
			Description: "Achieve state of the art inference performance with modern accelerators on Kubernetes.",
			Logo:        "https://github.com/llm-d.png",
			Charts: []chartManifest{
				{
					Ref:     fmt.Sprintf("https://github.com/llm-d-incubation/llm-d-infra/releases/download/v%[1]s/llm-d-infra-v%[1]s.tgz", versions.LLMDInfra),
					Version: versions.LLMDInfra,
					Labels: map[string]string{
						release.TypeLabel: "extension",
					},
					ValuesMap: map[string]string{
						"gateway.gatewayParameters.resources.limits.cpu":    "2",
						"gateway.gatewayParameters.resources.limits.memory": "1Gi",
						"gateway.service.type":                              "NodePort",
					},
				},
			},
		},
	}

	instance = []base{
		{
			Name:        "kubevirt-infra",
			Namespace:   "kubevirt",
			DisplayName: "KubeVirt",
			Description: "Kubernetes Virtualization API and runtime in order to define and manage virtual machines.",
			Logo:        "https://github.com/kubevirt.png",
			Charts: []chartManifest{
				{
					Ref:     fmt.Sprintf("https://github.com/otterscale/charts/releases/download/kubevirt-infra-%[1]s/kubevirt-infra-%[1]s.tgz", versions.KubeVirtInfra),
					Version: versions.KubeVirtInfra,
					Labels: map[string]string{
						release.TypeLabel: "extension",
					},
				},
			},
		},
	}

	storage = []base{
		{
			Name:        "samba-operator",
			Namespace:   "samba-operator",
			DisplayName: "Samba",
			Description: "An operator for Samba as a service on PVCs in Kubernetes.",
			Logo:        "https://github.com/otterscale.png",
			Charts: []chartManifest{
				{
					Ref:     fmt.Sprintf("https://github.com/otterscale/charts/releases/download/samba-operator-%[1]s/samba-operator-%[1]s.tgz", versions.SambaOperator),
					Version: versions.SambaOperator,
					Labels: map[string]string{
						release.TypeLabel: "extension",
					},
				},
			},
		},
	}
)

func (uc *UseCase) base(name string) (base, error) {
	all := []base{}
	all = append(all, general...)
	all = append(all, registry...)
	all = append(all, model...)
	all = append(all, instance...)
	all = append(all, storage...)

	for _, p := range all {
		if p.Name == name {
			return p, nil
		}
	}

	return base{}, fmt.Errorf("extension %s not found", name)
}
