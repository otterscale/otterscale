package standalone

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func (uc *UseCase) createCOS(ctx context.Context, scope string) error {
	// consume
	username := uc.conf.JujuUsername()
	offerURLs := []string{
		username + "/cos.global-prometheus",
		username + "/cos.global-grafana",
	}

	for _, url := range offerURLs {
		if err := uc.relation.Consume(ctx, scope, url); err != nil {
			return err
		}
	}

	// integrate
	endpointList := [][]string{
		{"grafana-agent:send-remote-write", "global-prometheus:receive-remote-write"},
		{"grafana-agent:grafana-dashboards-provider", "global-grafana:grafana-dashboard"},
	}

	eg, egctx := errgroup.WithContext(ctx)

	for _, endpoints := range endpointList {
		eg.Go(func() error {
			return uc.relation.Create(egctx, scope, endpoints)
		})
	}

	return eg.Wait()
}
