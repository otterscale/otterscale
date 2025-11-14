package juju

import (
	"context"

	"github.com/juju/juju/api/client/application"

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

// func (r *relationRepo) Consume(_ context.Context, scope string, args *crossmodel.ConsumeApplicationArgs) error {
// 	conn, err := r.juju.connection(scope)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = application.NewClient(conn).Consume(*args)
// 	return err
// }

// func (r *applicationOffers) GetConsumeDetails(_ context.Context, url string) (params.ConsumeOfferDetails, error) {
// 	conn, err := r.juju.connection("controller")
// 	if err != nil {
// 		return params.ConsumeOfferDetails{}, err
// 	}
// 	return api.NewClient(conn).GetConsumeDetails(url)
// }
