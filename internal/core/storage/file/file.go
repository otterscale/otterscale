package file

import (
	"github.com/otterscale/otterscale/internal/core/facility"
	"github.com/otterscale/otterscale/internal/core/facility/action"
)

type UseCase struct {
	subvolumeGroup    SubvolumeGroupRepo
	subvolume         SubvolumeRepo
	subvolumeSnapshot SubvolumeSnapshotRepo
	volume            VolumeRepo

	action   action.ActionRepo
	facility facility.FacilityRepo
}

func NewUseCase(subvolumeGroup SubvolumeGroupRepo, subvolume SubvolumeRepo, subvolumeSnapshot SubvolumeSnapshotRepo, volume VolumeRepo, action action.ActionRepo, facility facility.FacilityRepo) *UseCase {
	return &UseCase{
		subvolumeGroup:    subvolumeGroup,
		subvolume:         subvolume,
		subvolumeSnapshot: subvolumeSnapshot,
		volume:            volume,
		action:            action,
		facility:          facility,
	}
}
