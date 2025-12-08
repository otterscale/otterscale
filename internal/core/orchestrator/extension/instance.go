package extension

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
			units, err := uc.getContainerdUnits(ctx, scope)
			if err != nil {
				return err
			}

			for _, unit := range units {
				ok, err := uc.hasDeviceOwnershipConfig(ctx, scope, unit)
				if err != nil {
					return err
				}

				if ok {
					continue
				}

				if err := uc.setDeviceOwnershipConfig(ctx, scope, unit); err != nil {
					return err
				}

				if err := uc.reconcileContainerd(ctx, scope); err != nil {
					return err
				}
			}
			return nil
		},
	},
}

func (uc *UseCase) getContainerdUnits(ctx context.Context, scope string) ([]string, error) {
	facilities, err := uc.facility.List(ctx, scope, "")
	if err != nil {
		return nil, err
	}

	units := []string{}

	for i := range facilities {
		status := facilities[i].Status

		if status == nil || !strings.Contains(status.Charm, "kubernetes") {
			continue
		}

		for name := range status.Units {
			unit := status.Units[name]

			for subName := range unit.Subordinates {
				if !strings.Contains(subName, "containerd") {
					continue
				}

				units = append(units, subName)
			}
		}
	}

	return units, nil
}

func (uc *UseCase) hasDeviceOwnershipConfig(ctx context.Context, scope, unitName string) (bool, error) {
	path := uc.containerdConfigPath(unitName)
	cmd := fmt.Sprintf("cat %s", path)

	r, err := uc.action.Execute(ctx, scope, unitName, cmd)
	if err != nil {
		return false, err
	}

	stdout, ok := r["stdout"].(string)
	if !ok {
		return false, errors.New("containerd config stdout not found")
	}

	return strings.Contains(stdout, "device_ownership_from_security_context = true"), nil
}

func (uc *UseCase) setDeviceOwnershipConfig(ctx context.Context, scope, unitName string) error {
	path := uc.containerdConfigPath(unitName)
	cmd := fmt.Sprintf(`sed -i '/\[plugins."io.containerd.grpc.v1.cri"\]/a \ \ \ \ device_ownership_from_security_context = true' %s`, path)

	_, err := uc.action.Execute(ctx, scope, unitName, cmd)
	return err
}

func (uc *UseCase) reconcileContainerd(ctx context.Context, scope string) error {
	if err := uc.setContainerdCustomRegistries(ctx, scope, true); err != nil {
		return err
	}
	return uc.setContainerdCustomRegistries(ctx, scope, false)
}

func (uc *UseCase) containerdConfigPath(unitName string) string {
	return fmt.Sprintf("/var/lib/juju/agents/unit-%s/charm/templates/config_v2.toml", strings.ReplaceAll(unitName, "/", "-"))
}
