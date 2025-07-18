package core

import (
	"context"
)

type ApplicationUseCase struct {
	kubeApps    KubeAppsRepo
	kubeCore    KubeCoreRepo
	kubeStorage KubeStorageRepo
	chart       ChartRepo
	release     ReleaseRepo
	action      ActionRepo
	facility    FacilityRepo
	scope       ScopeRepo
	client      ClientRepo
}

func NewApplicationUseCase(kubeApps KubeAppsRepo, kubeCore KubeCoreRepo, kubeStorage KubeStorageRepo, chart ChartRepo, release ReleaseRepo, action ActionRepo, facility FacilityRepo, scope ScopeRepo, client ClientRepo) *ApplicationUseCase {
	return &ApplicationUseCase{
		kubeApps:    kubeApps,
		kubeCore:    kubeCore,
		kubeStorage: kubeStorage,
		chart:       chart,
		release:     release,
		action:      action,
		facility:    facility,
		scope:       scope,
		client:      client,
	}
}

func (uc *ApplicationUseCase) GetPublicAddress(ctx context.Context, uuid, facility string) (string, error) {
	leader, err := uc.facility.GetLeader(ctx, uuid, facility)
	if err != nil {
		return "", err
	}
	unitInfo, err := uc.facility.GetUnitInfo(ctx, uuid, leader)
	if err != nil {
		return "", err
	}
	return unitInfo.PublicAddress, nil
}
