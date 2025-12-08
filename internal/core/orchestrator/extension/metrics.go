package extension

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"

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
			targets, err := uc.getPrometheusScrapeTargetK8sTargets(ctx)
			if err != nil {
				return err
			}

			target, err := uc.waitAndGetPrometheusTarget(ctx, scope)
			if err != nil {
				return err
			}

			targets = append(targets, target)

			slices.Sort(targets)
			targets = slices.Compact(targets)

			return uc.setPrometheusScrapeTargetK8sTargets(ctx, targets)
		},
	},
}

func (uc *UseCase) getPrometheusScrapeTargetK8sTargets(ctx context.Context) ([]string, error) {
	config, err := uc.facility.Config(ctx, "cos", "prometheus-scrape-target-k8s")
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

	return slices.DeleteFunc(strings.Split(value, ","), func(target string) bool {
		return target == ""
	}), nil
}

func (uc *UseCase) getPrometheusTargetNodePort(ctx context.Context, scope string) (int32, error) {
	svc, err := uc.service.Get(ctx, scope, "monitoring", "kube-prometheus-stack-prometheus")
	if err != nil {
		return 0, err
	}

	ports := svc.Spec.Ports

	for i := range ports {
		if ports[i].Name == "http-web" {
			return ports[i].NodePort, nil
		}
	}

	return 0, errors.New("prometheus service has no http-web port defined")
}

func (uc *UseCase) setPrometheusScrapeTargetK8sTargets(ctx context.Context, targets []string) error {
	return uc.facility.Update(ctx, "cos", "prometheus-scrape-target-k8s", map[string]string{"targets": strings.Join(targets, ",")})
}

func (uc *UseCase) waitAndGetPrometheusTarget(ctx context.Context, scope string) (string, error) {
	const (
		timeout  = 5 * time.Minute
		interval = 5 * time.Second
	)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	ip, err := uc.node.InternalIP(ctx, scope)
	if err != nil {
		return "", err
	}

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()

		case <-ticker.C:
			port, err := uc.getPrometheusTargetNodePort(ctx, scope)
			if k8serrors.IsNotFound(err) {
				continue
			}
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("%s:%d", ip, port), nil
		}
	}
}
