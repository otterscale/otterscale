package juju

import (
	"context"
	"errors"
	"time"

	"github.com/juju/juju/api/client/machinemanager"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/model"
	"github.com/juju/juju/rpc/params"

	"github.com/otterscale/otterscale/internal/core/machine"
)

// Note: Juju API do not support context.
type machineManagerRepo struct {
	juju *Juju
}

func NewMachineManagerRepo(juju *Juju) machine.MachineManagerRepo {
	return &machineManagerRepo{
		juju: juju,
	}
}

var _ machine.MachineManagerRepo = (*machineManagerRepo)(nil)

func (r *machineManagerRepo) AddMachines(_ context.Context, scope, uuid, fqdn, baseOS, baseChannel string) error {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return err
	}

	params := []params.AddMachineParams{
		{
			Placement: &instance.Placement{Scope: uuid, Directive: fqdn},
			Jobs:      []model.MachineJob{model.JobHostUnits},
			Base:      &params.Base{Name: baseOS, Channel: baseChannel},
		},
	}

	results, err := machinemanager.NewClient(conn).AddMachines(params)
	if err != nil {
		return err
	}

	errs := []error{}

	for _, result := range results {
		if result.Error != nil {
			errs = append(errs, result.Error)
		}
	}

	return errors.Join(errs...)
}

func (r *machineManagerRepo) DestroyMachines(_ context.Context, scope string, force, keep, dryRun bool, maxWait *time.Duration, machines ...string) error {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return err
	}

	results, err := machinemanager.NewClient(conn).DestroyMachinesWithParams(force, keep, dryRun, maxWait, machines...)
	if err != nil {
		return err
	}

	errs := []error{}

	for _, result := range results {
		if result.Error != nil {
			errs = append(errs, result.Error)
		}
	}

	return errors.Join(errs...)
}
