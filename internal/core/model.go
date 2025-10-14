package core

type ModelUseCase struct {
	action   ActionRepo
	facility FacilityRepo
	release  ReleaseRepo
}

func NewModelUseCase(action ActionRepo, facility FacilityRepo, release ReleaseRepo) *ModelUseCase {
	return &ModelUseCase{
		action:   action,
		facility: facility,
		release:  release,
	}
}

// TODO: add back when llm infra is supported
// func (uc *ModelUseCase) CheckInfrastructureStatus(ctx context.Context, uuid, facility string) (int32, error) {
// 	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
// 	if err != nil {
// 		return 0, err
// 	}
// 	rel, err := uc.release.Get(config, LLMd, LLMd)
// 	if err != nil {
// 		if errors.Is(err, driver.ErrReleaseNotFound) {
// 			return llmdInfraHealthNotInstalled, nil
// 		}
// 		return 0, err
// 	}
// 	switch {
// 	case rel.Info.Status.IsPending():
// 		return llmdInfraHealthPending, nil
// 	case rel.Info.Status == release.StatusDeployed:
// 		return llmdInfraHealthOK, nil
// 	case rel.Info.Status == release.StatusFailed:
// 		return llmdInfraHealthFailed, nil
// 	}
// 	return 0, nil
// }
