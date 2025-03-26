package juju

import (
	"context"

	api "github.com/juju/juju/api/client/client"
	"github.com/juju/juju/rpc/params"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type client struct{}

func NewClient() service.JujuClient {
	return &client{}
}

var _ service.JujuClient = (*client)(nil)

func (r *client) Status(ctx context.Context, uuid string) (*params.FullStatus, error) {
	conn, err := newConnection(uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn, nil).Status(ctx, nil)
}
