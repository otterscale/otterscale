package core

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"slices"

	"connectrpc.com/connect"

	"github.com/otterscale/otterscale/internal/config"
)

type traefikProxiedEndpoints struct {
	Prometheus struct {
		URL string `json:"url"`
	} `json:"prometheus/0"`
}

type EnvironmentUseCase struct {
	conf *config.Config

	action   ActionRepo
	facility FacilityRepo
	scope    ScopeRepo

	prometheusURL *url.URL
}

func NewEnvironmentUseCase(conf *config.Config, action ActionRepo, facility FacilityRepo, scope ScopeRepo) *EnvironmentUseCase {
	return &EnvironmentUseCase{
		conf:          conf,
		action:        action,
		scope:         scope,
		facility:      facility,
		prometheusURL: &url.URL{},
	}
}

func (uc *EnvironmentUseCase) CheckHealth(_ context.Context) (int32, error) {
	if !isMAASConfigured(uc.conf) {
		return environmentHealthNotInstalled, nil
	}
	return environmentHealthOK, nil
}

func (uc *EnvironmentUseCase) UpdateConfig(_ context.Context, conf *config.Config) error {
	uc.conf.MAAS = conf.MAAS
	uc.conf.Juju = conf.Juju
	uc.conf.MicroK8s = conf.MicroK8s
	return uc.conf.Override(uc.conf)
}

func (uc *EnvironmentUseCase) GetConfigHelmRepos() []string {
	return uc.conf.Kube.HelmRepositoryURLs
}

func (uc *EnvironmentUseCase) UpdateConfigHelmRepos(urls []string) error {
	uc.conf.Kube.HelmRepositoryURLs = urls
	return uc.conf.Override(uc.conf)
}

func (uc *EnvironmentUseCase) GetPrometheusURL() *url.URL {
	return uc.prometheusURL
}

func (uc *EnvironmentUseCase) DiscoverPrometheusURL(ctx context.Context) (*url.URL, error) {
	if uc.prometheusURL.Scheme != "" {
		return uc.prometheusURL, nil
	}

	scopes, err := uc.scope.List(ctx)
	if err != nil {
		return nil, err
	}

	cosNames := []string{"cos", "cos-lite", "cos-dev"}
	cosScopes := slices.DeleteFunc(scopes, func(s Scope) bool {
		return !slices.Contains(cosNames, s.Name)
	})

	for i := range cosScopes {
		leader, err := uc.facility.GetLeader(ctx, cosScopes[i].Name, "traefik")
		if err != nil {
			continue
		}

		result, err := runAction(ctx, uc.action, cosScopes[i].Name, leader, "show-proxied-endpoints", nil)
		if err != nil {
			continue
		}

		var endpoints traefikProxiedEndpoints
		proxiedEndpointsStr, ok := result.Output["proxied-endpoints"].(string)
		if !ok {
			continue
		}

		if err := json.Unmarshal([]byte(proxiedEndpointsStr), &endpoints); err != nil {
			continue
		}

		prometheusURL, err := url.Parse(endpoints.Prometheus.URL)
		if err != nil {
			continue
		}

		*uc.prometheusURL = *prometheusURL
		return uc.prometheusURL, nil
	}

	return nil, connect.NewError(connect.CodeNotFound, errors.New("prometheus URL not found in any COS scope"))
}

func isMAASConfigured(conf *config.Config) bool {
	return conf.MAAS.Key != "::"
}
