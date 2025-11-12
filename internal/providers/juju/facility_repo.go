package juju

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/api/client/client"
	"github.com/juju/juju/core/base"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/names/v5"

	"github.com/otterscale/otterscale/internal/core/facility"
)

type facilityRepo struct {
	juju *Juju
}

func NewFacilityRepo(juju *Juju) facility.FacilityRepo {
	return &facilityRepo{
		juju: juju,
	}
}

var _ facility.FacilityRepo = (*facilityRepo)(nil)

func (r *facilityRepo) List(ctx context.Context, scope, jujuID string) ([]facility.Facility, error) {
	conn, err := r.juju.Connection(scope)
	if err != nil {
		return nil, err
	}

	args := &client.StatusArgs{}

	if jujuID != "" {
		args.Patterns = []string{"machine", jujuID}
	}

	fullStatus, err := client.NewClient(conn, nil).Status(args)
	if err != nil {
		return nil, err
	}

	return r.toFacility(fullStatus.Applications), nil
}

func (r *facilityRepo) Create(_ context.Context, scope, name, configYAML, charmName, channel string, revision int, series, directive, placementScope string) error {
	conn, err := r.juju.Connection(scope)
	if err != nil {
		return err
	}

	base, err := base.GetBaseFromSeries(series)
	if err != nil {
		return err
	}

	args := application.DeployFromRepositoryArg{
		ApplicationName: name,
		ConfigYAML:      configYAML,
		CharmName:       charmName,
		Base:            &base,
		Trust:           true,
	}

	if channel != "" {
		args.Channel = &channel
	}

	if revision != 0 {
		args.Revision = &revision
	}

	if directive != "" {
		args.Placement = append(args.Placement, &instance.Placement{
			Scope:     placementScope,
			Directive: directive,
		})
	}

	_, _, errs := application.NewClient(conn).DeployFromRepository(args)
	return errors.Join(errs...)
}

// Note: This function has not been tested.
func (r *facilityRepo) Update(_ context.Context, scope, name, configYAML string) error {
	conn, err := r.juju.Connection(scope)
	if err != nil {
		return err
	}

	return application.NewClient(conn).SetConfig("", name, configYAML, nil)
}

// Note: This function has not been tested.
func (r *facilityRepo) Delete(_ context.Context, scope, name string, destroyStorage, force bool) error {
	conn, err := r.juju.Connection(scope)
	if err != nil {
		return err
	}

	results, err := application.NewClient(conn).DestroyApplications(application.DestroyApplicationsParams{
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

func (r *facilityRepo) Resolve(_ context.Context, scope, unitName string) error {
	conn, err := r.juju.Connection(scope)
	if err != nil {
		return err
	}

	units := []string{unitName}
	return application.NewClient(conn).ResolveUnitErrors(units, false, true)
}

func (r *facilityRepo) Config(_ context.Context, scope, name string) (map[string]any, error) {
	conn, err := r.juju.Connection(scope)
	if err != nil {
		return nil, err
	}

	app, err := application.NewClient(conn).Get("", name)
	if err != nil {
		return nil, err
	}
	return app.CharmConfig, nil
}

func (r *facilityRepo) PublicAddress(_ context.Context, scope, name string) (string, error) {
	conn, err := r.juju.Connection(scope)
	if err != nil {
		return "", err
	}

	leader, err := application.NewClient(conn).Leader(name)
	if err != nil {
		return "", err
	}

	tag := names.NewUnitTag(leader)
	units := []names.UnitTag{tag}

	unitInfos, err := application.NewClient(conn).UnitsInfo(units)
	if err != nil {
		return "", err
	}

	for i := range unitInfos {
		if unitInfos[i].Tag == tag.String() {
			return unitInfos[i].PublicAddress, nil
		}
	}

	return "", connect.NewError(connect.CodeNotFound, fmt.Errorf("unit info %q not found", name))
}

func (r *facilityRepo) toFacility(apps map[string]params.ApplicationStatus) []facility.Facility {
	ret := []facility.Facility{}

	return ret
}
