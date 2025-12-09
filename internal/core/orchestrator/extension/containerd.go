package extension

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/otterscale/otterscale/internal/core/machine"
)

type unitInfo struct {
	Hostname string
	HasGPU   bool
}

func (uc *UseCase) containerdConfigPath() string {
	return "/etc/containerd/config.toml"
}

func (uc *UseCase) getContainerdConfig(ctx context.Context, scope, unit string) (string, error) {
	path := uc.containerdConfigPath()
	cmd := fmt.Sprintf("cat %s", path)

	r, err := uc.action.Execute(ctx, scope, unit, cmd)
	if err != nil {
		return "", err
	}

	stdout, ok := r["stdout"].(string)
	if !ok {
		return "", errors.New("containerd config stdout not found")
	}

	return stdout, nil
}

func (uc *UseCase) containerdTemplatePath(unitName string) string {
	sanitizedName := strings.ReplaceAll(unitName, "/", "-")
	return fmt.Sprintf("/var/lib/juju/agents/unit-%s/charm/templates/config_v2.toml", sanitizedName)
}

func (uc *UseCase) setContainerdTemplate(ctx context.Context, scope, unit string, hasGPU bool) error {
	data := base64.StdEncoding.EncodeToString([]byte(renderContainerdTemplate(hasGPU)))
	path := uc.containerdTemplatePath(unit)

	cmd := fmt.Sprintf(`echo %q | base64 -d > %s`, data, path)

	_, err := uc.action.Execute(ctx, scope, unit, cmd)
	return err
}

func (uc *UseCase) getCustomRegistries(ctx context.Context, scope string) ([]registryConfig, error) {
	config, err := uc.getConfig(ctx, scope, "containerd")
	if err != nil {
		return nil, err
	}

	customRegistries, err := getValue[string](config, "custom_registries")
	if err != nil {
		return nil, err
	}

	var registries []registryConfig

	if err := json.Unmarshal([]byte(customRegistries), &registries); err != nil {
		return nil, err
	}

	return registries, nil
}

func (uc *UseCase) reconcileContainerd(ctx context.Context, scope, unit string) error {
	registries, err := uc.getCustomRegistries(ctx, scope)
	if err != nil {
		return err
	}

	if err := uc.hackAndWaitContainerdConfigChanged(ctx, scope, unit, registries); err != nil {
		return err
	}

	// set custom registries back
	return uc.setContainerdCustomRegistries(ctx, scope, registries)
}

func (uc *UseCase) patchContainerdTemplates(ctx context.Context, scope string, wait bool) error {
	ctrUnitInfo, err := uc.getContainerdUnitInfo(ctx, scope)
	if err != nil {
		return err
	}

	if wait {
		if !uc.anyGPU(ctrUnitInfo) {
			return errors.New("no GPU found in any unit")
		}

		if err := uc.waitGPUOperatorReady(ctx, scope); err != nil {
			return err
		}
	}

	eg, egctx := errgroup.WithContext(ctx)

	for ctr, unitInfo := range ctrUnitInfo {
		eg.Go(func() error {
			hasGPU, err := uc.nodeContainsGPU(egctx, scope, unitInfo.Hostname)
			if err != nil {
				return err
			}

			if err := uc.setContainerdTemplate(egctx, scope, ctr, hasGPU); err != nil {
				return err
			}

			return uc.reconcileContainerd(egctx, scope, ctr)
		})
	}

	return eg.Wait()
}

func (uc *UseCase) nodeContainsGPU(ctx context.Context, scope, name string) (bool, error) {
	node, err := uc.node.Get(ctx, scope, name)
	if err != nil {
		return false, err
	}

	present, ok := node.GetLabels()["nvidia.com/gpu.deploy.operator-validator"]
	if ok {
		return present == "true", nil
	}

	return false, nil
}

func (uc *UseCase) getContainerdUnitInfo(ctx context.Context, scope string) (map[string]unitInfo, error) {
	machineUnitInfo, err := uc.findMachineUnitInfo(ctx, scope)
	if err != nil {
		return nil, err
	}

	return uc.findContainerdUnitInfo(ctx, scope, machineUnitInfo)
}

func (uc *UseCase) findContainerdUnitInfo(ctx context.Context, scope string, machineUnitInfo map[string]unitInfo) (map[string]unitInfo, error) {
	facilities, err := uc.facility.List(ctx, scope, "")
	if err != nil {
		return nil, err
	}

	ctrUnitInfo := map[string]unitInfo{}

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

				if unitInfo, ok := machineUnitInfo[unit.Machine]; ok {
					ctrUnitInfo[subName] = unitInfo
				}
			}
		}
	}

	return ctrUnitInfo, nil
}

func (uc *UseCase) findMachineUnitInfo(ctx context.Context, scope string) (map[string]unitInfo, error) {
	machines, err := uc.machine.List(ctx)
	if err != nil {
		return nil, err
	}

	machineUnitInfo := map[string]unitInfo{}

	for i := range machines {
		machineScope, err := uc.machine.ExtractScope(&machines[i])
		if err != nil {
			continue
		}

		if machineScope != scope {
			continue
		}

		jujuID, err := uc.machine.ExtractJujuID(&machines[i])
		if err != nil {
			continue
		}

		gpus, err := uc.nodeDevice.ListGPUs(ctx, machines[i].SystemID)
		if err != nil {
			return nil, err
		}

		machineUnitInfo[jujuID] = unitInfo{
			Hostname: machines[i].Hostname,
			HasGPU:   uc.containsNvidiaGPU(gpus),
		}
	}

	return machineUnitInfo, nil
}

func (uc *UseCase) containsNvidiaGPU(gpus []machine.GPU) bool {
	for i := range gpus {
		if strings.EqualFold(gpus[i].VendorID, "10de") {
			return true
		}
	}
	return false
}

func (uc *UseCase) anyGPU(ctrUnitInfo map[string]unitInfo) bool {
	for _, unitInfo := range ctrUnitInfo {
		if unitInfo.HasGPU {
			return true
		}
	}
	return false
}

func (uc *UseCase) waitGPUOperatorReady(ctx context.Context, scope string) error {
	const interval = 5 * time.Second

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			daemonSet, err := uc.daemonSet.Get(ctx, scope, "gpu-operator", "nvidia-operator-validator")
			if k8serrors.IsNotFound(err) {
				continue // waiting until it appears or context is cancelled.
			}
			if err != nil {
				return err
			}

			status := daemonSet.Status

			if status.DesiredNumberScheduled == status.CurrentNumberScheduled {
				if status.DesiredNumberScheduled == status.NumberReady {
					return nil
				}
			}
		}
	}
}

func (uc *UseCase) hackAndWaitContainerdConfigChanged(ctx context.Context, scope, unit string, registries []registryConfig) error {
	const interval = 5 * time.Second

	before, err := uc.getContainerdConfig(ctx, scope, unit)
	if err != nil {
		return err
	}

	// add hack registries
	hack := slices.Clone(registries)
	hack = append(hack, registryConfig{
		URL:                "http://hack.reconcile",
		InsecureSkipVerify: true,
	})

	if err := uc.setContainerdCustomRegistries(ctx, scope, hack); err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			after, err := uc.getContainerdConfig(ctx, scope, unit)
			if err != nil {
				return err
			}

			if before != after {
				return nil
			}
		}
	}
}
