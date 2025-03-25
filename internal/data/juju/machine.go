package juju

import (
	"context"

	"github.com/juju/juju/api/client/machinemanager"
	"github.com/juju/juju/rpc/params"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type machine struct{}

func NewMachine() service.JujuMachine {
	return &machine{}
}

var _ service.JujuMachine = (*machine)(nil)

func (r *machine) AddMachines(ctx context.Context, uuid string, params []params.AddMachineParams) ([]params.AddMachinesResult, error) {
	conn, err := newConnection(uuid)
	if err != nil {
		return nil, err
	}
	return machinemanager.NewClient(conn).AddMachines(ctx, params)
}
