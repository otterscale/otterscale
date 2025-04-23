package juju

import (
	"context"

	api "github.com/juju/juju/api/client/client"
	"github.com/juju/juju/rpc/params"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type client struct {
	jujuMap JujuMap
}

func NewClient(jujuMap JujuMap) service.JujuClient {
	return &client{
		jujuMap: jujuMap,
	}
}

var _ service.JujuClient = (*client)(nil)

func (r *client) Status(_ context.Context, uuid string, patterns []string) (*params.FullStatus, error) {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn, nil).Status(&api.StatusArgs{Patterns: patterns})
}
