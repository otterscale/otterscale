package kube

import (
	"context"
	"fmt"
	"strconv"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"

	oscore "github.com/otterscale/otterscale/internal/core"
)

const (
	SourceHTTP     = "http"
	SourceBlank    = "blank"
	SourceRegistry = "registry"
	SourceUpload   = "upload"
	SourceS3       = "s3"
	SourceVDDK     = "vddk"
	SourcePVC      = "pvc"
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

// CreateDataVolume creates a DataVolume. Depending on the source_type, constructs different DataVolumeSources.
// Sets PVC or Storage Spec as needed, then calls the CDI API to create it.
func (r *virtDV) CreateDataVolume(ctx context.Context, config *rest.Config, namespace, name, sourceType, source, vmName string, sizeBytes int64, isBootable bool) (*oscore.DataVolume, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	dvSpec := &v1beta1.DataVolumeSpec{}
	var dvSource *v1beta1.DataVolumeSource

	// Decide the DataVolumeSource according to source_type
	switch sourceType {
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
		// Use an existing PVC directly, and create a PVC spec for DataVolume reference
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
		// Unsupported source_type returns an error directly
		return nil, fmt.Errorf("unsupported source_type: %s", sourceType)
	}

	// If not PVC, use StorageSpec to fill in the size information
	if sourceType != "PVC" {
		dvSpec.Storage = oscore.StorageSpec(sizeBytes)
	}
	dvSpec.Source = dvSource

	// Create the DataVolume object
	dv := &v1beta1.DataVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"otterscale.io/is_bootable":    strconv.FormatBool(isBootable),
				"otterscale.io/virtualmachine": vmName,
			},
		},
		Spec: *dvSpec,
	}

	opts := metav1.CreateOptions{}
	return virtClient.CdiClient().CdiV1beta1().DataVolumes(namespace).Create(ctx, dv, opts)
}

// GetDataVolume gets the specified DataVolume and synchronizes its PVC resource and storageClass info.
func (r *virtDV) GetDataVolume(ctx context.Context, config *rest.Config, namespace, name string, pvc *v1.PersistentVolumeClaim) (*oscore.DataVolume, error) {
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

// ListDataVolume lists all DataVolumes in the specified namespace.
func (r *virtDV) ListDataVolumes(ctx context.Context, config *rest.Config, namespace string) ([]oscore.DataVolume, error) {
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

func (r *virtDV) ListDataVolumesByLabel(ctx context.Context, config *rest.Config, namespace, label string) ([]oscore.DataVolume, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{
		LabelSelector: label,
	}
	dvs, err := virtClient.CdiClient().CdiV1beta1().DataVolumes(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return dvs.Items, nil
}

// DeleteDataVolume deletes the specified DataVolume.
func (r *virtDV) DeleteDataVolume(ctx context.Context, config *rest.Config, namespace, name string) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}
	opts := metav1.DeleteOptions{}
	return virtClient.CdiClient().CdiV1beta1().DataVolumes(namespace).Delete(ctx, name, opts)
}

// ExtendDataVolume expands the existing PVC and the corresponding DataVolume.
func (r *virtDV) ExtendDataVolume(ctx context.Context, config *rest.Config, namespace string, pvc *v1.PersistentVolumeClaim, sizeBytes resource.Quantity) error {
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
