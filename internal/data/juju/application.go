package juju

import (
	"context"

	api "github.com/juju/juju/api/client/application"
	"github.com/juju/juju/rpc/params"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type application struct{}

func NewApplication() service.JujuApplication {
	return &application{}
}

var _ service.JujuApplication = (*application)(nil)

func (r *application) ResolveUnitErrors(ctx context.Context, uuid string, units []string) error {
	conn, err := newConnection(uuid)
	if err != nil {
		return err
	}
	return api.NewClient(conn).ResolveUnitErrors(ctx, units, true, true)
}

func (r *application) CreateRelation(ctx context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error) {
	conn, err := newConnection(uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).AddRelation(ctx, endpoints, nil)
}

func (r *application) DeleteRelation(ctx context.Context, uuid string, id int) error {
	conn, err := newConnection(uuid)
	if err != nil {
		return err
	}
	return api.NewClient(conn).DestroyRelationId(ctx, id, nil, nil)
}
