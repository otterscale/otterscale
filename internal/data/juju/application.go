package juju

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	api "github.com/juju/juju/api/client/application"
	"github.com/juju/juju/core/base"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/names/v5"

	"github.com/otterscale/otterscale/internal/core"
)

type application struct {
	juju *Juju
}

func NewApplication(juju *Juju) core.FacilityRepo {
	return &application{
		juju: juju,
	}
}

var _ core.FacilityRepo = (*application)(nil)

func (r *application) Create(_ context.Context, scope, name, configYAML, charmName, channel string, revision, number int, base *base.Base, placements []instance.Placement, constraint *constraints.Value, trust bool) (*api.DeployInfo, error) {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return nil, err
	}

	args := api.DeployFromRepositoryArg{
		ApplicationName: name,
		ConfigYAML:      configYAML,
		CharmName:       charmName,
		Base:            base,
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

	for _, p := range placements {
		args.Placement = append(args.Placement, &p)
	}

	if constraint != nil {
		args.Cons = *constraint
	}

	result, _, errs := api.NewClient(conn).DeployFromRepository(args)
	err = errors.Join(errs...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Note: This function has not been tested.
func (r *application) Update(_ context.Context, scope, name, configYAML string) error {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return err
	}
	return api.NewClient(conn).SetConfig("", name, configYAML, nil)
}

// Note: This function has not been tested.
func (r *application) Delete(_ context.Context, scope, name string, destroyStorage, force bool) error {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return err
	}

	results, err := api.NewClient(conn).DestroyApplications(api.DestroyApplicationsParams{
		Applications:   []string{name},
		DestroyStorage: destroyStorage,
		Force:          force,
	})
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

func (r *application) Expose(_ context.Context, scope, name string, endpoints map[string]params.ExposedEndpoint) error {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return err
	}
	return api.NewClient(conn).Expose(name, endpoints)
}

func (r *application) AddUnits(_ context.Context, scope, name string, number int, placements []instance.Placement) ([]string, error) {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return nil, err
	}

	params := api.AddUnitsParams{
		ApplicationName: name,
		NumUnits:        number,
	}
	for _, placement := range placements {
		params.Placement = append(params.Placement, &placement)
	}
	return api.NewClient(conn).AddUnits(params)
}

func (r *application) ResolveUnitErrors(_ context.Context, scope string, units []string) error {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return err
	}
	return api.NewClient(conn).ResolveUnitErrors(units, false, true)
}

func (r *application) CreateRelation(_ context.Context, scope string, endpoints []string) (*params.AddRelationResults, error) {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return nil, err
	}
	return api.NewClient(conn).AddRelation(endpoints, nil)
}

func (r *application) DeleteRelation(_ context.Context, scope string, id int) error {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return err
	}
	return api.NewClient(conn).DestroyRelationId(id, nil, nil)
}

func (r *application) GetConfig(_ context.Context, scope, name string) (map[string]any, error) {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return nil, err
	}

	app, err := api.NewClient(conn).Get("", name)
	if err != nil {
		return nil, err
	}
	return app.CharmConfig, nil
}

func (r *application) GetLeader(_ context.Context, scope, name string) (string, error) {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return "", err
	}
	return api.NewClient(conn).Leader(name)
}

func (r *application) GetUnitInfo(_ context.Context, scope, name string) (*api.UnitInfo, error) {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return nil, err
	}

	tag := names.NewUnitTag(name)
	units, err := api.NewClient(conn).UnitsInfo([]names.UnitTag{tag})
	if err != nil {
		return nil, err
	}

	for i := range units {
		if units[i].Tag == tag.String() {
			return &units[i], nil
		}
	}
	return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("unit info %q not found", name))
}

func (r *application) Consume(_ context.Context, scope string, args *crossmodel.ConsumeApplicationArgs) error {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return err
	}
	_, err = api.NewClient(conn).Consume(*args)
	return err
}
