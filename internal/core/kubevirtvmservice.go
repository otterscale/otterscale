package core

import (
	"context"
	"fmt"
	"strings"
	
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

const VirtualMachineLabelKey = "otterscale.io/virtualmachine"
const VirtualMachineServiceLabelKey = "otterscale.io/virtualmachineservice"

func (uc *KubeVirtUseCase) CreateVirtualMachineService(ctx context.Context, uuid, facility, namespace, name string, svcspec *corev1.ServiceSpec) (*corev1.Service, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	if svcspec.Selector == nil {
		svcspec.Selector = map[string]string{}
	}
	if svcspec.Selector[VirtualMachineLabelKey] == "" {
		svcspec.Selector[VirtualMachineLabelKey] = name
	}

	for i := range svcspec.Ports {
		svcspec.Ports[i].Name = fmt.Sprintf("%s-%d", strings.ToLower(string(svcspec.Ports[i].Protocol)), svcspec.Ports[i].Port)
		svcspec.Ports[i].TargetPort = intstr.FromInt(int(svcspec.Ports[i].Port))
	}

	svc, err := uc.kubeCore.CreateVirtualMachineService(ctx, config, namespace, name, svcspec)
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func (uc *KubeVirtUseCase) GetVirtualMachineService(ctx context.Context, uuid, facility, namespace, name string) (*corev1.Service, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeCore.GetService(ctx, config, namespace, name)
}

func (uc *KubeVirtUseCase) ListVirtualMachineServices(ctx context.Context, uuid, facility, namespace string) ([]corev1.Service, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

    return uc.kubeCore.ListServicesByOptions(ctx, config, namespace, VirtualMachineServiceLabelKey, "")

}

func (uc *KubeVirtUseCase) UpdateVirtualMachineService(ctx context.Context, uuid, facility, namespace, name string, svcspec *corev1.ServiceSpec) (*corev1.Service, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	currentService, err := uc.GetVirtualMachineService(ctx, uuid, facility, namespace, name)
	if err != nil {
		return nil, err
	}

	newSpec := currentService.Spec

	newSpec.Ports = []corev1.ServicePort{}
	for _, p := range svcspec.Ports {
		serviceport := corev1.ServicePort{
			Name:       fmt.Sprintf("%s-%d", strings.ToLower(string(p.Protocol)), p.Port),
			Protocol:   p.Protocol,
			Port:       p.Port,
			TargetPort: intstr.FromInt(int(p.Port)),
		}
		newSpec.Ports = append(newSpec.Ports, serviceport)
	}
	newSpec.Type = svcspec.Type

	return uc.kubeCore.UpdateVirtualMachineService(ctx, config, namespace, name, &newSpec)
}

func (uc *KubeVirtUseCase) DeleteVirtualMachineService(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeCore.DeleteService(ctx, config, namespace, name)
}

func (uc *KubeVirtUseCase) ExposeVirtualMachine(ctx context.Context, uuid, facility, namespace, name string, svcspec *corev1.ServiceSpec) (*corev1.Service, error) {
	_, err := uc.GetVirtualMachineService(ctx, uuid, facility, namespace, name)

	if apierrors.IsNotFound(err) {
		return uc.CreateVirtualMachineService(ctx, uuid, facility, namespace, name, svcspec)
	}
	if err != nil {
		return nil, err
	}
	
	return uc.UpdateVirtualMachineService(ctx, uuid, facility, namespace, name, svcspec)
}
