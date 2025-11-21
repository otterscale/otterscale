package extension

import (
	"fmt"

	"github.com/otterscale/otterscale/internal/core/application/release"
)

const chartRepoURL = "https://otterscale.github.io/charts"

type chartManifest struct {
	ID          string // If a value is present, it indicates a child chart.
	Namespace   string
	RepoURL     string
	Labels      map[string]string
	Annotations map[string]string
	ValuesMap   map[string]string
	PostFunc    func(scope string) error
}

type crdManifest struct {
	Version              string
	RepoURL              string
	AnnotationVersionKey string
}

type base struct {
	ID          string
	Name        string
	Description string
	Logo        string
	Charts      []chartManifest
	CRD         *crdManifest
}

var (
	general = []base{
		{
			ID:          "gateway-api",
			Name:        "Gateway API",
			Description: "Gateway API is an official Kubernetes project focused on L4 and L7 routing in Kubernetes.",
			Logo:        "https://github.com/kubernetes-sigs.png",
			CRD: &crdManifest{
				Version:              "v1.3.0",
				RepoURL:              "https://github.com/kubernetes-sigs/gateway-api.git/config/crd",
				AnnotationVersionKey: "gateway.networking.k8s.io/bundle-version",
			},
		},
		{
			ID:          "gateway-api-inference-extension",
			Name:        "Gateway API Inference Extension",
			Description: "Gateway API Inference Extension is an official Kubernetes project that optimizes self-hosting Generative Models on Kubernetes.",
			Logo:        "https://github.com/kubernetes-sigs.png",
			CRD: &crdManifest{
				Version:              "v1.1.0",
				RepoURL:              "https://github.com/kubernetes-sigs/gateway-api-inference-extension.git/config/crd",
				AnnotationVersionKey: "inference.networking.k8s.io/bundle-version",
			},
		},
		{
			ID:          "base",
			Name:        "Istio",
			Description: "Istio is an open source service mesh that layers transparently onto existing distributed applications.",
			Logo:        "https://github.com/istio.png",
			Charts: []chartManifest{
				{
					Namespace: "istio-system",
					RepoURL:   "https://istio-release.storage.googleapis.com/charts",
					Labels: map[string]string{
						release.TypeLabel: "extension",
					},
				},
				{
					ID:        "istiod",
					Namespace: "istio-system",
					RepoURL:   "https://istio-release.storage.googleapis.com/charts",
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
			ID:          "kube-prometheus-stack",
			Name:        "Prometheus Stack",
			Description: "Installs the kube-prometheus stack for easy, end-to-end Kubernetes cluster monitoring using the Prometheus Operator.",
			Logo:        "https://github.com/prometheus-community.png",
			Charts: []chartManifest{
				{
					Namespace: "monitoring",
					RepoURL:   "https://prometheus-community.github.io/helm-charts",
					Labels: map[string]string{
						release.TypeLabel: "extension",
					},
				},
			},
		},
	}

	registry = []base{
		{
			ID:          "kubevirt-infra", // helm chart name
			Name:        "Registry",
			Description: "",
			Logo:        "https://github.com/distribution.png",
			Charts: []chartManifest{
				{
					Namespace: "distribution",
					RepoURL:   chartRepoURL,
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
			ID:          "gpu-operator",
			Name:        "GPU Operator",
			Description: "GPU Operator creates, configures, and manages GPUs in Kubernetes.",
			Logo:        "https://github.com/otterscale.png",
			Charts: []chartManifest{
				{
					Namespace: "gpu-operator",
					RepoURL:   chartRepoURL,
					Labels: map[string]string{
						release.TypeLabel: "extension",
					},
				},
			},
		},
		{
			ID:          "llm-d-infra",
			Name:        "llm-d",
			Description: "Achieve state of the art inference performance with modern accelerators on Kubernetes.",
			Logo:        "https://github.com/llm-d.png",
			Charts: []chartManifest{
				{
					Namespace: "llm-d",
					RepoURL:   "https://llm-d-incubation.github.io/llm-d-infra",
					Labels: map[string]string{
						release.TypeLabel: "extension",
					},
					ValuesMap: map[string]string{
						"nameOverride": "llm-d-gateway",
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
			ID:          "kubevirt-infra",
			Name:        "KubeVirt",
			Description: "Kubernetes Virtualization API and runtime in order to define and manage virtual machines.",
			Logo:        "https://github.com/kubevirt.png",
			Charts: []chartManifest{
				{
					Namespace: "kubevirt",
					RepoURL:   chartRepoURL,
					Labels: map[string]string{
						release.TypeLabel: "extension",
					},
				},
			},
		},
	}

	storage = []base{
		{
			ID:          "samba-operator",
			Name:        "Samba",
			Description: "An operator for Samba as a service on PVCs in Kubernetes.",
			Logo:        "https://github.com/otterscale.png",
			Charts: []chartManifest{
				{
					Namespace: "samba-operator",
					RepoURL:   chartRepoURL,
					Labels: map[string]string{
						release.TypeLabel: "extension",
					},
				},
			},
		},
	}
)

func (uc *UseCase) base(id string) (base, error) {
	all := []base{}
	all = append(all, general...)
	all = append(all, registry...)
	all = append(all, model...)
	all = append(all, instance...)
	all = append(all, storage...)

	for _, p := range all {
		if p.ID == id {
			return p, nil
		}
	}

	return base{}, fmt.Errorf("extension %s not found", id)
}
