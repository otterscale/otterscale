package environment

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core/facility/action"
	"github.com/otterscale/otterscale/internal/core/scope"
)

type traefikProxiedEndpoints struct {
	Prometheus struct {
		URL string `json:"url"`
	} `json:"prometheus/0"`
}

type UseCase struct {
	conf *config.Config

	action action.ActionRepo
	scope  scope.ScopeRepo

	prometheusURL *url.URL
}

func NewUseCase(conf *config.Config, action action.ActionRepo, scope scope.ScopeRepo) *UseCase {
	return &UseCase{
		conf:          conf,
		action:        action,
		scope:         scope,
		prometheusURL: &url.URL{},
	}
}

func (uc *UseCase) CheckHealth() (HealthStatus, error) {
	if !isMAASConfigured(uc.conf) {
		return HealthStatusNotInstalled, nil
	}
	return HealthStatusOK, nil
}

// TODO: update kubernetes config map
func (uc *UseCase) UpdateConfig(_ *config.Schema) error {
	return errors.ErrUnsupported
}

func (uc *UseCase) GetPrometheusURL() *url.URL {
	return uc.prometheusURL
}

func (uc *UseCase) DiscoverPrometheusURL(ctx context.Context) (*url.URL, error) {
	if uc.prometheusURL.Scheme != "" {
		return uc.prometheusURL, nil
	}

	result, err := uc.action.Run(ctx, "cos", "traefik", "show-proxied-endpoints", nil)
	if err != nil {
		return nil, err
	}

	var endpoints traefikProxiedEndpoints
	proxiedEndpointsStr, ok := result["proxied-endpoints"].(string)
	if !ok {
		return nil, err
	}

	if err := json.Unmarshal([]byte(proxiedEndpointsStr), &endpoints); err != nil {
		return nil, err
	}

	prometheusURL, err := url.Parse(endpoints.Prometheus.URL)
	if err != nil {
		return nil, err
	}

	*uc.prometheusURL = *prometheusURL
	return uc.prometheusURL, nil
}

func isMAASConfigured(conf *config.Config) bool {
	return conf.MAASKey() != "::"
}
