package juju

import (
	"context"

	"github.com/juju/juju/api/client/machinemanager"
	"github.com/juju/juju/rpc/params"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type machine struct {
	jujuMap JujuMap
}

func NewMachine(jujuMap JujuMap) service.JujuMachine {
	return &machine{
		jujuMap: jujuMap,
	}
}

var _ service.JujuMachine = (*machine)(nil)

func (r *machine) AddMachines(_ context.Context, uuid string, params []params.AddMachineParams) ([]params.AddMachinesResult, error) {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return nil, err
	}
	return machinemanager.NewClient(conn).AddMachines(params)
}
