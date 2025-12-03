package extension

import (
	"fmt"

	"github.com/otterscale/otterscale/internal/core/versions"
)

var (
	gaCRDRef          = fmt.Sprintf("https://github.com/kubernetes-sigs/gateway-api.git/config/crd?ref=v%s", versions.GatewayAPI)
	gaieCRDRef        = fmt.Sprintf("https://github.com/kubernetes-sigs/gateway-api-inference-extension.git/config/crd?ref=v%s", versions.GatewayAPIInferenceExtension)
	istioBaseChartRef = fmt.Sprintf("https://istio-release.storage.googleapis.com/charts/base-%s.tgz", versions.Istio)
	istiodChartRef    = fmt.Sprintf("https://istio-release.storage.googleapis.com/charts/istiod-%s.tgz", versions.Istio)
)

var serviceMeshComponents = []component{
	{
		ID:          "gateway-api",
		DisplayName: "Gateway API",
		Description: "Gateway API is an official Kubernetes project focused on L4 and L7 routing in Kubernetes.",
		Logo:        "https://github.com/kubernetes-sigs.png",
		CRD: &crdComponent{
			Ref:               gaCRDRef,
			Version:           versions.GatewayAPI,
			VersionAnnotation: "gateway.networking.k8s.io/bundle-version",
		},
	},
	{
		ID:          "gateway-api-inference-extension",
		DisplayName: "Gateway API Inference Extension",
		Description: "Gateway API Inference Extension is an official Kubernetes project that optimizes self-hosting Generative Models on Kubernetes.",
		Logo:        "https://github.com/kubernetes-sigs.png",
		CRD: &crdComponent{
			Ref:               gaieCRDRef,
			Version:           versions.GatewayAPIInferenceExtension,
			VersionAnnotation: "inference.networking.k8s.io/bundle-version",
		},
	},
	{
		ID:          "istio-base",
		DisplayName: "Istio",
		Description: "Istio is an open source service mesh that layers transparently onto existing distributed applications.",
		Logo:        "https://github.com/istio.png",
		Chart: &chartComponent{
			Name:      "istio-base",
			Namespace: "istio-system",
			Ref:       istioBaseChartRef,
			Version:   versions.Istio,
		},
	},
	{
		ID:          "istiod",
		DisplayName: "Istio",
		Description: "Istio is an open source service mesh that layers transparently onto existing distributed applications.",
		Logo:        "https://github.com/istio.png",
		Chart: &chartComponent{
			Name:      "istiod",
			Namespace: "istio-system",
			Ref:       istiodChartRef,
			Version:   versions.Istio,
			ValuesMap: map[string]string{
				"env.ENABLE_GATEWAY_API_INFERENCE_EXTENSION": "true",
			},
		},
	},
}
