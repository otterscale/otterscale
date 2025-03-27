package juju

import (
	"context"

	api "github.com/juju/juju/api/client/application"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type application struct{}

func NewApplication() service.JujuApplication {
	return &application{}
}

var _ service.JujuApplication = (*application)(nil)

func (r *application) Create(ctx context.Context, uuid, charmName, appName, channel string, revision, number int, config map[string]string, constraint constraints.Value, placements []*instance.Placement, trust bool) error {
	conn, err := newConnection(uuid)
	if err != nil {
		return err
	}

	args := api.DeployFromRepositoryArg{
		CharmName:       charmName,
		ApplicationName: appName,
		// ConfigYAML:      configYAML, // TODO: YAML
		Cons:      constraint,
		Placement: placements,
		Trust:     trust,
	}

	if channel != "" {
		args.Channel = &channel
	}

	if revision != 0 {
		args.Revision = &revision
	}

	if number != 0 {
		args.NumUnits = &number
	}

	_, _, errs := api.NewClient(conn).DeployFromRepository(ctx, args)
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *application) Update(ctx context.Context, uuid, name string, config map[string]string) error {
	conn, err := newConnection(uuid)
	if err != nil {
		return err
	}
	return api.NewClient(conn).SetConfig(ctx, name, "", config)
}

func (r *application) Delete(ctx context.Context, uuid, name string, destroyStorage, force bool) error {
	conn, err := newConnection(uuid)
	if err != nil {
		return err
	}
	dars, err := api.NewClient(conn).DestroyApplications(ctx, api.DestroyApplicationsParams{
		Applications:   []string{name},
		DestroyStorage: destroyStorage,
		Force:          force,
	})
	if err != nil {
		return err
	}
	for _, dar := range dars {
		if dar.Error != nil {
			return dar.Error
		}
	}
	return nil
}

func (r *application) Expose(ctx context.Context, uuid, name string, endpoints map[string]params.ExposedEndpoint) error {
	conn, err := newConnection(uuid)
	if err != nil {
		return err
	}
	return api.NewClient(conn).Expose(ctx, name, endpoints)
}

func (r *application) AddUnits(ctx context.Context, uuid, name string, number int, placements []*instance.Placement) ([]string, error) {
	conn, err := newConnection(uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).AddUnits(ctx, api.AddUnitsParams{
		ApplicationName: name,
		NumUnits:        number,
		Placement:       placements,
	})
}

func (r *application) ResolveUnitErrors(ctx context.Context, uuid string, units []string) error {
	conn, err := newConnection(uuid)
	if err != nil {
		return err
	}
	return api.NewClient(conn).ResolveUnitErrors(ctx, units, true, true)
}

func (r *application) CreateRelation(ctx context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error) {
	conn, err := newConnection(uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).AddRelation(ctx, endpoints, nil)
}

func (r *application) DeleteRelation(ctx context.Context, uuid string, id int) error {
	conn, err := newConnection(uuid)
	if err != nil {
		return err
	}
	return api.NewClient(conn).DestroyRelationId(ctx, id, nil, nil)
}

func (r *application) GetConfigs(ctx context.Context, uuid string, names ...string) ([]map[string]interface{}, error) {
	conn, err := newConnection(uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).GetConfig(ctx, names...)
}
