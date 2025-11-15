package kubevirt

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/instance/vmi"
)

type virtualMachineInstanceMigrationRepo struct {
	kubevirt *KubeVirt
}

func NewVirtualMachineInstanceMigrationRepo(kubevirt *KubeVirt) vmi.VirtualMachineInstanceMigrationRepo {
	return &virtualMachineInstanceMigrationRepo{
		kubevirt: kubevirt,
	}
}

var _ vmi.VirtualMachineInstanceMigrationRepo = (*virtualMachineInstanceMigrationRepo)(nil)

func (r *virtualMachineInstanceMigrationRepo) Create(ctx context.Context, scope, namespace string, vmim *vmi.VirtualMachineInstanceMigration) (*vmi.VirtualMachineInstanceMigration, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.KubevirtV1().VirtualMachineInstanceMigrations(namespace).Create(ctx, vmim, opts)
}
