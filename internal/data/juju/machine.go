package juju

import (
	"context"
	"errors"
	"time"

	api "github.com/juju/juju/api/client/machinemanager"
	"github.com/juju/juju/rpc/params"

	"github.com/openhdc/otterscale/internal/core"
)

type machine struct {
	juju *Juju
}

func NewMachine(juju *Juju) core.MachineManagerRepo {
	return &machine{
		juju: juju,
	}
}

var _ core.MachineManagerRepo = (*machine)(nil)

func (r *machine) AddMachines(_ context.Context, uuid string, params []params.AddMachineParams) error {
	conn, err := r.juju.connection(uuid)
	if err != nil {
		return err
	}

	results, err := api.NewClient(conn).AddMachines(params)
	if err != nil {
		return err
	}
	errs := []error{}
	for _, result := range results {
		errs = append(errs, result.Error)
	}
	return errors.Join(errs...)
}

func (r *machine) DestroyMachines(_ context.Context, uuid string, force, keep, dryRun bool, maxWait *time.Duration, machines ...string) error {
	conn, err := r.juju.connection(uuid)
	if err != nil {
		return err
	}

	results, err := api.NewClient(conn).DestroyMachinesWithParams(force, keep, dryRun, maxWait, machines...)
	if err != nil {
		return err
	}
	errs := []error{}
	for _, result := range results {
		errs = append(errs, result.Error)
	}
	return errors.Join(errs...)
}
