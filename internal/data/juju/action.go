package juju

import (
	"context"

	api "github.com/juju/juju/api/client/action"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type action struct{}

func NewAction() service.JujuAction {
	return &action{}
}

var _ service.JujuAction = (*action)(nil)

func (r *action) List(ctx context.Context, uuid, appName string) (map[string]api.ActionSpec, error) {
	conn, err := newConnection(uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).ApplicationCharmActions(ctx, appName)
}
