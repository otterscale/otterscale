package core

import (
	"context"

	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/names/v5"
)

func (uc *OrchestratorUseCase) createCOS(ctx context.Context, scope, prefix string) error {
	// consume
	offerURLs := []string{
		uc.conf.Juju.Username + "/cos.global-prometheus",
		uc.conf.Juju.Username + "/cos.global-grafana",
	}
	for _, url := range offerURLs {
		if err := uc.consumeRemoteOffer(ctx, scope, url); err != nil {
			return err
		}
	}

	// integrate
	name := prefix + "-" + "grafana-agent"
	endpointList := [][]string{
		{name + ":send-remote-write", "global-prometheus:receive-remote-write"},
		{name + ":grafana-dashboards-provider", "global-grafana:grafana-dashboard"},
	}
	return uc.createEssentialRelations(ctx, scope, endpointList)
}

func (uc *OrchestratorUseCase) consumeRemoteOffer(ctx context.Context, scope, url string) error {
	consumeDetails, err := uc.facilityOffers.GetConsumeDetails(ctx, url)
	if err != nil {
		return err
	}
	offerURL, err := crossmodel.ParseOfferURL(consumeDetails.Offer.OfferURL)
	if err != nil {
		return err
	}
	offerURL.Source = uc.conf.Juju.Controller
	consumeDetails.Offer.OfferURL = offerURL.String()

	args := &crossmodel.ConsumeApplicationArgs{
		Offer:    *consumeDetails.Offer,
		Macaroon: consumeDetails.Macaroon,
	}
	if consumeDetails.ControllerInfo != nil {
		controllerTag, err := names.ParseControllerTag(consumeDetails.ControllerInfo.ControllerTag)
		if err != nil {
			return err
		}
		args.ControllerInfo = &crossmodel.ControllerInfo{
			ControllerTag: controllerTag,
			Alias:         consumeDetails.ControllerInfo.Alias,
			Addrs:         consumeDetails.ControllerInfo.Addrs,
			CACert:        consumeDetails.ControllerInfo.CACert,
		}
	}
	return uc.facility.Consume(ctx, scope, args)
}
