package juju

import (
	"context"
	"fmt"

	api "github.com/juju/juju/api/client/action"
	"github.com/juju/names/v5"

	"github.com/otterscale/otterscale/internal/core"
)

type action struct {
	juju *Juju
}

func NewAction(juju *Juju) core.ActionRepo {
	return &action{
		juju: juju,
	}
}

var _ core.ActionRepo = (*action)(nil)

func (r *action) List(_ context.Context, scope, appName string) (map[string]core.FacilityActionSpec, error) {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).ApplicationCharmActions(appName)
}

func (r *action) RunCommand(_ context.Context, scope, unitName, command string) (string, error) {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return "", err
	}
	parallel := true
	args := api.RunParams{
		Commands: command,
		Units:    []string{unitName},
		Parallel: &parallel,
	}
	enqueued, err := api.NewClient(conn).Run(args)
	if err != nil {
		return "", err
	}
	if len(enqueued.Actions) == 0 || enqueued.Actions[0].Action == nil {
		return "", fmt.Errorf("failed to run command %q on %s", command, unitName)
	}
	return enqueued.Actions[0].Action.ID, nil
}

func (r *action) RunAction(_ context.Context, scope, unitName, actionName string, parameters map[string]any) (string, error) {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return "", err
	}
	actions := []api.Action{
		{
			Receiver:   names.NewUnitTag(unitName).String(),
			Name:       actionName,
			Parameters: parameters,
		},
	}
	enqueued, err := api.NewClient(conn).EnqueueOperation(actions)
	if err != nil {
		return "", err
	}
	if len(enqueued.Actions) == 0 || enqueued.Actions[0].Action == nil {
		return "", fmt.Errorf("failed to run action %q on %s", actionName, unitName)
	}
	return enqueued.Actions[0].Action.ID, nil
}

func (r *action) GetResult(_ context.Context, scope, id string) (*api.ActionResult, error) {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return nil, err
	}
	results, err := api.NewClient(conn).Actions([]string{id})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 || results[0].Action == nil {
		return nil, fmt.Errorf("failed to get action result %q", id)
	}
	return &results[0], nil
}
