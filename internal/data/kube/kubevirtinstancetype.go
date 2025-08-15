package kube

import (
	"context"

	oscore "github.com/openhdc/otterscale/internal/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type virtInstancetype struct {
	kubevirt *kubevirt
}

func NewVirtInstancetype(kube *Kube, kubevirt *kubevirt) oscore.KubeVirtInstancetypeRepo {
	return &virtInstancetype{
		kubevirt: kubevirt,
	}
}

var _ oscore.KubeVirtInstancetypeRepo = (*virtInstancetype)(nil)

func (r *virtInstancetype) CreateInstancetype(ctx context.Context, config *rest.Config, name string) (*oscore.Instancetype, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}
	return virtClient.VirtualMachineClusterInstancetype().Create(ctx, obh, opts)
}

func (r *virtInstancetype) GetInstancetype(ctx context.Context, config *rest.Config, name string) (*oscore.Instancetype, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	return virtClient.VirtualMachineClusterInstancetype().Get(ctx, name, opts)
}

func (r *virtInstancetype) ListInstancetypes(ctx context.Context, config *rest.Config) ([]oscore.Instancetype, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	dvs, err := virtClient.CdiClient().CdiV1beta1().DataVolumes(namespace).List(ctx, opts)

	return dvs.Items, nil
}

func (r *virtInstancetype) DeleteInstancetype(ctx context.Context, config *rest.Config, name string) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}
	opts := metav1.DeleteOptions{}
	return virtClient.VirtualMachineClusterInstancetype().Delete(ctx, name, opts)
}
