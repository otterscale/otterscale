package juju

import (
	"context"

	api "github.com/juju/juju/api/client/action"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type action struct {
	juju *Juju
}

func NewAction(juju *Juju) service.JujuAction {
	return &action{
		juju: juju,
	}
}

var _ service.JujuAction = (*action)(nil)

func (r *action) List(_ context.Context, uuid, appName string) (map[string]api.ActionSpec, error) {
	conn, err := r.juju.connection(uuid)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return api.NewClient(conn).ApplicationCharmActions(appName)
}
