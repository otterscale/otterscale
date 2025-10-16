package juju

import (
	"context"

	api "github.com/juju/juju/api/client/client"
	"github.com/juju/juju/rpc/params"

	"github.com/otterscale/otterscale/internal/core"
)

type client struct {
	juju *Juju
}

func NewClient(juju *Juju) core.ClientRepo {
	return &client{
		juju: juju,
	}
}

var _ core.ClientRepo = (*client)(nil)

func (r *client) Status(_ context.Context, scope string, patterns []string) (*params.FullStatus, error) {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn, nil).Status(&api.StatusArgs{Patterns: patterns})
}
