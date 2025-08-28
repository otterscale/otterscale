package core

import (
	"context"

	"k8s.io/client-go/rest"
)

type KubeVirtDVRepo interface {
	CreateDataVolume(ctx context.Context, config *rest.Config, namespace, name string, source_type string, source string, sizeBytes int64, is_bootable bool) (*DataVolume, error)
	GetDataVolume(ctx context.Context, config *rest.Config, namespace, name string) (*DataVolume, error)
	ListDataVolume(ctx context.Context, config *rest.Config, namespace string) ([]DataVolume, error)
	DeleteDataVolume(ctx context.Context, config *rest.Config, namespace, name string) error
	ExtendDataVolume(ctx context.Context, config *rest.Config, namespace, name string, sizeBytes int64) error
}

// Data Volume Operations
func (uc *KubeVirtUseCase) CreateDataVolume(ctx context.Context, uuid, facility, namespace string, name string, source_type string, source string, sizeBytes int64, is_bootable bool) (*DataVolume, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	return uc.kubeVirtDV.CreateDataVolume(ctx, config, namespace, name, source_type, source, sizeBytes, is_bootable)
}

func (uc *KubeVirtUseCase) GetDataVolume(ctx context.Context, uuid, facility, namespace string, name string) (*DataVolume, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtDV.GetDataVolume(ctx, config, namespace, name)
}

func (uc *KubeVirtUseCase) ListDataVolumes(ctx context.Context, uuid, facility, namespace string) ([]DataVolume, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtDV.ListDataVolume(ctx, config, namespace)
}

func (uc *KubeVirtUseCase) DeleteDataVolume(ctx context.Context, uuid, facility, namespace string, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtDV.DeleteDataVolume(ctx, config, namespace, name)
}

func (uc *KubeVirtUseCase) ExtendDataVolume(ctx context.Context, uuid, facility, namespace string, name string, sizeBytes int64) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	var dv *DataVolume
	dv, err = uc.kubeVirtDV.GetDataVolume(ctx, config, namespace, name)
	if err != nil {
		return err
	}

	pvcName := dv.Status.ClaimName
	if pvcName == "" {
		pvcName = dv.Name // fallback（通常PVC同名）
	}

	return uc.kubeVirtDV.ExtendDataVolume(ctx, config, namespace, pvcName, sizeBytes)
}

// Extracts source, sourceType, and sizeBytes from a DataVolume
func ExtractDataVolumeInfo(dv *DataVolume) (source string, sourceType string, sizeBytes int64, accessMode string, storageClassName string) {
	if dv.Spec.PVC != nil {
		if dv.Spec.PVC.Resources.Requests != nil {
			accessMode = string(dv.Spec.PVC.AccessModes[0])
			storageClassName = *dv.Spec.PVC.StorageClassName
			size, found := dv.Spec.PVC.Resources.Requests["storage"]
			if found {
				sizeBytes = size.Value()
			}
		}
	} else if dv.Spec.Storage != nil {
		if dv.Spec.Storage.Resources.Requests != nil {
			accessMode = string(dv.Spec.Storage.AccessModes[0])
			storageClassName = *dv.Spec.Storage.StorageClassName
			size, found := dv.Spec.Storage.Resources.Requests["storage"]
			if found {
				sizeBytes = size.Value()
			}
		}
	}

	if dv.Spec.Source.HTTP != nil {
		source = dv.Spec.Source.HTTP.URL
		sourceType = "HTTP"
	} else if dv.Spec.Source.Upload != nil {
		source = ""
		sourceType = "Upload"
	} else if dv.Spec.Source.S3 != nil {
		source = dv.Spec.Source.S3.URL
		sourceType = "S3"
	} else if dv.Spec.Source.VDDK != nil {
		source = dv.Spec.Source.VDDK.URL
		sourceType = string(dv.Spec.Source.VDDK.UUID)
	}

	return
}
