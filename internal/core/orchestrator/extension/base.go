package extension

import (
	"fmt"

	"github.com/otterscale/otterscale/internal/core/application/release"
)

const chartRepoURL = "https://otterscale.github.io/charts"

type base struct {
	name        string
	namespace   string
	version     string
	repoURL     string
	labels      map[string]string
	annotations map[string]string
	valuesMap   map[string]string
}

var (
	general = []base{
		{
			name:      "kube-prometheus-stack",
			namespace: "monitoring",
			repoURL:   "https://prometheus-community.github.io/helm-charts",
			valuesMap: map[string]string{
				release.TypeLabel: "extension",
			},
		},
	}

	model = []base{
		{
			name:      "gpu-operator",
			namespace: "nvidia-gpu-operator",
			repoURL:   "https://nvidia.github.io/gpu-operator",
			valuesMap: map[string]string{
				release.TypeLabel: "extension",
			},
		},
		{
			name:      "llm-d-infra",
			namespace: "llm-d",
			repoURL:   "https://llm-d-incubation.github.io/llm-d-infra",
			valuesMap: map[string]string{
				release.TypeLabel: "extension",
				"nameOverride":    "llm-d-infra",
				"gateway.gatewayParameters.resources.limits.cpu":    "4",
				"gateway.gatewayParameters.resources.limits.memory": "2Gi",
				"gateway.service.type":                              "NodePort",
			},
		},
	}

	instance = []base{
		{
			name:      "kubevirt-infra",
			namespace: "kubevirt",
			repoURL:   chartRepoURL,
			valuesMap: map[string]string{
				release.TypeLabel: "extension",
			},
		},
	}

	storage = []base{
		{
			name:      "samba-operator",
			namespace: "samba-operator",
			repoURL:   chartRepoURL,
			valuesMap: map[string]string{
				release.TypeLabel: "extension",
			},
		},
	}
)

func (uc *ExtensionUseCase) base(name string) (base, error) {
	all := []base{}
	all = append(all, general...)
	all = append(all, model...)
	all = append(all, instance...)
	all = append(all, storage...)

	for _, p := range all {
		if p.name == name {
			return p, nil
		}
	}

	return base{}, fmt.Errorf("extension %s not found", name)
}
