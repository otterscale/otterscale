package extension

import (
	"fmt"

	"github.com/otterscale/otterscale/internal/core/versions"
)

var (
	modelChartRef     = fmt.Sprintf("https://github.com/otterscale/charts/releases/download/gpu-operator-%[1]s/gpu-operator-%[1]s.tgz", versions.GPUOperator)
	llmdInfraChartRef = fmt.Sprintf("https://github.com/llm-d-incubation/llm-d-infra/releases/download/v%[1]s/llm-d-infra-v%[1]s.tgz", versions.LLMDInfra)
)

var modelComponents = []component{
	{
		ID:          "gpu-operator",
		DisplayName: "GPU Operator",
		Description: "GPU Operator creates, configures, and manages GPUs in Kubernetes.",
		Logo:        "https://github.com/otterscale.png",
		Chart: &chartComponent{
			Name:      "gpu-operator",
			Namespace: "gpu-operator",
			Ref:       modelChartRef,
			Version:   versions.GPUOperator,
			ValuesMap: map[string]string{
				"gpu-operator.driver.version":                      "580.95.05",
				"gpu-operator.driver.upgradePolicy.autoUpgrade":    "false",
				"gpu-operator.dcgmExporter.serviceMonitor.enabled": "true",
				"hami.prometheus.serviceMonitor.enabled":           "true",
			},
		},
		Dependencies: []string{"kube-prometheus-stack"},
	},
	{
		ID:          "llm-d-infra",
		DisplayName: "llm-d",
		Description: "Achieve state of the art inference performance with modern accelerators on Kubernetes.",
		Logo:        "https://github.com/llm-d.png",
		Chart: &chartComponent{
			Name:      "llm-d-infra",
			Namespace: "llm-d",
			Ref:       llmdInfraChartRef,
			Version:   versions.LLMDInfra,
			ValuesMap: map[string]string{
				"gateway.gatewayParameters.resources.limits.cpu":    "2",
				"gateway.gatewayParameters.resources.limits.memory": "1Gi",
				"gateway.service.type":                              "NodePort",
			},
		},
	},
}
