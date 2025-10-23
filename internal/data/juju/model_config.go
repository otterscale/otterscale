package juju

import (
	"context"

	api "github.com/juju/juju/api/client/modelconfig"

	"github.com/otterscale/otterscale/internal/core"
)

type modelConfig struct {
	juju *Juju
}

func NewModelConfig(juju *Juju) core.ScopeConfigRepo {
	return &modelConfig{
		juju: juju,
	}
}

var _ core.ScopeConfigRepo = (*modelConfig)(nil)

func (r *modelConfig) List(_ context.Context, scope string) (map[string]any, error) {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).ModelGet()
}

func (r *modelConfig) Set(_ context.Context, scope string, config map[string]any) error {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return err
	}
	return api.NewClient(conn).ModelSet(config)
}

func (r *modelConfig) Unset(_ context.Context, scope string, keys ...string) error {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return err
	}
	return api.NewClient(conn).ModelUnset(keys...)
}
