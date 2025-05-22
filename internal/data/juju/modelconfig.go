package juju

import (
	"context"

	api "github.com/juju/juju/api/client/modelconfig"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type modelConfig struct {
	juju *Juju
}

func NewModelConfig(juju *Juju) service.JujuModelConfig {
	return &modelConfig{
		juju: juju,
	}
}

var _ service.JujuModelConfig = (*modelConfig)(nil)

func (r *modelConfig) List(_ context.Context, uuid string) (map[string]any, error) {
	conn, err := r.juju.connection(uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).ModelGet()
}

func (r *modelConfig) Set(_ context.Context, uuid string, config map[string]any) error {
	conn, err := r.juju.connection(uuid)
	if err != nil {
		return err
	}
	return api.NewClient(conn).ModelSet(config)
}

func (r *modelConfig) Unset(_ context.Context, uuid string, keys ...string) error {
	conn, err := r.juju.connection(uuid)
	if err != nil {
		return err
	}
	return api.NewClient(conn).ModelUnset(keys...)
}
