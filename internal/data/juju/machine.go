package juju

import (
	"context"
	"time"

	api "github.com/juju/juju/api/client/machinemanager"
	"github.com/juju/juju/rpc/params"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type machine struct {
	juju *Juju
}

func NewMachine(juju *Juju) service.JujuMachine {
	return &machine{
		juju: juju,
	}
}

var _ service.JujuMachine = (*machine)(nil)

func (r *machine) AddMachines(_ context.Context, uuid string, params []params.AddMachineParams) ([]params.AddMachinesResult, error) {
	conn, err := r.juju.connection(uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).AddMachines(params)
}

func (r *machine) DestroyMachines(_ context.Context, uuid string, force, keep, dryRun bool, maxWait *time.Duration, machines ...string) ([]params.DestroyMachineResult, error) {
	conn, err := r.juju.connection(uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).DestroyMachinesWithParams(force, keep, dryRun, maxWait, machines...)
}
