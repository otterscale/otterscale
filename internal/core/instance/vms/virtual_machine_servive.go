package vms

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/service"
	"github.com/otterscale/otterscale/internal/core/instance/vm"
)

type VirtualMachineServiceUseCase struct {
	service service.ServiceRepo
}

func NewVirtualMachineServiceUseCase(service service.ServiceRepo) *VirtualMachineServiceUseCase {
	return &VirtualMachineServiceUseCase{
		service: service,
	}
}

func (uc *VirtualMachineServiceUseCase) CreateVirtualMachineService(ctx context.Context, scope, namespace, name, vmName string, ports []service.Port) (*service.Service, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				vm.VirtualMachineNameLabel: vmName,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: ports,
			Type:  corev1.ServiceTypeNodePort,
		},
	}

	return uc.service.Create(ctx, scope, namespace, service)
}

func (uc *VirtualMachineServiceUseCase) UpdateVirtualMachineService(ctx context.Context, scope, namespace, name string, ports []service.Port) (*service.Service, error) {
	service, err := uc.service.Get(ctx, scope, namespace, name)
	if err != nil {
		return nil, err
	}

	service.Spec.Ports = ports

	return uc.service.Update(ctx, scope, namespace, service)
}

func (uc *VirtualMachineServiceUseCase) DeleteVirtualMachineService(ctx context.Context, scope, namespace, name string) error {
	return uc.service.Delete(ctx, scope, namespace, name)
}
