package purge

import (
	"context"
	"fmt"
	"strings"

	"github.com/otterscale/otterscale/internal/core/facility"
	"github.com/otterscale/otterscale/internal/core/facility/action"
	"github.com/otterscale/otterscale/internal/core/machine"
)

type PurgeUseCase struct {
	action   action.ActionRepo
	facility facility.FacilityRepo
	machine  machine.MachineRepo
}

func NewPurgeUseCase(action action.ActionRepo, facility facility.FacilityRepo, machine machine.MachineRepo) *PurgeUseCase {
	return &PurgeUseCase{
		action:   action,
		facility: facility,
		machine:  machine,
	}
}

// TODO: improve performance by parallel execution
// TODO: osd devices in unit config or app config
func (uc *PurgeUseCase) PurgeDisk(ctx context.Context, machineID string) error {
	machine, err := uc.machine.Get(ctx, machineID)
	if err != nil {
		return err
	}

	scope, err := uc.machine.ExtractScope(machine)
	if err != nil {
		return err
	}

	jujuID, err := uc.machine.ExtractJujuID(machine)
	if err != nil {
		return err
	}

	apps, err := uc.facility.List(ctx, scope, jujuID)
	if err != nil {
		return err
	}

	for _, app := range apps {
		if app.Name != scope+"-ceph-osd" {
			continue
		}

		config, err := uc.facility.Config(ctx, scope, app.Name)
		if err != nil {
			continue
		}

		info, ok := config["osd-devices"].(map[string]interface{})
		if !ok {
			continue
		}

		val, ok := info["value"]
		if !ok || val == nil {
			continue
		}

		osdDevices := strings.SplitSeq(val.(string), " ")

		for osdDevice := range osdDevices {
			for unitName := range app.Status.Units {
				if _, err := uc.action.Execute(ctx, scope, unitName, fmt.Sprintf("sudo dd if=/dev/zero of=%s bs=1M count=200000", osdDevice)); err != nil {
					continue
				}
				break
			}
		}
	}

	return nil
}
