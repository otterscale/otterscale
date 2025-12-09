package extension

import (
	"context"
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
			return uc.setPrometheusTargets(ctx, scope)
		},
	},
}

func (uc *UseCase) setPrometheusTargets(ctx context.Context, scope string) error {
	targets, err := uc.getPrometheusTargets(ctx)
	if err != nil {
		return err
	}

	newTarget, err := uc.waitPrometheusNewTarget(ctx, scope)
	if err != nil {
		return err
	}

	targets = append(targets, newTarget)

	// de-duplicate
	slices.Sort(targets)
	targets = slices.Compact(targets)

	return uc.setConfig(ctx, "cos", "prometheus-scrape-target-k8s", "targets", strings.Join(targets, ","))
}

func (uc *UseCase) getPrometheusTargets(ctx context.Context) ([]string, error) {
	config, err := uc.getConfig(ctx, "cos", "prometheus-scrape-target-k8s")
	if err != nil {
		return nil, err
	}

	targets, err := getValue[string](config, "targets")
	if err != nil {
		return nil, err
	}

	return slices.DeleteFunc(strings.Split(targets, ","), func(target string) bool {
		return target == ""
	}), nil
}

func (uc *UseCase) waitPrometheusNewTarget(ctx context.Context, scope string) (string, error) {
	const interval = 5 * time.Second

	ip, err := uc.node.InternalIP(ctx, scope)
	if err != nil {
		return "", err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()

		case <-ticker.C:
			port, err := uc.getNodePort(ctx, scope, "monitoring", "kube-prometheus-stack-prometheus", "http-web")
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
