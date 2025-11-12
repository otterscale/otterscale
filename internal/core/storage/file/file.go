package file

import "github.com/otterscale/otterscale/internal/core/facility"

type FileUseCase struct {
	volume            VolumeRepo
	subvolumeGroup    SubvolumeGroupRepo
	subvolume         SubvolumeRepo
	subvolumeSnapshot SubvolumeSnapshotRepo

	action   facility.ActionRepo
	facility facility.FacilityRepo
}

func NewFileUseCase(volume VolumeRepo, subvolumeGroup SubvolumeGroupRepo, subvolume SubvolumeRepo, subvolumeSnapshot SubvolumeSnapshotRepo, action facility.ActionRepo, facility facility.FacilityRepo) *FileUseCase {
	return &FileUseCase{
		volume:            volume,
		subvolumeGroup:    subvolumeGroup,
		subvolume:         subvolume,
		subvolumeSnapshot: subvolumeSnapshot,
		action:            action,
		facility:          facility,
	}
}
