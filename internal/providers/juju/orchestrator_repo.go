package juju

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/juju/juju/api/client/client"

	"github.com/otterscale/otterscale/internal/core/machine"
)

// Note: Juju API do not support context.
type orchestratorRepo struct {
	juju *Juju
}

func NewOrchestratorRepo(juju *Juju) machine.OrchestratorRepo {
	return &orchestratorRepo{
		juju: juju,
	}
}

var _ machine.OrchestratorRepo = (*orchestratorRepo)(nil)

func (r *orchestratorRepo) AgentStatus(_ context.Context, scope, jujuID string) (*machine.AgentStatus, error) {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return nil, err
	}

	args := &client.StatusArgs{
		Patterns: []string{"machine", jujuID},
	}

	fullStatus, err := client.NewClient(conn, nil).Status(args)
	if err != nil {
		return nil, err
	}

	machineStatus, ok := fullStatus.Machines[jujuID]
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("agent with ID %q not found", jujuID))
	}

	return &machine.AgentStatus{
		Name:    machineStatus.AgentStatus.Status,
		Message: machineStatus.InstanceStatus.Info,
	}, nil
}
