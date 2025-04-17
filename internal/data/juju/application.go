package juju

import (
	"context"
	"fmt"

	api "github.com/juju/juju/api/client/application"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/names/v6"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type application struct {
	jujuMap JujuMap
}

func NewApplication(jujuMap JujuMap) service.JujuApplication {
	return &application{
		jujuMap: jujuMap,
	}
}

var _ service.JujuApplication = (*application)(nil)

func (r *application) Create(ctx context.Context, uuid, name string, configYAML string, charmName, channel string, revision, number int, placements []instance.Placement, constraint constraints.Value, trust bool) (*api.DeployInfo, error) {
	conn, err := r.jujuMap.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}

	ps := []*instance.Placement{}
	for _, p := range placements {
		ps = append(ps, &p)
	}

	args := api.DeployFromRepositoryArg{
		CharmName:       charmName,
		ApplicationName: name,
		ConfigYAML:      configYAML,
		Cons:            constraint,
		Placement:       ps,
		Trust:           trust,
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

	di, _, errs := api.NewClient(conn).DeployFromRepository(ctx, args)
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return &di, nil
}

func (r *application) Update(ctx context.Context, uuid, name string, configYAML string) error {
	conn, err := r.jujuMap.Get(ctx, uuid)
	if err != nil {
		return err
	}
	return api.NewClient(conn).SetConfig(ctx, name, configYAML, nil)
}

func (r *application) Delete(ctx context.Context, uuid, name string, destroyStorage, force bool) error {
	conn, err := r.jujuMap.Get(ctx, uuid)
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
	conn, err := r.jujuMap.Get(ctx, uuid)
	if err != nil {
		return err
	}
	return api.NewClient(conn).Expose(ctx, name, endpoints)
}

func (r *application) AddUnits(ctx context.Context, uuid, name string, number int, placements []instance.Placement) ([]string, error) {
	conn, err := r.jujuMap.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}
	ps := []*instance.Placement{}
	for _, p := range placements {
		ps = append(ps, &p)
	}
	return api.NewClient(conn).AddUnits(ctx, api.AddUnitsParams{
		ApplicationName: name,
		NumUnits:        number,
		Placement:       ps,
	})
}

func (r *application) ResolveUnitErrors(ctx context.Context, uuid string, units []string) error {
	conn, err := r.jujuMap.Get(ctx, uuid)
	if err != nil {
		return err
	}
	return api.NewClient(conn).ResolveUnitErrors(ctx, units, true, true)
}

func (r *application) CreateRelation(ctx context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error) {
	conn, err := r.jujuMap.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).AddRelation(ctx, endpoints, nil)
}

func (r *application) DeleteRelation(ctx context.Context, uuid string, id int) error {
	conn, err := r.jujuMap.Get(ctx, uuid)
	if err != nil {
		return err
	}
	return api.NewClient(conn).DestroyRelationId(ctx, id, nil, nil)
}

func (r *application) GetConfig(ctx context.Context, uuid string, name string) (map[string]any, error) {
	conn, err := r.jujuMap.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}
	app, err := api.NewClient(conn).Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return app.CharmConfig, nil
}

func (r *application) GetLeader(ctx context.Context, uuid, name string) (string, error) {
	conn, err := r.jujuMap.Get(ctx, uuid)
	if err != nil {
		return "", err
	}
	return api.NewClient(conn).Leader(ctx, name)
}

func (r *application) GetUnitInfo(ctx context.Context, uuid, name string) (*api.UnitInfo, error) {
	conn, err := r.jujuMap.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}
	uis, err := api.NewClient(conn).UnitsInfo(ctx, []names.UnitTag{names.NewUnitTag(name)})
	if err != nil {
		return nil, err
	}
	if len(uis) == 1 {
		return &uis[0], nil
	}
	return nil, fmt.Errorf("unit info %q not found", name)
}
