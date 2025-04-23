package juju

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/juju/juju/api/client/application"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/names/v5"

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

func (r *application) Create(_ context.Context, uuid, name, configYAML, charmName, channel string, revision, number int, placements []instance.Placement, constraint *constraints.Value, trust bool) (*api.DeployInfo, error) {
	conn, err := r.jujuMap.Get(uuid)
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
		Cons:            *constraint,
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

	di, _, errs := api.NewClient(conn).DeployFromRepository(args)
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return &di, nil
}

func (r *application) Update(_ context.Context, uuid, name, configYAML string) error {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return err
	}
	return api.NewClient(conn).SetConfig("", name, configYAML, nil)
}

func (r *application) Delete(_ context.Context, uuid, name string, destroyStorage, force bool) error {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return err
	}
	dars, err := api.NewClient(conn).DestroyApplications(api.DestroyApplicationsParams{
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

func (r *application) Expose(_ context.Context, uuid, name string, endpoints map[string]params.ExposedEndpoint) error {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return err
	}
	return api.NewClient(conn).Expose(name, endpoints)
}

func (r *application) AddUnits(_ context.Context, uuid, name string, number int, placements []instance.Placement) ([]string, error) {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return nil, err
	}
	ps := []*instance.Placement{}
	for _, p := range placements {
		ps = append(ps, &p)
	}
	return api.NewClient(conn).AddUnits(api.AddUnitsParams{
		ApplicationName: name,
		NumUnits:        number,
		Placement:       ps,
	})
}

func (r *application) ResolveUnitErrors(_ context.Context, uuid string, units []string) error {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return err
	}
	return api.NewClient(conn).ResolveUnitErrors(units, true, true)
}

func (r *application) CreateRelation(_ context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error) {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).AddRelation(endpoints, nil)
}

func (r *application) DeleteRelation(_ context.Context, uuid string, id int) error {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return err
	}
	return api.NewClient(conn).DestroyRelationId(id, nil, nil)
}

func (r *application) GetConfig(_ context.Context, uuid, name string) (map[string]any, error) {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return nil, err
	}
	app, err := api.NewClient(conn).Get("", name)
	if err != nil {
		return nil, err
	}
	return app.CharmConfig, nil
}

func (r *application) GetLeader(_ context.Context, uuid, name string) (string, error) {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return "", err
	}
	return api.NewClient(conn).Leader(name)
}

func (r *application) GetUnitInfo(_ context.Context, uuid, name string) (*api.UnitInfo, error) {
	conn, err := r.jujuMap.Get(uuid)
	if err != nil {
		return nil, err
	}
	uis, err := api.NewClient(conn).UnitsInfo([]names.UnitTag{names.NewUnitTag(name)})
	if err != nil {
		return nil, err
	}
	if len(uis) == 1 {
		return &uis[0], nil
	}
	return nil, status.Errorf(codes.NotFound, "unit info %q not found", name)
}
