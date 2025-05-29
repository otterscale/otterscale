package core

import (
	"context"
	"errors"
	"sync"

	"connectrpc.com/connect"

	"github.com/openhdc/otterscale/internal/config"
)

type EnvironmentStatus struct {
	Phase   string
	Message string
}

type EnvironmentUseCase struct {
	conf      *config.Config
	statusMap sync.Map
}

func NewEnvironmentUseCase(conf *config.Config) *EnvironmentUseCase {
	return &EnvironmentUseCase{
		conf: conf,
	}
}

//nolint:mnd
func (uc *EnvironmentUseCase) CheckHealthy(ctx context.Context) (int32, error) {
	if uc.conf.MAAS.Key == "::" {
		return 21, nil // NOT_INSTALLED
	}
	return 11, nil // OK
}

func (uc *EnvironmentUseCase) LoadStatus(ctx context.Context) (*EnvironmentStatus, error) {
	v, ok := uc.statusMap.Load("")
	if ok {
		return v.(*EnvironmentStatus), nil
	}
	return nil, connect.NewError(connect.CodeNotFound, errors.New("status not found"))
}

func (uc *EnvironmentUseCase) StoreStatus(ctx context.Context, phase, message string) {
	uc.statusMap.Store("", &EnvironmentStatus{
		Phase:   phase,
		Message: message,
	})
}

func (uc *EnvironmentUseCase) UpdateConfig(ctx context.Context, conf *config.Config) error {
	return uc.conf.Override(conf)
}

func (uc *EnvironmentUseCase) UpdateConfigHelmRepos(ctx context.Context, urls []string) error {
	uc.conf.Kube.HelmRepositoryURLs = urls
	return uc.conf.Override(uc.conf)
}
