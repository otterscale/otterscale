package extension

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/otterscale/otterscale/internal/core/versions"
)

var registryChartRef = fmt.Sprintf("https://github.com/otterscale/charts/releases/download/registry-%[1]s/registry-%[1]s.tgz", versions.Registry)

var registryComponents = []component{
	{
		ID:          "registry",
		DisplayName: "Registry",
		Description: "OCI compliant container image registry for storing and distributing container images.",
		Logo:        "https://github.com/distribution.png",
		Chart: &chartComponent{
			Name:      "registry",
			Namespace: "registry",
			Ref:       registryChartRef,
			Version:   versions.Registry,
			ValuesMap: map[string]string{
				"dockerRegistry.image.repository":          "registry",
				"dockerRegistry.image.tag":                 "3",
				"dockerRegistry.replicaCount":              "1",
				"dockerRegistry.service.type":              "NodePort",
				"dockerRegistry.service.nodePort":          "0",
				"dockerRegistry.persistence.size":          "600Gi",
				"dockerRegistry.persistence.storageClass":  "ceph-ext4",
				"dockerRegistry.resources.requests.memory": "256Mi",
				"dockerRegistry.resources.requests.cpu":    "250m",
				"dockerRegistry.resources.limits.memory":   "512Mi",
				"dockerRegistry.resources.limits.cpu":      "500m",
			},
		},
		Dependencies: []string{"kube-prometheus-stack"},
		PostFunc: func(uc *UseCase, ctx context.Context, scope string) error {
			return uc.setContainerdCustomRegistries(ctx, scope, false)
		},
	},
}

type registryConfig struct {
	URL                string `json:"url"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify"`
}

func (uc *UseCase) setContainerdCustomRegistries(ctx context.Context, scope string, hack bool) error {
	url, err := uc.repository.GetRegistryURL(scope)
	if err != nil {
		return err
	}

	registries := []registryConfig{
		{
			URL:                "http://" + url,
			InsecureSkipVerify: true,
		},
	}

	if hack {
		registries = append(registries, registryConfig{
			URL:                "http://hack.reconcile",
			InsecureSkipVerify: true,
		})
	}

	value, err := json.Marshal(registries)
	if err != nil {
		return err
	}

	return uc.facility.Update(ctx, scope, scope+"-containerd", map[string]string{"custom_registries": string(value)})
}
