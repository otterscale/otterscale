package extension

import (
	"context"
	"fmt"
	"slices"
	"strings"

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
				"prometheus.service.type":                                                                   "NodePort",
				"prometheus.service.nodePort":                                                               "0",
				"nodeExporter.enabled":                                                                      "true",
				"grafana.enabled":                                                                           "false",
			},
		},
		PostFunc: func(uc *UseCase, ctx context.Context, scope string) error {
			cosScope := "cos"
			appName := "prometheus-scrape-target-k8s"
			namespace := "monitoring"
			service := "kube-prometheus-stack-prometheus"
			portName := "http-web"

			targets, err := uc.getPrometheusScrapeTargetK8sTargets(ctx, cosScope, appName)
			if err != nil {
				return err
			}

			target, err := uc.getPrometheusTarget(ctx, scope, namespace, service, portName)
			if err != nil {
				return err
			}

			targets = append(targets, target)

			slices.Sort(targets)
			targets = slices.Compact(targets)

			return uc.setPrometheusScrapeTargetK8sTargets(ctx, cosScope, appName, targets)
		},
	},
}

func (uc *UseCase) getPrometheusScrapeTargetK8sTargets(ctx context.Context, scope, name string) ([]string, error) {
	config, err := uc.facility.Config(ctx, scope, name)
	if err != nil {
		return nil, err
	}

	targets, ok := config["targets"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid type for targets field")
	}

	value, ok := targets["value"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid type for targets.value field")
	}

	return strings.Split(value, ","), nil
}

func (uc *UseCase) getPrometheusTarget(ctx context.Context, scope, namespace, name, portName string) (string, error) {
	ip, err := uc.node.InternalIP(ctx, scope)
	if err != nil {
		return "", err
	}

	svc, err := uc.service.Get(ctx, scope, namespace, name)
	if err != nil {
		return "", err
	}

	ports := svc.Spec.Ports

	for i := range ports {
		if ports[i].Name == portName {
			return fmt.Sprintf("%s:%d", ip, ports[i].NodePort), nil
		}
	}

	return "", fmt.Errorf("prometheus service has no %s port defined", portName)
}

func (uc *UseCase) setPrometheusScrapeTargetK8sTargets(ctx context.Context, scope string, name string, targets []string) error {
	return uc.facility.Update(ctx, scope, name, map[string]string{"targets": strings.Join(targets, ",")})
}
