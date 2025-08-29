package core

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	v1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

// KubeVirtDVRepo 定義 DataVolume 相關的資料層介面
type KubeVirtDVRepo interface {
	// CreateDataVolume 建立 DataVolume
	CreateDataVolume(ctx context.Context, config *rest.Config, namespace, name string,
		source_type string, source string, sizeBytes int64, is_bootable bool) (*DataVolume, error)

	// GetDataVolume 取得指定的 DataVolume
	GetDataVolume(ctx context.Context, config *rest.Config, namespace, name string) (*DataVolume, error)

	// ListDataVolume 列出指定 namespace 中的所有 DataVolume
	ListDataVolume(ctx context.Context, config *rest.Config, namespace string) ([]DataVolume, error)

	// DeleteDataVolume 刪除指定的 DataVolume
	DeleteDataVolume(ctx context.Context, config *rest.Config, namespace, name string) error

	// ExtendDataVolume 為已有的 PVC（同步 DataVolume）擴容
	ExtendDataVolume(ctx context.Context, config *rest.Config, namespace, name string, sizeBytes int64) error
}

// CreateDataVolume 透過 KubeVirtDVRepo 建立 DataVolume，先取得 kubeConfig 再委派給 repo
func (uc *KubeVirtUseCase) CreateDataVolume(ctx context.Context, uuid, facility,
	namespace string, name string, source_type string, source string, sizeBytes int64, is_bootable bool) (*DataVolume, error) {

	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	return uc.kubeVirtDV.CreateDataVolume(ctx, config, namespace, name,
		source_type, source, sizeBytes, is_bootable)
}

// GetDataVolume 取得 DataVolume 並回傳 domain model
func (uc *KubeVirtUseCase) GetDataVolume(ctx context.Context, uuid, facility,
	namespace string, name string) (*DataVolume, error) {

	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtDV.GetDataVolume(ctx, config, namespace, name)
}

// ListDataVolumes 列出指定 namespace 下的所有 DataVolume
func (uc *KubeVirtUseCase) ListDataVolumes(ctx context.Context, uuid, facility,
	namespace string) ([]DataVolume, error) {

	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtDV.ListDataVolume(ctx, config, namespace)
}

// DeleteDataVolume 刪除指定的 DataVolume
func (uc *KubeVirtUseCase) DeleteDataVolume(ctx context.Context, uuid, facility,
	namespace string, name string) error {

	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtDV.DeleteDataVolume(ctx, config, namespace, name)
}

// ExtendDataVolume 為 DataVolume 內的 PVC 進行容量擴充
func (uc *KubeVirtUseCase) ExtendDataVolume(ctx context.Context, uuid, facility,
	namespace string, name string, sizeBytes int64) error {

	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	// 先取得 DataVolume，找出實際的 PVC 名稱
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

// GetDataVolumeConditions 取得 DataVolume 的 Bound 條件訊息
func GetDataVolumeConditions(dv *DataVolume) (condition_message string,
	condition_reason string, condition_status string) {

	for _, condition := range dv.Status.Conditions {
		if condition.Type == "Bound" {
			return condition.Message, condition.Reason, string(condition.Status)
		}
	}
	return
}

// GetPVCInfo 從 PVC Spec 取得大小、存取模式與 storage class
func GetPVCInfo(pvc *v1.PersistentVolumeClaimSpec) (sizeBytes int64,
	accessMode string, storageClassName string) {

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

// GetStorageInfo 從 DataVolume.Storage 取得大小、存取模式與 storage class
func GetStorageInfo(storage *v1beta1.StorageSpec) (sizeBytes int64,
	accessMode string, storageClassName string) {

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
ExtractDataVolumeInfo 從 DataVolume 取得 source、sourceType、大小、存取模式與 storage class。
DataVolume 同時可能有 PVC 與 Storage 兩種規格，但同時只會出現其中一種。
*/
func ExtractDataVolumeInfo(dv *DataVolume) (source string, sourceType string, sizeBytes int64, accessMode string, storageClassName string) {
	if dv == nil {
		return
	}

	// 先根據使用的 volume 類型取得容量資訊
	switch {
	case dv.Spec.PVC != nil:
		sizeBytes, accessMode, storageClassName = GetPVCInfo(dv.Spec.PVC)
	case dv.Spec.Storage != nil:
		sizeBytes, accessMode, storageClassName = GetStorageInfo(dv.Spec.Storage)
	}

	// 再根據 source 類型取得來源 URL（或 UUID）
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
