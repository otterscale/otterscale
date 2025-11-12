package juju

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/juju/juju/api/client/client"
	"github.com/otterscale/otterscale/internal/core/machine"
)

type orchestratorRepo struct {
	juju *Juju
}

func NewOrchestratorRepo(juju *Juju) machine.OrchestratorRepo {
	return &orchestratorRepo{
		juju: juju,
	}
}

var _ machine.OrchestratorRepo = (*orchestratorRepo)(nil)

func (r *orchestratorRepo) AgentStatus(ctx context.Context, scope, jujuID string) (string, error) {
	conn, err := r.juju.Connection(scope)
	if err != nil {
		return "", err
	}

	args := &client.StatusArgs{
		Patterns: []string{"machine", jujuID},
	}

	fullStatus, err := client.NewClient(conn, nil).Status(args)
	if err != nil {
		return "", err
	}

	machineStatus, ok := fullStatus.Machines[jujuID]
	if !ok {
		return "", connect.NewError(connect.CodeNotFound, fmt.Errorf("agent with ID %q not found", jujuID))
	}

	return machineStatus.AgentStatus.Status, nil
}
