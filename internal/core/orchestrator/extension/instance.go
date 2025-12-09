package extension

import (
	"context"
	"fmt"

	"github.com/otterscale/otterscale/internal/core/versions"
)

var kubevirtInfraChartRef = fmt.Sprintf("https://github.com/otterscale/charts/releases/download/kubevirt-infra-%[1]s/kubevirt-infra-%[1]s.tgz", versions.KubeVirtInfra)

var instanceComponents = []component{
	{
		ID:          "kubevirt-infra",
		DisplayName: "KubeVirt",
		Description: "Kubernetes Virtualization API and runtime in order to define and manage virtual machines.",
		Logo:        "https://github.com/kubevirt.png",
		Chart: &chartComponent{
			Name:      "kubevirt-infra",
			Namespace: "kubevirt",
			Ref:       kubevirtInfraChartRef,
			Version:   versions.KubeVirtInfra,
			ValuesMap: map[string]string{
				"kubevirt.serviceMonitor.enabled": "true",
			},
		},
		Dependencies: []string{"kube-prometheus-stack"},
		PostFunc: func(uc *UseCase, ctx context.Context, scope string) error {
			return uc.patchContainerdTemplates(ctx, scope, false)
		},
	},
}
