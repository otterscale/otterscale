package extension

import (
	"fmt"

	"github.com/otterscale/otterscale/internal/core/versions"
)

var kubePrometheusStackChartRef = fmt.Sprintf("https://github.com/prometheus-community/helm-charts/releases/download/kube-prometheus-stack-%[1]s/kube-prometheus-stack-%[1]s.tgz", versions.KubePrometheusStack)

var metricsComponents = []component{
	{
		ID:          "kube-prometheus-stack",
		DisplayName: "Prometheus Stack",
		Description: "Installs the kube-prometheus stack for easy, end-to-end Kubernetes cluster monitoring using the Prometheus Operator.",
		Logo:        "https://github.com/prometheus-community.png",
		Chart: &chartComponent{
			Name:      "kube-prometheus-stack",
			Namespace: "monitoring",
			Ref:       kubePrometheusStackChartRef,
			Version:   versions.KubePrometheusStack,
			ValuesMap: map[string]string{
				"prometheus.prometheusSpec.externalLabels.juju_model":                                       "{{ .Scope }}",
				"prometheus.prometheusSpec.externalLabels.juju_model_uuid":                                  "{{ .Scope.UUID }}",
				"prometheus.prometheusSpec.enableRemoteWriteReceiver":                                       "true",
				"prometheus.prometheusSpec.retention":                                                       "365d",
				"prometheus.prometheusSpec.retentionSize":                                                   "40GiB",
				"prometheus.prometheusSpec.storageSpec.volumeClaimTemplate.spec.storageClassName":           "ceph-ext4",
				"prometheus.prometheusSpec.storageSpec.volumeClaimTemplate.spec.accessModes[0]":             "ReadWriteOnce",
				"prometheus.prometheusSpec.storageSpec.volumeClaimTemplate.spec.resources.requests.storage": "40Gi",
				"prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues":                         "false",
				"nodeExporter.enabled": "true",
				"grafana.enabled":      "false",
			},
		},
	},
}
