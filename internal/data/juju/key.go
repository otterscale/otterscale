package juju

import (
	"context"

	"github.com/juju/juju/api/client/keymanager"
	"github.com/juju/juju/rpc/params"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type key struct {
	jujuMap JujuMap
}

func NewKey(jujuMap JujuMap) service.JujuKey {
	return &key{
		jujuMap: jujuMap,
	}
}

var _ service.JujuKey = (*key)(nil)

func (r *key) Add(_ context.Context, uuid, key string) ([]params.ErrorResult, error) {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return nil, err
	}
	return keymanager.NewClient(conn).AddKeys(r.jujuMap.Username, key)
}
