package file

import (
	"github.com/otterscale/otterscale/internal/core/facility"
	"github.com/otterscale/otterscale/internal/core/facility/action"
)

type UseCase struct {
	volume            VolumeRepo
	subvolumeGroup    SubvolumeGroupRepo
	subvolume         SubvolumeRepo
	subvolumeSnapshot SubvolumeSnapshotRepo

	action   action.ActionRepo
	facility facility.FacilityRepo
}

func NewUseCase(volume VolumeRepo, subvolumeGroup SubvolumeGroupRepo, subvolume SubvolumeRepo, subvolumeSnapshot SubvolumeSnapshotRepo, action action.ActionRepo, facility facility.FacilityRepo) *UseCase {
	return &UseCase{
		volume:            volume,
		subvolumeGroup:    subvolumeGroup,
		subvolume:         subvolume,
		subvolumeSnapshot: subvolumeSnapshot,
		action:            action,
		facility:          facility,
	}
}
