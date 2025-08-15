package kube

import (
	"context"

	oscore "github.com/openhdc/otterscale/internal/core"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (r *virtDV) CreateDataVolume(ctx context.Context, config *rest.Config, namespace, name string, source_type string, source string, sizeBytes int64) (*oscore.DataVolume, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}
	var dvSource *v1beta1.DataVolumeSource

	switch {
	case source_type == "HTTP":
		dvSource = &v1beta1.DataVolumeSource{
			HTTP: &v1beta1.DataVolumeSourceHTTP{URL: source},
		}
	case source_type == "PVC":
		dvSource = &v1beta1.DataVolumeSource{
			PVC: &v1beta1.DataVolumeSourcePVC{Namespace: namespace, Name: source},
		}
	case source_type == "Blank":
		dvSource = &v1beta1.DataVolumeSource{Blank: &v1beta1.DataVolumeBlankImage{}}
	case source_type == "Registry":
		dvSource = &v1beta1.DataVolumeSource{
			Registry: &v1beta1.DataVolumeSourceRegistry{URL: &source},
		}
	case source_type == "Upload":
		dvSource = &v1beta1.DataVolumeSource{Upload: &v1beta1.DataVolumeSourceUpload{}}
	case source_type == "S3":
		dvSource = &v1beta1.DataVolumeSource{S3: &v1beta1.DataVolumeSourceS3{URL: source}}
	case source_type == "VDDK":
		dvSource = &v1beta1.DataVolumeSource{VDDK: &v1beta1.DataVolumeSourceVDDK{URL: source}}
	default:
		return nil, err
	}

	pvcSpec := &v1.PersistentVolumeClaimSpec{
		AccessModes: []v1.PersistentVolumeAccessMode{
			v1.ReadWriteMany,
		},
		Resources: v1.VolumeResourceRequirements{
			Requests: v1.ResourceList{
				v1.ResourceStorage: *resource.NewQuantity(sizeBytes, resource.BinarySI),
			},
		},
	}

	dv := &v1beta1.DataVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1beta1.DataVolumeSpec{
			Source: dvSource,
			PVC:    pvcSpec,
		},
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
	return virtClient.CdiClient().CdiV1beta1().DataVolumes(namespace).Get(ctx, name, opts)
}

func (r *virtDV) ListDataVolume(ctx context.Context, config *rest.Config, namespace string) ([]oscore.DataVolume, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	dvs, err := virtClient.CdiClient().CdiV1beta1().DataVolumes(namespace).List(ctx, opts)

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
