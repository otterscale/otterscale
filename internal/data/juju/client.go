package juju

import (
	"context"

	api "github.com/juju/juju/api/client/client"
	"github.com/juju/juju/rpc/params"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type client struct {
	juju *Juju
}

func NewClient(juju *Juju) service.JujuClient {
	return &client{
		juju: juju,
	}
}

var _ service.JujuClient = (*client)(nil)

func (r *client) Status(_ context.Context, uuid string, patterns []string) (*params.FullStatus, error) {
	conn, err := r.juju.connection(uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn, nil).Status(&api.StatusArgs{Patterns: patterns})
}
