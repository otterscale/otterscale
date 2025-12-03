package juju

import (
	"context"

	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/api/client/applicationoffers"
	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/names/v5"

	"github.com/otterscale/otterscale/internal/core/facility"
)

// Note: Juju API do not support context.
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
	offers, err := r.offers(url)
	if err != nil {
		return err
	}

	return r.consume(scope, offers)
}

func (r *relationRepo) offers(url string) (*params.ConsumeOfferDetails, error) {
	conn, err := r.juju.connection("controller")
	if err != nil {
		return nil, err
	}

	offers, err := applicationoffers.NewClient(conn).GetConsumeDetails(url)
	if err != nil {
		return nil, err
	}

	offerURL, err := crossmodel.ParseOfferURL(offers.Offer.OfferURL)
	if err != nil {
		return nil, err
	}

	offerURL.Source = r.juju.conf.JujuController()
	offers.Offer.OfferURL = offerURL.String()

	return &offers, nil
}

func (r *relationRepo) consume(scope string, offers *params.ConsumeOfferDetails) error {
	conn, err := r.juju.connection(scope)
	if err != nil {
		return err
	}

	args := crossmodel.ConsumeApplicationArgs{
		Offer:    *offers.Offer,
		Macaroon: offers.Macaroon,
	}

	if offers.ControllerInfo != nil {
		controllerTag, err := names.ParseControllerTag(offers.ControllerInfo.ControllerTag)
		if err != nil {
			return err
		}

		args.ControllerInfo = &crossmodel.ControllerInfo{
			ControllerTag: controllerTag,
			Alias:         offers.ControllerInfo.Alias,
			Addrs:         offers.ControllerInfo.Addrs,
			CACert:        offers.ControllerInfo.CACert,
		}
	}

	_, err = application.NewClient(conn).Consume(args)
	return err
}
