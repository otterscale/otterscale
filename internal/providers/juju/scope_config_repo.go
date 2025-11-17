package juju

import (
	"context"

	"github.com/juju/juju/api/client/modelconfig"
	"github.com/juju/juju/api/client/modelmanager"

	"github.com/otterscale/otterscale/internal/core/configuration"
)

// Note: Juju API do not support context.
type scopeConfigRepo struct {
	juju *Juju
}

func NewScopeConfigRepo(juju *Juju) configuration.ScopeConfigRepo {
	return &scopeConfigRepo{
		juju: juju,
	}
}

var _ configuration.ScopeConfigRepo = (*scopeConfigRepo)(nil)

func (r *scopeConfigRepo) Set(_ context.Context, config map[string]any) error {
	conn, err := r.juju.connection("controller")
	if err != nil {
		return err
	}

	summaries, err := modelmanager.NewClient(conn).ListModelSummaries(r.juju.conf.JujuUsername(), true)
	if err != nil {
		return err
	}

	for i := range summaries {
		conn, err := r.juju.connection(summaries[i].Name)
		if err != nil {
			return err
		}

		if err := modelconfig.NewClient(conn).ModelSet(config); err != nil {
			return err
		}
	}

	return nil
}
