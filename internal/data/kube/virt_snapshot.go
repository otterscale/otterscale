package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	virtcore "kubevirt.io/api/core"
	snapshotv1beta1 "kubevirt.io/api/snapshot/v1beta1"

	oscore "github.com/otterscale/otterscale/internal/core"
)

type virtSnapshot struct {
	kube *Kube
}

func NewVirtSnapshot(kube *Kube) oscore.KubeVirtSnapshotRepo {
	return &virtSnapshot{
		kube: kube,
	}
}

var _ oscore.KubeVirtSnapshotRepo = (*virtSnapshot)(nil)

func (r *virtSnapshot) CreateVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name, vmName string) (*oscore.VirtualMachineSnapshot, error) {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return nil, err
	}
	apiGroup := virtcore.GroupName
	kind := "VirtualMachine"
	snapshot := &snapshotv1beta1.VirtualMachineSnapshot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: snapshotv1beta1.VirtualMachineSnapshotSpec{
			Source: corev1.TypedLocalObjectReference{
				APIGroup: &apiGroup,
				Kind:     kind,
				Name:     vmName,
			},
		},
	}
	opts := metav1.CreateOptions{}
	return clientset.SnapshotV1beta1().VirtualMachineSnapshots(namespace).Create(ctx, snapshot, opts)
}

func (r *virtSnapshot) DeleteVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string) error {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return err
	}
	opts := metav1.DeleteOptions{}
	return clientset.SnapshotV1beta1().VirtualMachineSnapshots(namespace).Delete(ctx, name, opts)
}

func (r *virtSnapshot) CreateVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name, target, snapshot string) (*oscore.VirtualMachineRestore, error) {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return nil, err
	}
	apiGroup := virtcore.GroupName
	kind := "VirtualMachine"
	restore := &snapshotv1beta1.VirtualMachineRestore{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: snapshotv1beta1.VirtualMachineRestoreSpec{
			Target: corev1.TypedLocalObjectReference{
				APIGroup: &apiGroup,
				Kind:     kind,
				Name:     target,
			},
			VirtualMachineSnapshotName: snapshot,
		},
	}
	opts := metav1.CreateOptions{}
	return clientset.SnapshotV1beta1().VirtualMachineRestores(namespace).Create(ctx, restore, opts)
}

func (r *virtSnapshot) DeleteVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string) error {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return err
	}
	opts := metav1.DeleteOptions{}
	return clientset.SnapshotV1beta1().VirtualMachineRestores(namespace).Delete(ctx, name, opts)
}
