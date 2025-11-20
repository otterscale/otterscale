package extension

import (
	"fmt"

	"github.com/otterscale/otterscale/internal/core/application/release"
)

const chartRepoURL = "https://otterscale.github.io/charts"

type chartManifest struct {
	Namespace   string
	RepoURL     string
	Labels      map[string]string
	Annotations map[string]string
	ValuesMap   map[string]string
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
	Chart       *chartManifest
	CRD         *crdManifest
}

var (
	general = []base{
		{
			ID:          "gateway-api",
			Name:        "Gateway API",
			Description: "Gateway API is an official Kubernetes project focused on L4 and L7 routing in Kubernetes.",
			Logo:        "https://raw.githubusercontent.com/kubernetes-sigs/gateway-api/refs/tags/v1.3.0/site-src/images/logo/logo.svg",
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
			Logo:        "",
			CRD: &crdManifest{
				Version:              "v1.1.0",
				RepoURL:              "https://github.com/kubernetes-sigs/gateway-api-inference-extension.git/config/crd",
				AnnotationVersionKey: "inference.networking.k8s.io/bundle-version",
			},
		},
		{
			ID:          "kube-prometheus-stack",
			Name:        "Prometheus Stack",
			Description: "Installs the kube-prometheus stack for easy, end-to-end Kubernetes cluster monitoring using the Prometheus Operator.",
			Logo:        "https://github.com/prometheus-community.png",
			Chart: &chartManifest{
				Namespace: "monitoring",
				RepoURL:   "https://prometheus-community.github.io/helm-charts",
				ValuesMap: map[string]string{
					release.TypeLabel: "extension",
				},
			},
		},
	}

	model = []base{
		{
			ID:          "gpu-operator",
			Name:        "NVIDIA GPU Operator",
			Description: "NVIDIA GPU Operator creates, configures, and manages GPUs in Kubernetes.",
			Logo:        "https://github.com/nvidia.png",
			Chart: &chartManifest{
				Namespace: "nvidia-gpu-operator",
				RepoURL:   "https://nvidia.github.io/gpu-operator",
				ValuesMap: map[string]string{
					release.TypeLabel: "extension",
				},
			},
		},
		{
			ID:          "llm-d-infra",
			Name:        "llm-d",
			Description: "Achieve state of the art inference performance with modern accelerators on Kubernetes.",
			Logo:        "https://github.com/llm-d.png",
			Chart: &chartManifest{
				Namespace: "llm-d",
				RepoURL:   "https://llm-d-incubation.github.io/llm-d-infra",
				ValuesMap: map[string]string{
					release.TypeLabel: "extension",
					"nameOverride":    "llm-d-infra",
					"gateway.gatewayParameters.resources.limits.cpu":    "4",
					"gateway.gatewayParameters.resources.limits.memory": "2Gi",
					"gateway.service.type":                              "NodePort",
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
			Chart: &chartManifest{
				Namespace: "kubevirt",
				RepoURL:   chartRepoURL,
				ValuesMap: map[string]string{
					release.TypeLabel: "extension",
				},
			},
		},
	}

	storage = []base{
		{
			ID:          "samba-operator",
			Name:        "Samba",
			Description: "An operator for Samba as a service on PVCs in Kubernetes.",
			Logo:        "",
			Chart: &chartManifest{
				Namespace: "samba-operator",
				RepoURL:   chartRepoURL,
				ValuesMap: map[string]string{
					release.TypeLabel: "extension",
				},
			},
		},
	}
)

func (uc *UseCase) base(id string) (base, error) {
	all := []base{}
	all = append(all, general...)
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
