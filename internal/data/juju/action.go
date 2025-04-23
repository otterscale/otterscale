package juju

import (
	"context"

	api "github.com/juju/juju/api/client/action"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type action struct {
	jujuMap JujuMap
}

func NewAction(jujuMap JujuMap) service.JujuAction {
	return &action{
		jujuMap: jujuMap,
	}
}

var _ service.JujuAction = (*action)(nil)

func (r *action) List(_ context.Context, uuid, appName string) (map[string]api.ActionSpec, error) {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).ApplicationCharmActions(appName)
}
