package juju

import (
	"context"

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

func (r *action) ListResults(_ context.Context, uuid, id string) ([]api.ActionResult, error) {
	conn, err := r.juju.connection(uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).Actions([]string{id})
}

func (r *action) ListOperations(_ context.Context, uuid, unitName, actionName string) ([]api.Operation, error) {
	conn, err := r.juju.connection(uuid)
	if err != nil {
		return nil, err
	}
	args := api.OperationQueryArgs{
		Units:       []string{unitName},
		ActionNames: []string{actionName},
	}
	ops, err := api.NewClient(conn).ListOperations(args)
	if err != nil {
		return nil, err
	}
	return ops.Operations, nil
}

func (r *action) Run(_ context.Context, uuid, unitName, command string) ([]api.ActionResult, error) {
	conn, err := r.juju.connection(uuid)
	if err != nil {
		return nil, err
	}
	parallel := true
	args := api.RunParams{
		Commands: command,
		Units:    []string{unitName},
		Parallel: &parallel,
	}
	enqueued, err := api.NewClient(conn).Run(args)
	if err != nil {
		return nil, err
	}
	return enqueued.Actions, nil
}
