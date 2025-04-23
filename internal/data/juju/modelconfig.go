package juju

import (
	"context"

	"github.com/juju/juju/api/client/modelconfig"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type modelConfig struct {
	jujuMap JujuMap
}

func NewModelConfig(jujuMap JujuMap) service.JujuModelConfig {
	return &modelConfig{
		jujuMap: jujuMap,
	}
}

var _ service.JujuModelConfig = (*modelConfig)(nil)

func (r *modelConfig) List(_ context.Context, uuid string) (map[string]any, error) {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return nil, err
	}
	return modelconfig.NewClient(conn).ModelGet()
}

func (r *modelConfig) Set(_ context.Context, uuid string, config map[string]any) error {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return err
	}
	return modelconfig.NewClient(conn).ModelSet(config)
}

func (r *modelConfig) Unset(_ context.Context, uuid string, keys ...string) error {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return err
	}
	return modelconfig.NewClient(conn).ModelUnset(keys...)
}
