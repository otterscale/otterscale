package core

import (
	"context"
	"errors"

	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
)

type LargeLanguageModelUseCase struct {
	release  ReleaseRepo
	action   ActionRepo
	facility FacilityRepo
}

func NewLargeLanguageModelUseCase(release ReleaseRepo, action ActionRepo, facility FacilityRepo) *LargeLanguageModelUseCase {
	return &LargeLanguageModelUseCase{
		release:  release,
		action:   action,
		facility: facility,
	}
}

func (uc *LargeLanguageModelUseCase) CheckInfrastructureStatus(ctx context.Context, uuid, facility string) (int32, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return 0, err
	}
	rel, err := uc.release.Get(config, LLMd, LLMd)
	if err != nil {
		if errors.Is(err, driver.ErrReleaseNotFound) {
			return llmdInfraHealthNotInstalled, nil
		}
		return 0, err
	}
	switch {
	case rel.Info.Status.IsPending():
		return llmdInfraHealthPending, nil
	case rel.Info.Status == release.StatusDeployed:
		return llmdInfraHealthOK, nil
	case rel.Info.Status == release.StatusFailed:
		return llmdInfraHealthFailed, nil
	}
	return 0, nil
}
