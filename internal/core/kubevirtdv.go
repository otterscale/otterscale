package core

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	v1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
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

func GetDataVolumeConditions(dv *DataVolume) (condition_message string, condition_reason string, condition_status string) {
	for _, condition := range dv.Status.Conditions {
		if condition.Type == "Bound" {
			return condition.Message, condition.Reason, string(condition.Status)
		}
	}
	return
}

func GetPVCInfo(pvc *v1.PersistentVolumeClaimSpec) (sizeBytes int64, accessMode string, storageClassName string) {
	if pvc.Resources.Requests != nil {
		size, found := pvc.Resources.Requests["storage"]
		if found {
			sizeBytes = size.Value()
		}
	}
	if pvc.AccessModes != nil {
		accessMode = string(pvc.AccessModes[0])
	}
	if pvc.StorageClassName != nil {
		storageClassName = *pvc.StorageClassName
	}
	return
}

func GetStorageInfo(storage *v1beta1.StorageSpec) (sizeBytes int64, accessMode string, storageClassName string) {
	if storage.Resources.Requests != nil {
		size, found := storage.Resources.Requests["storage"]
		if found {
			sizeBytes = size.Value()
		}
	}
	if storage.AccessModes != nil {
		accessMode = string(storage.AccessModes[0])
	}
	if storage.StorageClassName != nil {
		storageClassName = *storage.StorageClassName
	}
	return
}

// Extracts source, sourceType, and sizeBytes from a DataVolume
func ExtractDataVolumeInfo(dv *DataVolume) (source string, sourceType string, sizeBytes int64, accessMode string, storageClassName string) {
	if dv == nil {
		return
	}

	switch {
	case dv.Spec.PVC != nil:
		sizeBytes, accessMode, storageClassName = GetPVCInfo(dv.Spec.PVC)
	case dv.Spec.Storage != nil:
		sizeBytes, accessMode, storageClassName = GetStorageInfo(dv.Spec.Storage)
	}

	switch {
	case dv.Spec.Source.HTTP != nil:
		source = dv.Spec.Source.HTTP.URL
		sourceType = "HTTP"
	case dv.Spec.Source.Upload != nil:
		source = ""
		sourceType = "Upload"
	case dv.Spec.Source.S3 != nil:
		source = dv.Spec.Source.S3.URL
		sourceType = "S3"
	case dv.Spec.Source.VDDK != nil:
		source = dv.Spec.Source.VDDK.URL
		sourceType = string(dv.Spec.Source.VDDK.UUID)
	}
	return
}
