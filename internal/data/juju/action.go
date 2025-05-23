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
