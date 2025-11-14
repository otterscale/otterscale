package extension

import (
	"fmt"

	"github.com/otterscale/otterscale/internal/core/application/release"
)

const chartRepoURL = "https://otterscale.github.io/charts"

type base struct {
	Name        string
	Namespace   string
	Version     string
	RepoURL     string
	Labels      map[string]string
	Annotations map[string]string
	ValuesMap   map[string]string
}

var (
	general = []base{
		{
			Name:      "kube-prometheus-stack",
			Namespace: "monitoring",
			RepoURL:   "https://prometheus-community.github.io/helm-charts",
			ValuesMap: map[string]string{
				release.TypeLabel: "extension",
			},
		},
	}

	model = []base{
		{
			Name:      "gpu-operator",
			Namespace: "nvidia-gpu-operator",
			RepoURL:   "https://nvidia.github.io/gpu-operator",
			ValuesMap: map[string]string{
				release.TypeLabel: "extension",
			},
		},
		{
			Name:      "llm-d-infra",
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
	}

	instance = []base{
		{
			Name:      "kubevirt-infra",
			Namespace: "kubevirt",
			RepoURL:   chartRepoURL,
			ValuesMap: map[string]string{
				release.TypeLabel: "extension",
			},
		},
	}

	storage = []base{
		{
			Name:      "samba-operator",
			Namespace: "samba-operator",
			RepoURL:   chartRepoURL,
			ValuesMap: map[string]string{
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
		if p.Name == name {
			return p, nil
		}
	}

	return base{}, fmt.Errorf("extension %s not found", name)
}
