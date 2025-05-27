package core

import (
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	maas := uc.conf.ConfigSet.GetMaas()
	if maas.GetKey() == "::" {
		return 21, nil // NOT_INSTALLED
	}
	return 11, nil // OK
}

func (uc *EnvironmentUseCase) LoadStatus(ctx context.Context) (*EnvironmentStatus, error) {
	v, ok := uc.statusMap.Load("")
	if ok {
		return v.(*EnvironmentStatus), nil
	}
	return nil, status.Error(codes.NotFound, "status not found")
}

func (uc *EnvironmentUseCase) StoreStatus(ctx context.Context, phase, message string) {
	uc.statusMap.Store("", &EnvironmentStatus{
		Phase:   phase,
		Message: message,
	})
}

func (uc *EnvironmentUseCase) UpdateConfig(ctx context.Context, set *config.ConfigSet) error {
	return uc.conf.Override(set)
}

func (uc *EnvironmentUseCase) UpdateConfigHelmRepos(ctx context.Context, urls []string) error {
	kube := uc.conf.ConfigSet.GetKube()
	kube.SetHelmRepositoryUrls(urls)
	return uc.conf.Override(uc.conf.ConfigSet)
}
