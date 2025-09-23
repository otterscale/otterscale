package kube

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	virtv1 "kubevirt.io/api/core/v1"

	"github.com/google/uuid"
	oscore "github.com/otterscale/otterscale/internal/core"
)

type virt struct {
	kube *Kube
}

func NewVirt(kube *Kube) oscore.KubeVirtRepo {
	return &virt{
		kube: kube,
	}
}

var _ oscore.KubeVirtRepo = (*virt)(nil)

func (r *virt) StartVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return err
	}
	opts := &virtv1.StartOptions{}
	return clientset.KubevirtV1().VirtualMachines(namespace).Start(ctx, name, opts)
}

func (r *virt) StopVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return err
	}
	opts := &virtv1.StopOptions{}
	return clientset.KubevirtV1().VirtualMachines(namespace).Stop(ctx, name, opts)
}

func (r *virt) RestartVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return err
	}
	opts := &virtv1.RestartOptions{}
	return clientset.KubevirtV1().VirtualMachines(namespace).Restart(ctx, name, opts)
}

func (r *virt) PauseInstance(ctx context.Context, config *rest.Config, namespace, name string) error {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return err
	}
	opts := &virtv1.PauseOptions{}
	return clientset.KubevirtV1().VirtualMachineInstances(namespace).Pause(ctx, name, opts)
}

func (r *virt) UnpauseInstance(ctx context.Context, config *rest.Config, namespace, name string) error {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return err
	}
	opts := &virtv1.UnpauseOptions{}
	return clientset.KubevirtV1().VirtualMachineInstances(namespace).Unpause(ctx, name, opts)
}

func (r *virt) MigrateInstance(ctx context.Context, config *rest.Config, namespace, name, hostname string) (*oscore.VirtualMachineInstanceMigration, error) {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return nil, err
	}
	migration := &virtv1.VirtualMachineInstanceMigration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-migration-%s", name, uuid.New().String()),
			Namespace: namespace,
		},
		Spec: virtv1.VirtualMachineInstanceMigrationSpec{
			VMIName: name,
			AddedNodeSelector: map[string]string{
				"kubernetes.io/hostname": hostname,
			},
		},
	}
	opts := metav1.CreateOptions{}
	return clientset.KubevirtV1().VirtualMachineInstanceMigrations(namespace).Create(ctx, migration, opts)
}
