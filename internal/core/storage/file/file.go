package file

import (
	"github.com/otterscale/otterscale/internal/core/facility"
	"github.com/otterscale/otterscale/internal/core/facility/action"
)

type FileUseCase struct {
	volume            VolumeRepo
	subvolumeGroup    SubvolumeGroupRepo
	subvolume         SubvolumeRepo
	subvolumeSnapshot SubvolumeSnapshotRepo

	action   action.ActionRepo
	facility facility.FacilityRepo
}

func NewFileUseCase(volume VolumeRepo, subvolumeGroup SubvolumeGroupRepo, subvolume SubvolumeRepo, subvolumeSnapshot SubvolumeSnapshotRepo, action action.ActionRepo, facility facility.FacilityRepo) *FileUseCase {
	return &FileUseCase{
		volume:            volume,
		subvolumeGroup:    subvolumeGroup,
		subvolume:         subvolume,
		subvolumeSnapshot: subvolumeSnapshot,
		action:            action,
		facility:          facility,
	}
}
