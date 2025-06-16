package juju

import (
	"context"
	"fmt"

	api "github.com/juju/juju/api/client/action"

	"github.com/openhdc/otterscale/internal/core"
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

func (r *action) List(_ context.Context, uuid, appName string) (map[string]core.ActionSpec, error) {
	conn, err := r.juju.connection(uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).ApplicationCharmActions(appName)
}

func (r *action) Run(_ context.Context, uuid, unitName, command string) (string, error) {
	conn, err := r.juju.connection(uuid)
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
		return "", fmt.Errorf("failed to run action %s on %s", command, unitName)
	}
	return enqueued.Actions[0].Action.ID, nil
}

func (r *action) GetResult(_ context.Context, uuid, id string) (*api.ActionResult, error) {
	conn, err := r.juju.connection(uuid)
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
