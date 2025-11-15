package standalone

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func (uc *UseCase) createCOS(ctx context.Context, scope string) error {
	// consume
	offerURLs := []string{
		uc.conf.Juju.Username + "/cos.global-prometheus",
		uc.conf.Juju.Username + "/cos.global-grafana",
	}

	for _, url := range offerURLs {
		if err := uc.relation.Consume(ctx, scope, url); err != nil {
			return err
		}
	}

	// integrate
	name := scope + "-" + "grafana-agent"
	endpointList := [][]string{
		{name + ":send-remote-write", "global-prometheus:receive-remote-write"},
		{name + ":grafana-dashboards-provider", "global-grafana:grafana-dashboard"},
	}

	eg, egctx := errgroup.WithContext(ctx)

	for _, endpoints := range endpointList {
		eg.Go(func() error {
			return uc.relation.Create(egctx, scope, endpoints)
		})
	}

	return eg.Wait()
}
