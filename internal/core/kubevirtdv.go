package core

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/rest"
	v1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

type KubeVirtDVRepo interface {
	CreateDataVolume(ctx context.Context, config *rest.Config, namespace, name, sourceType, source, vmName string, sizeBytes int64, isBootable bool) (*DataVolume, error)
	GetDataVolume(ctx context.Context, config *rest.Config, namespace, name string, pvc *v1.PersistentVolumeClaim) (*DataVolume, error)
	ListDataVolumes(ctx context.Context, config *rest.Config, namespace string) ([]DataVolume, error)
	ListDataVolumesByOptions(ctx context.Context, config *rest.Config, namespace, label, field string) ([]DataVolume, error)
	DeleteDataVolume(ctx context.Context, config *rest.Config, namespace, name string) error
	ExtendDataVolume(ctx context.Context, config *rest.Config, namespace string, pvc *v1.PersistentVolumeClaim, sizeBytes resource.Quantity) error
}

// CreateDataVolume creates a DataVolume via KubeVirtDVRepo. It obtains kubeConfig first and then delegates to the repo.
func (uc *KubeVirtUseCase) CreateDataVolume(ctx context.Context, uuid, facility, namespace, name, sourceType, source string, sizeBytes int64, isBootable bool) (*DataVolume, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	return uc.kubeVirtDV.CreateDataVolume(ctx, config, namespace, name,
		sourceType, source, "", sizeBytes, isBootable)
}

// GetDataVolume retrieves a DataVolume and returns the domain model.
func (uc *KubeVirtUseCase) GetDataVolume(ctx context.Context, uuid, facility, namespace, name string) (*DataVolume, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	pvc, err := uc.kubeCore.GetPersistentVolumeClaims(ctx, config, namespace, name)
	if err != nil {
		return nil, err
	}

	dv, err := uc.kubeVirtDV.GetDataVolume(ctx, config, namespace, name, pvc)
	if err != nil {
		return nil, err
	}

	return dv, nil
}

// DeleteDataVolume deletes the specified DataVolume.
func (uc *KubeVirtUseCase) ListDataVolumes(ctx context.Context, uuid, facility, namespace string) ([]DataVolume, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtDV.ListDataVolumes(ctx, config, namespace)
}

// ExtendDataVolume expands the capacity of the PVC inside the DataVolume.
func (uc *KubeVirtUseCase) DeleteDataVolume(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtDV.DeleteDataVolume(ctx, config, namespace, name)
}

// ExtendDataVolume expands the capacity of the PVC inside the DataVolume.
func (uc *KubeVirtUseCase) ExtendDataVolume(ctx context.Context, uuid, facility, namespace, name string, sizeBytes int64) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	pvc, err := uc.kubeCore.GetPersistentVolumeClaims(ctx, config, namespace, name)
	if err != nil {
		return err
	}

	current := pvc.Spec.Resources.Requests[v1.ResourceStorage]
	desired := *resource.NewQuantity(sizeBytes, resource.BinarySI)
	// If the current size is greater than or equal to the requested size, return an error directly.
	if current.Cmp(desired) >= 0 {
		return fmt.Errorf("current size >= requested size, no need to extend")
	}

	return uc.kubeVirtDV.ExtendDataVolume(ctx, config, namespace, pvc, desired)
}

// GetDataVolumeConditions retrieves the Bound condition message of a DataVolume.
func GetDataVolumeConditions(dv *DataVolume) (conditionMessage, conditionReason, conditionStatus string) {
	for _, condition := range dv.Status.Conditions {
		if condition.Type == "Bound" {
			return condition.Message, condition.Reason, string(condition.Status)
		}
	}
	return
}

// GetPVCInfo retrieves the size, access mode, and storage class from PVC Spec.
func GetPVCInfo(pvc *v1.PersistentVolumeClaimSpec) (sizeBytes int64, accessMode, storageClassName string) {
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

// GetStorageInfo retrieves the size, access mode, and storage class from DataVolume.Storage.
func GetStorageInfo(storage *v1beta1.StorageSpec) (sizeBytes int64, accessMode, storageClassName string) {
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

/*
ExtractDataVolumeInfo retrieves source, sourceType, size, access mode, and storage class from a DataVolume.
A DataVolume may have PVC or Storage specification, but only one will appear at a time.
*/
func ExtractDataVolumeInfo(dv *DataVolume) (source, sourceType, accessMode, storageClassName string, sizeBytes int64) {
	if dv == nil {
		return
	}

	// First, get capacity info based on the volume type used.
	switch {
	case dv.Spec.PVC != nil:
		sizeBytes, accessMode, storageClassName = GetPVCInfo(dv.Spec.PVC)
	case dv.Spec.Storage != nil:
		sizeBytes, accessMode, storageClassName = GetStorageInfo(dv.Spec.Storage)
	}

	// Then, get the source URL (or UUID) based on the source type.
	switch {
	case dv.Spec.Source.HTTP != nil:
		source = dv.Spec.Source.HTTP.URL
		sourceType = "HTTP"
	case dv.Spec.Source.Blank != nil:
		source = ""
		sourceType = "BLANK"
	case dv.Spec.Source.PVC != nil:
		source = dv.Spec.Source.PVC.Name
		sourceType = "PVC"
	}
	return
}

// StorageSpec helper function for generating StorageSpec (non-PVC)
func StorageSpec(size int64) *v1beta1.StorageSpec {
	return &v1beta1.StorageSpec{
		AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
		Resources: v1.VolumeResourceRequirements{
			Requests: v1.ResourceList{
				v1.ResourceStorage: *resource.NewQuantity(size, resource.BinarySI),
			},
		},
	}
}

// PvcResizePatch generates JSON-Patch for PVC resizing.
func PvcResizePatch(desired string) ([]byte, error) {
	ops := []map[string]interface{}{
		{"op": "replace", "path": "/spec/resources/requests/storage", "value": desired},
		{
			"op": "add", "path": "/metadata/annotations/otterscale.io~1last-updated",
			"value": time.Now().Format(time.RFC3339),
		},
	}
	return json.Marshal(ops)
}

// DataVolumeLastUpdatedPatch adds last-updated annotation to DataVolume.
func DataVolumeLastUpdatedPatch() ([]byte, error) {
	patch := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": map[string]string{
				"otterscale.io/last-updated": time.Now().Format(time.RFC3339),
			},
		},
	}
	return json.Marshal(patch)
}

// SyncDataVolumeSpec synchronizes resource information from PVC to DataVolume.Spec (PVC or Storage).
func SyncDataVolumeSpec(dv *v1beta1.DataVolume, pvc *v1.PersistentVolumeClaim) {
	req := pvc.Spec.Resources.Requests
	if dv.Spec.PVC != nil {
		if dv.Spec.PVC.Resources.Requests == nil {
			dv.Spec.PVC.Resources.Requests = make(map[v1.ResourceName]resource.Quantity)
		}
		dv.Spec.PVC.Resources.Requests = req
		dv.Spec.PVC.StorageClassName = pvc.Spec.StorageClassName
	}
	if dv.Spec.Storage != nil {
		if dv.Spec.Storage.Resources.Requests == nil {
			dv.Spec.Storage.Resources.Requests = make(map[v1.ResourceName]resource.Quantity)
		}
		dv.Spec.Storage.Resources.Requests = req
		dv.Spec.Storage.StorageClassName = pvc.Spec.StorageClassName
	}
}
