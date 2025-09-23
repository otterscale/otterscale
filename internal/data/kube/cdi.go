package kube

import (
	"context"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"

	oscore "github.com/otterscale/otterscale/internal/core"
)

type cdi struct {
	kube *Kube
}

func NewCDI(kube *Kube) oscore.KubeCDIRepo {
	return &cdi{
		kube: kube,
	}
}

var _ oscore.KubeCDIRepo = (*cdi)(nil)

func (r *cdi) ListDataVolumes(ctx context.Context, config *rest.Config, namespace string) ([]oscore.DataVolume, error) {
	clientset, err := r.kube.cdiClientset(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	list, err := clientset.CdiV1beta1().DataVolumes(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *cdi) GetDataVolume(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.DataVolume, error) {
	clientset, err := r.kube.cdiClientset(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	return clientset.CdiV1beta1().DataVolumes(namespace).Get(ctx, name, opts)
}

func (r *cdi) CreateDataVolume(ctx context.Context, config *rest.Config, namespace, name string, srcType oscore.SourceType, srcData string, size int64, bootImage bool) (*oscore.DataVolume, error) {
	clientset, err := r.kube.cdiClientset(config)
	if err != nil {
		return nil, err
	}
	spec := r.buildDataVolumeSpec(srcType, srcData, namespace, size)
	dataVolume := &cdiv1beta1.DataVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				oscore.DataVolumeBootImageLabel: strconv.FormatBool(bootImage),
			},
		},
		Spec: spec,
	}
	opts := metav1.CreateOptions{}
	return clientset.CdiV1beta1().DataVolumes(namespace).Create(ctx, dataVolume, opts)
}

func (r *cdi) DeleteDataVolume(ctx context.Context, config *rest.Config, namespace, name string) error {
	clientset, err := r.kube.cdiClientset(config)
	if err != nil {
		return err
	}
	opts := metav1.DeleteOptions{}
	return clientset.CdiV1beta1().DataVolumes(namespace).Delete(ctx, name, opts)
}

func (r *cdi) buildDataVolumeSpec(srcType oscore.SourceType, srcData, namespace string, size int64) cdiv1beta1.DataVolumeSpec {
	if srcType == oscore.SourceTypePVC {
		return r.buildPVCSpec(srcData, namespace, size)
	}
	return r.buildStorageSpec(srcType, srcData, size)
}

func (r *cdi) buildPVCSpec(srcData, namespace string, size int64) cdiv1beta1.DataVolumeSpec {
	source := &cdiv1beta1.DataVolumeSource{
		PVC: &cdiv1beta1.DataVolumeSourcePVC{
			Namespace: namespace,
			Name:      srcData,
		},
	}
	pvc := &corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		Resources: corev1.VolumeResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceStorage: *resource.NewQuantity(size, resource.BinarySI),
			},
		},
	}
	return cdiv1beta1.DataVolumeSpec{
		Source: source,
		PVC:    pvc,
	}
}

func (r *cdi) buildStorageSpec(srcType oscore.SourceType, srcData string, size int64) cdiv1beta1.DataVolumeSpec {
	source := &cdiv1beta1.DataVolumeSource{}
	switch srcType {
	case oscore.SourceTypeBlank:
		source.Blank = &cdiv1beta1.DataVolumeBlankImage{}
	case oscore.SourceTypeHTTP:
		source.HTTP = &cdiv1beta1.DataVolumeSourceHTTP{URL: srcData}
	}
	storage := &cdiv1beta1.StorageSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		Resources: corev1.VolumeResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceStorage: *resource.NewQuantity(size, resource.BinarySI),
			},
		},
	}
	return cdiv1beta1.DataVolumeSpec{
		Source:  source,
		Storage: storage,
	}
}
