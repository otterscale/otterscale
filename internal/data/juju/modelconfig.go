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

func (r *modelConfig) List(ctx context.Context, uuid string) (map[string]any, error) {
	conn, err := r.jujuMap.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return modelconfig.NewClient(conn).ModelGet(ctx)
}

func (r *modelConfig) Set(ctx context.Context, uuid string, config map[string]any) error {
	conn, err := r.jujuMap.Get(ctx, uuid)
	if err != nil {
		return err
	}
	return modelconfig.NewClient(conn).ModelSet(ctx, config)
}

func (r *modelConfig) Unset(ctx context.Context, uuid string, keys ...string) error {
	conn, err := r.jujuMap.Get(ctx, uuid)
	if err != nil {
		return err
	}
	return modelconfig.NewClient(conn).ModelUnset(ctx, keys...)
}
