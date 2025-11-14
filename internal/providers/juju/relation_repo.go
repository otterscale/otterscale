package juju

import (
	"context"

	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/api/client/applicationoffers"
	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/names/v5"

	"github.com/otterscale/otterscale/internal/core/facility"
)

type relationRepo struct {
	juju *Juju
}

func NewRelationRepo(juju *Juju) facility.RelationRepo {
	return &relationRepo{
		juju: juju,
	}
}

var _ facility.RelationRepo = (*relationRepo)(nil)

func (r *relationRepo) Create(_ context.Context, scope string, endpoints []string) error {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return err
	}

	_, err = application.NewClient(conn).AddRelation(endpoints, nil)
	return err
}

func (r *relationRepo) Delete(_ context.Context, scope string, id int) error {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return err
	}

	return application.NewClient(conn).DestroyRelationId(id, nil, nil)
}

func (r *relationRepo) Consume(_ context.Context, scope, url string) error {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return err
	}

	details, err := applicationoffers.NewClient(conn).GetConsumeDetails(url)
	if err != nil {
		return err
	}

	offerURL, err := crossmodel.ParseOfferURL(details.Offer.OfferURL)
	if err != nil {
		return err
	}

	offerURL.Source = r.juju.controller()
	details.Offer.OfferURL = offerURL.String()

	args := crossmodel.ConsumeApplicationArgs{
		Offer:    *details.Offer,
		Macaroon: details.Macaroon,
	}

	if details.ControllerInfo != nil {
		controllerTag, err := names.ParseControllerTag(details.ControllerInfo.ControllerTag)
		if err != nil {
			return err
		}

		args.ControllerInfo = &crossmodel.ControllerInfo{
			ControllerTag: controllerTag,
			Alias:         details.ControllerInfo.Alias,
			Addrs:         details.ControllerInfo.Addrs,
			CACert:        details.ControllerInfo.CACert,
		}
	}

	_, err = application.NewClient(conn).Consume(args)
	return err
}
