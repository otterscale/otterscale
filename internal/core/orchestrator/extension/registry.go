package extension

import (
	"context"
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
				"dockerRegistry.persistence.size":          "600Gi",
				"dockerRegistry.persistence.storageClass":  "ceph-ext4",
				"dockerRegistry.resources.requests.memory": "256Mi",
				"dockerRegistry.resources.requests.cpu":    "250m",
				"dockerRegistry.resources.limits.memory":   "512Mi",
				"dockerRegistry.resources.limits.cpu":      "500m",
			},
		},
		Dependencies: []string{"kube-prometheus-stack"},
		PostFunc: func(ctx context.Context, scope string) error {
			return nil
		},
	},
}
