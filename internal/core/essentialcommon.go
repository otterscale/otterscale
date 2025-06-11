package core

import (
	"context"

	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/names/v5"

	"github.com/openhdc/otterscale/internal/config"
)

var (
	commonCharms = []EssentialCharm{
		{Name: "ch:ceph-csi", LXD: true},
		{Name: "ch:grafana-agent", LXD: false},
		{Name: "ch:hardware-observer", LXD: false},
	}

	commonRelations = [][]string{
		{"ceph-csi", "ceph-mon"},
		{"ceph-csi", "kubernetes-control-plane"},
		{"grafana-agent:cos-agent", "ceph-mon:cos-agent"},
		{"grafana-agent:cos-agent", "kubeapi-load-balancer:cos-agent"},
		{"grafana-agent:cos-agent", "kubernetes-control-plane:cos-agent"},
		{"grafana-agent:cos-agent", "kubernetes-worker:cos-agent"},
		{"grafana-agent:cos-agent", "hardware-observer:cos-agent"},
	}
)

func CreateCommon(ctx context.Context, serverRepo ServerRepo, machineRepo MachineRepo, facilityRepo FacilityRepo, facilityOffersRepo FacilityOffersRepo, conf *config.Config, uuid, prefix string, configs map[string]string) error {
	if err := createEssential(ctx, serverRepo, machineRepo, facilityRepo, uuid, "", prefix, commonCharms, configs); err != nil {
		return err
	}
	if err := createCOS(ctx, facilityRepo, facilityOffersRepo, conf, uuid, prefix); err != nil {
		return err
	}
	return createEssentialRelations(ctx, facilityRepo, uuid, toEndpointList(prefix, commonRelations))
}

func newCommonConfigs(prefix string) (map[string]string, error) {
	configs := map[string]map[string]any{
		"ceph-csi": {
			"default-storage":      "ceph-ext4",
			"provisioner-replicas": 1,
		},
	}
	return NewCharmConfigs(prefix, configs)
}

func createCOS(ctx context.Context, facilityRepo FacilityRepo, facilityOffersRepo FacilityOffersRepo, conf *config.Config, uuid, prefix string) error {
	// consume
	offerURLs := []string{
		conf.Juju.Username + "/cos.global-prometheus",
		conf.Juju.Username + "/cos.global-grafana",
	}
	for _, url := range offerURLs {
		if err := consumeRemoteOffer(ctx, facilityRepo, facilityOffersRepo, uuid, url, conf.Juju.Controller); err != nil {
			return err
		}
	}

	// integrate
	name := toEssentialName(prefix, "grafana-agent")
	endpointList := [][]string{
		{name + ":send-remote-write", "global-prometheus:receive-remote-write"},
		{name + ":grafana-dashboards-provider", "global-grafana:grafana-dashboard"},
	}
	return createEssentialRelations(ctx, facilityRepo, uuid, endpointList)
}

func consumeRemoteOffer(ctx context.Context, facilityRepo FacilityRepo, facilityOffersRepo FacilityOffersRepo, uuid, url, controller string) error {
	consumeDetails, err := facilityOffersRepo.GetConsumeDetails(ctx, url)
	if err != nil {
		return err
	}
	offerURL, err := crossmodel.ParseOfferURL(consumeDetails.Offer.OfferURL)
	if err != nil {
		return err
	}
	offerURL.Source = controller
	consumeDetails.Offer.OfferURL = offerURL.String()

	args := crossmodel.ConsumeApplicationArgs{
		Offer:    *consumeDetails.Offer,
		Macaroon: consumeDetails.Macaroon,
	}
	if consumeDetails.ControllerInfo != nil {
		controllerTag, err := names.ParseControllerTag(consumeDetails.ControllerInfo.ControllerTag)
		if err != nil {
			panic(err)
		}
		args.ControllerInfo = &crossmodel.ControllerInfo{
			ControllerTag: controllerTag,
			Alias:         consumeDetails.ControllerInfo.Alias,
			Addrs:         consumeDetails.ControllerInfo.Addrs,
			CACert:        consumeDetails.ControllerInfo.CACert,
		}
	}
	return facilityRepo.Consume(ctx, uuid, args)
}
