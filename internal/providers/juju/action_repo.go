package juju

import (
	"context"

	api "github.com/juju/juju/api/client/action"

	"github.com/otterscale/otterscale/internal/core/facility/action"
)

type actionRepo struct {
	juju *Juju
}

func NewActionRepo(juju *Juju) action.ActionRepo {
	return &actionRepo{
		juju: juju,
	}
}

var _ action.ActionRepo = (*actionRepo)(nil)

func (r *actionRepo) List(_ context.Context, scope, appName string) ([]action.Action, error) {
	conn, err := r.juju.Connection(scope)
	if err != nil {
		return nil, err
	}

	m, err := api.NewClient(conn).ApplicationCharmActions(appName)
	if err != nil {
		return nil, err
	}

	return r.toActions(m), nil
}

func (r *actionRepo) Run(ctx context.Context, scope, appName, actionName string, params map[string]any) (map[string]any, error) {
	return r.juju.Run(ctx, scope, appName, actionName, params)
}

func (r *actionRepo) Execute(ctx context.Context, scope, appName, command string) (map[string]any, error) {
	return r.juju.Execute(ctx, scope, appName, command)
}

func (r *actionRepo) toActions(m map[string]api.ActionSpec) []action.Action {
	ret := []action.Action{}
	return ret
}
