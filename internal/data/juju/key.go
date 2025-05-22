package juju

import (
	"context"

	api "github.com/juju/juju/api/client/keymanager"
	"github.com/juju/juju/rpc/params"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type key struct {
	juju *Juju
}

func NewKey(juju *Juju) service.JujuKey {
	return &key{
		juju: juju,
	}
}

var _ service.JujuKey = (*key)(nil)

func (r *key) Add(_ context.Context, uuid, key string) ([]params.ErrorResult, error) {
	conn, err := r.juju.connection(uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).AddKeys(r.juju.username(), key)
}
