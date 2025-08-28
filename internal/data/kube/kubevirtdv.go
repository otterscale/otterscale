package kube

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	oscore "github.com/openhdc/otterscale/internal/core"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
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

func (r *virtDV) CreateDataVolume(ctx context.Context, config *rest.Config, namespace, name string, source_type string, source string, sizeBytes int64, is_bootable bool) (*oscore.DataVolume, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	newStorageSpec := func(size int64) *v1beta1.StorageSpec {
		return &v1beta1.StorageSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteOnce,
			},
			Resources: v1.VolumeResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: *resource.NewQuantity(size, resource.BinarySI),
				},
			},
		}
	}

	dvSpec := &v1beta1.DataVolumeSpec{}
	var dvSource *v1beta1.DataVolumeSource

	switch source_type {
	case "HTTP":
		dvSource = &v1beta1.DataVolumeSource{HTTP: &v1beta1.DataVolumeSourceHTTP{URL: source}}
	case "Blank":
		dvSource = &v1beta1.DataVolumeSource{Blank: &v1beta1.DataVolumeBlankImage{}}
	case "Registry":
		dvSource = &v1beta1.DataVolumeSource{Registry: &v1beta1.DataVolumeSourceRegistry{URL: &source}}
	case "Upload":
		dvSource = &v1beta1.DataVolumeSource{Upload: &v1beta1.DataVolumeSourceUpload{}}
	case "S3":
		dvSource = &v1beta1.DataVolumeSource{S3: &v1beta1.DataVolumeSourceS3{URL: source}}
	case "VDDK":
		dvSource = &v1beta1.DataVolumeSource{VDDK: &v1beta1.DataVolumeSourceVDDK{URL: source}}
	case "PVC":
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
		return nil, err
	}

	if source_type != "PVC" {
		dvSpec.Storage = newStorageSpec(sizeBytes)
	}
	dvSpec.Source = dvSource

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

func (r *virtDV) GetDataVolume(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.DataVolume, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	dv, err := virtClient.CdiClient().CdiV1beta1().DataVolumes(namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}

	// Get PVC
	pvc, err := virtClient.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, dv.Status.ClaimName, opts)
	if err != nil {
		return nil, err
	}

	pvcRequests := pvc.Spec.Resources.Requests
	// ── 替換 DataVolume 中可能的兩個欄位 ─────────────────────────────
	// dv.Spec.PVC.Resources.Requests
	if dv.Spec.PVC != nil {
		if dv.Spec.PVC.Resources.Requests == nil {
			dv.Spec.PVC.Resources.Requests = make(map[v1.ResourceName]resource.Quantity)
		}
		dv.Spec.PVC.Resources.Requests = pvcRequests
		dv.Spec.PVC.StorageClassName = pvc.Spec.StorageClassName
	}

	// dv.Spec.Storage.Resources.Requests
	if dv.Spec.Storage != nil {
		if dv.Spec.Storage.Resources.Requests == nil {
			dv.Spec.Storage.Resources.Requests = make(map[v1.ResourceName]resource.Quantity)
		}
		dv.Spec.Storage.Resources.Requests = pvcRequests
		dv.Spec.Storage.StorageClassName = pvc.Spec.StorageClassName
	}

	return dv, nil
}

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

func (r *virtDV) DeleteDataVolume(ctx context.Context, config *rest.Config, namespace, name string) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}
	opts := metav1.DeleteOptions{}
	return virtClient.CdiClient().CdiV1beta1().DataVolumes(namespace).Delete(ctx, name, opts)
}

func (r *virtDV) ExtendDataVolume(ctx context.Context, config *rest.Config, namespace, pvcName string, sizeBytes int64) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}
	opts := metav1.GetOptions{}

	// Get PVC
	pvc, err := virtClient.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, opts)
	if err != nil {
		return fmt.Errorf("failed to get PVC: %w", err)
	}

	current := pvc.Spec.Resources.Requests[v1.ResourceStorage]
	desired := *resource.NewQuantity(sizeBytes, resource.BinarySI)
	if current.Cmp(desired) >= 0 {
		return fmt.Errorf("current size >= requested size, no need to extend")
	}

	// --- PVC Patch ---------------------------------------------------------
	patchOps := []map[string]interface{}{
		{
			"op":    "replace",
			"path":  "/spec/resources/requests/storage",
			"value": desired.String(), // resource.Quantity 會以字串形式送出
		},
		{
			"op":    "add",
			"path":  "/metadata/annotations/otterscale.io~1last-updated",
			"value": time.Now().Format(time.RFC3339),
		},
	}
	patchBytes, err := json.Marshal(patchOps)
	if err != nil {
		return fmt.Errorf("failed to marshal pvc patch: %w", err)
	}

	_, err = virtClient.CoreV1().PersistentVolumeClaims(namespace).
		Patch(ctx, pvcName, types.JSONPatchType, patchBytes, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("patch PVC failed: %w", err)
	}

	// --- DataVolume Annotation Patch ----------------------------------------
	annotationPatch := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": map[string]string{
				"otterscale.io/last-updated": time.Now().Format(time.RFC3339),
			},
		},
	}
	annoBytes, _ := json.Marshal(annotationPatch)

	_, err = virtClient.CdiClient().
		CdiV1beta1().
		DataVolumes(namespace).
		Patch(ctx, pvcName, types.MergePatchType, annoBytes, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("patch DataVolume annotation failed: %w", err)
	}

	return nil
}
