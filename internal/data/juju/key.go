package juju

import (
	"context"
	"errors"

	api "github.com/juju/juju/api/client/keymanager"

	"github.com/otterscale/otterscale/internal/core"
)

type key struct {
	juju *Juju
}

func NewKey(juju *Juju) core.KeyRepo {
	return &key{
		juju: juju,
	}
}

var _ core.KeyRepo = (*key)(nil)

func (r *key) Add(_ context.Context, scope, key string) error {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return err
	}

	results, err := api.NewClient(conn).AddKeys(r.juju.username(), key)
	if err != nil {
		return err
	}
	errs := []error{}
	for _, result := range results {
		if result.Error != nil {
			errs = append(errs, result.Error)
		}
	}
	return errors.Join(errs...)
}
