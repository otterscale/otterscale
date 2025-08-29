package kube

import (
	"context"
	"fmt"
	"strconv"

	oscore "github.com/openhdc/otterscale/internal/core"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

const (
	SourceHTTP     = "HTTP"
	SourceBlank    = "Blank"
	SourceRegistry = "Registry"
	SourceUpload   = "Upload"
	SourceS3       = "S3"
	SourceVDDK     = "VDDK"
	SourcePVC      = "PVC"
)

type virtDV struct {
	kubevirt *kubevirt
}

func NewVirtDV(kube *Kube, kubevirt *kubevirt) oscore.KubeVirtDVRepo {
	return &virtDV{
		kubevirt: kubevirt,
	}
}

var _ oscore.KubeVirtDVRepo = (*virtDV)(nil)

// -------------------------------------------------------------------
// Repository implementation (CRUD + Extend)
// -------------------------------------------------------------------

// CreateDataVolume 建立 DataVolume，根據 source_type 建構不同的 DataVolumeSource
// 並依需求設定 PVC 或 Storage Spec，最後呼叫 CDI API 建立。
func (r *virtDV) CreateDataVolume(ctx context.Context, config *rest.Config, namespace, name string, source_type string, source string, sizeBytes int64, is_bootable bool) (*oscore.DataVolume, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	dvSpec := &v1beta1.DataVolumeSpec{}
	var dvSource *v1beta1.DataVolumeSource

	// 依 source_type 決定 DataVolumeSource
	switch source_type {
	case SourceHTTP:
		dvSource = &v1beta1.DataVolumeSource{HTTP: &v1beta1.DataVolumeSourceHTTP{URL: source}}
	case SourceBlank:
		dvSource = &v1beta1.DataVolumeSource{Blank: &v1beta1.DataVolumeBlankImage{}}
	case SourceRegistry:
		dvSource = &v1beta1.DataVolumeSource{Registry: &v1beta1.DataVolumeSourceRegistry{URL: &source}}
	case SourceUpload:
		dvSource = &v1beta1.DataVolumeSource{Upload: &v1beta1.DataVolumeSourceUpload{}}
	case SourceS3:
		dvSource = &v1beta1.DataVolumeSource{S3: &v1beta1.DataVolumeSourceS3{URL: source}}
	case SourceVDDK:
		dvSource = &v1beta1.DataVolumeSource{VDDK: &v1beta1.DataVolumeSourceVDDK{URL: source}}
	case SourcePVC:
		// 直接使用已有 PVC，另外建立 PVC spec 供 DataVolume 參考
		dvSource = &v1beta1.DataVolumeSource{
			PVC: &v1beta1.DataVolumeSourcePVC{Namespace: namespace, Name: source},
		}
		pvcSpec := &v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			Resources: v1.VolumeResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: *resource.NewQuantity(sizeBytes, resource.BinarySI),
				},
			},
		}
		dvSpec.PVC = pvcSpec
	default:
		// 未支援的 source_type 直接回傳錯誤
		return nil, fmt.Errorf("unsupported source_type: %s", source_type)
	}

	// 若不是 PVC，使用 StorageSpec 填入大小資訊
	if source_type != "PVC" {
		dvSpec.Storage = oscore.StorageSpec(sizeBytes)
	}
	dvSpec.Source = dvSource

	// 建立 DataVolume 物件
	dv := &v1beta1.DataVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"is_bootable": strconv.FormatBool(is_bootable),
			},
		},
		Spec: *dvSpec,
	}

	opts := metav1.CreateOptions{}
	return virtClient.CdiClient().CdiV1beta1().DataVolumes(namespace).Create(ctx, dv, opts)
}

// GetDataVolume 取得指定的 DataVolume，並同步其 PVC 的資源與 storageClass 訊息。
func (r *virtDV) GetDataVolume(ctx context.Context, config *rest.Config, namespace string, name string, pvc *v1.PersistentVolumeClaim) (*oscore.DataVolume, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	dv, err := virtClient.CdiClient().CdiV1beta1().DataVolumes(namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}

	oscore.SyncDataVolumeSpec(dv, pvc)
	return dv, nil
}

// ListDataVolume 列出指定 namespace 中的所有 DataVolume。
func (r *virtDV) ListDataVolume(ctx context.Context, config *rest.Config, namespace string) ([]oscore.DataVolume, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	dvs, err := virtClient.CdiClient().CdiV1beta1().DataVolumes(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return dvs.Items, nil
}

// DeleteDataVolume 刪除指定的 DataVolume。
func (r *virtDV) DeleteDataVolume(ctx context.Context, config *rest.Config, namespace, name string) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}
	opts := metav1.DeleteOptions{}
	return virtClient.CdiClient().CdiV1beta1().DataVolumes(namespace).Delete(ctx, name, opts)
}

// ExtendDataVolume 為既有的 PVC 以及對應的 DataVolume 進行擴容。
func (r *virtDV) ExtendDataVolume(ctx context.Context, config *rest.Config, namespace string,
	pvc *v1.PersistentVolumeClaim, sizeBytes resource.Quantity) error {

	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}

	// --- PVC Patch ---------------------------------------------------------
	patchBytes, err := oscore.PvcResizePatch(sizeBytes.String())
	if err != nil {
		return fmt.Errorf("marshal pvc patch: %w", err)
	}
	_, err = virtClient.CoreV1().PersistentVolumeClaims(namespace).
		Patch(ctx, pvc.Name, types.JSONPatchType, patchBytes, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("patch PVC failed: %w", err)
	}

	// --- DataVolume Annotation Patch ----------------------------------------
	annoBytes, err := oscore.DataVolumeLastUpdatedPatch()
	if err != nil {
		return fmt.Errorf("marshal dv patch: %w", err)
	}
	_, err = virtClient.CdiClient().
		CdiV1beta1().
		DataVolumes(namespace).
		Patch(ctx, pvc.Name, types.MergePatchType, annoBytes, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("patch DataVolume annotation failed: %w", err)
	}

	return nil
}
