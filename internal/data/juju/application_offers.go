package juju

import (
	"context"

	api "github.com/juju/juju/api/client/applicationoffers"
	"github.com/juju/juju/rpc/params"

	"github.com/otterscale/otterscale/internal/core"
)

type applicationOffers struct {
	juju *Juju
}

func NewApplicationOffers(juju *Juju) core.FacilityOffersRepo {
	return &applicationOffers{
		juju: juju,
	}
}

var _ core.FacilityOffersRepo = (*applicationOffers)(nil)

func (r *applicationOffers) GetConsumeDetails(_ context.Context, url string) (params.ConsumeOfferDetails, error) {
	conn, err := r.juju.connection("")
	if err != nil {
		return params.ConsumeOfferDetails{}, err
	}
	return api.NewClient(conn).GetConsumeDetails(url)
}
