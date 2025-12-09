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

	"github.com/otterscale/otterscale/internal/core/machine"
)

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

func (uc *UseCase) getContainerdTemplate(ctx context.Context, scope, unit string) (string, error) {
	path := uc.containerdTemplatePath(unit)
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

func (uc *UseCase) setContainerdTemplate(ctx context.Context, scope, unit, config string) error {
	path := uc.containerdTemplatePath(unit)

	data := base64.StdEncoding.EncodeToString([]byte(config))
	cmd := fmt.Sprintf(`echo "%s" | base64 -d > %s`, data, path)

	_, err := uc.action.Execute(ctx, scope, unit, cmd)
	return err
}

func (uc *UseCase) reconcileContainerd(ctx context.Context, scope, unit string) error {
	// get custom registries
	config, err := uc.getConfig(ctx, scope, "containerd")
	if err != nil {
		return err
	}

	customRegistries, err := getValue[string](config, "custom_registries")
	if err != nil {
		return err
	}

	var registries []registryConfig
	if err := json.Unmarshal([]byte(customRegistries), &registries); err != nil {
		return err
	}

	// add hack registries
	hack := slices.Clone(registries)
	hack = append(hack, registryConfig{
		URL:                "http://hack.reconcile",
		InsecureSkipVerify: true,
	})

	hackValue, err := json.Marshal(hack)
	if err != nil {
		return err
	}

	// set hacked config
	if err := uc.hackAndWaitContainerdConfigChanged(ctx, scope, unit, string(hackValue)); err != nil {
		return err
	}

	// set config back
	oriValue, err := json.Marshal(registries)
	if err != nil {
		return err
	}

	return uc.setConfig(ctx, scope, "containerd", "custom_registries", string(oriValue))
}

func (uc *UseCase) patchContainerdTemplates(ctx context.Context, scope string, wait bool) error {
	unitWithGPU, err := uc.getContainerdUnitContainsGPU(ctx, scope)
	if err != nil {
		return err
	}

	for unit, hasGPU := range unitWithGPU {
		// update by gpu operator
		if wait {
			if err := uc.waitContainerdTemplateChanged(ctx, scope, unit); err != nil {
				return err
			}
		}

		config := renderContainerdTemplate(hasGPU)

		if err := uc.setContainerdTemplate(ctx, scope, unit, config); err != nil {
			return err
		}

		if err := uc.reconcileContainerd(ctx, scope, unit); err != nil {
			return err
		}
	}

	return nil
}

func (uc *UseCase) getContainerdUnitContainsGPU(ctx context.Context, scope string) (map[string]bool, error) {
	machineWithGPU, err := uc.findMachineWithGPU(ctx)
	if err != nil {
		return nil, err
	}

	return uc.findContainerdUnitWithGPU(ctx, scope, machineWithGPU)
}

func (uc *UseCase) findContainerdUnitWithGPU(ctx context.Context, scope string, machineWithGPU map[string]bool) (map[string]bool, error) {
	facilities, err := uc.facility.List(ctx, scope, "")
	if err != nil {
		return nil, err
	}

	unitWithGPU := map[string]bool{}

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

				if hasGPU, ok := machineWithGPU[unit.Machine]; ok {
					unitWithGPU[subName] = hasGPU
				}
			}
		}
	}

	return unitWithGPU, nil
}

func (uc *UseCase) findMachineWithGPU(ctx context.Context) (map[string]bool, error) {
	machines, err := uc.machine.List(ctx)
	if err != nil {
		return nil, err
	}

	machineWithGPU := map[string]bool{}

	for i := range machines {
		jujuID, err := uc.machine.ExtractJujuID(&machines[i])
		if err != nil {
			continue
		}

		gpus, err := uc.nodeDevice.ListGPUs(ctx, machines[i].SystemID)
		if err != nil {
			return nil, err
		}

		machineWithGPU[jujuID] = uc.containsNvidiaGPU(gpus)
	}

	return machineWithGPU, nil
}

func (uc *UseCase) containsNvidiaGPU(gpus []machine.GPU) bool {
	for i := range gpus {
		if strings.EqualFold(gpus[i].VendorID, "10de") {
			return true
		}
	}
	return false
}

func (uc *UseCase) waitContainerdTemplateChanged(ctx context.Context, scope, unit string) error {
	const interval = 5 * time.Second

	before, err := uc.getContainerdTemplate(ctx, scope, unit)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			after, err := uc.getContainerdTemplate(ctx, scope, unit)
			if err != nil {
				return err
			}

			if before != after {
				return nil
			}
		}
	}
}

func (uc *UseCase) hackAndWaitContainerdConfigChanged(ctx context.Context, scope, unit, hackValue string) error {
	const interval = 5 * time.Second

	before, err := uc.getContainerdConfig(ctx, scope, unit)
	if err != nil {
		return err
	}

	if err := uc.setConfig(ctx, scope, "containerd", "custom_registries", string(hackValue)); err != nil {
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
