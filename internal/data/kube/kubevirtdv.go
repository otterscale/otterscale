package kube

import (
	"context"

	oscore "github.com/openhdc/otterscale/internal/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	v1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
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

func (r *virtDV) CreateDataVolume(ctx context.Context, config *rest.Config, namespace, name string, spec *oscore.DataVolumeSpec) (*oscore.DataVolume, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	dv := &v1beta1.DataVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if spec != nil {
		dv.Spec = *spec
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

func (r *virtDV) ListDataVolume(ctx context.Context, config *rest.Config, namespace, name string) ([]oscore.DataVolume, error) {
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
