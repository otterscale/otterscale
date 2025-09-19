package core

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const VirtualMachineLabelKey = "otterscale.io/virtualmachine"

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
		svcspec.Ports[i].TargetPort = intstr.FromInt(int(svcspec.Ports[i].Port))
	}

	if svcspec.Type != corev1.ServiceTypeNodePort {
		for i := range svcspec.Ports {
			svcspec.Ports[i].NodePort = 0
		}
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
	var label = "otterscale.io/kind=vm-service"

    return uc.kubeCore.ListServicesByOptions(ctx, config, namespace, label, "")

}

func (uc *KubeVirtUseCase) UpdateVirtualMachineService(ctx context.Context, uuid, facility, namespace, name string, svcspec *corev1.ServiceSpec) (*corev1.Service, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	currentService, err := uc.kubeCore.GetService(ctx, config, namespace, name)
	if err != nil {
		return nil, err
	}

	newSpec := currentService.Spec

	newSpec.Ports = []corev1.ServicePort{}
	for _, p := range svcspec.Ports {
		serviceport := corev1.ServicePort{
			Name:       p.Name,
			Protocol:   p.Protocol,
			Port:       p.Port,
			TargetPort: intstr.FromInt(int(p.Port)),
		}
		if p.NodePort > 0 {
			serviceport.NodePort = p.NodePort
		}
		newSpec.Ports = append(newSpec.Ports, serviceport)
	}
	newSpec.Type = currentService.Spec.Type

	return uc.kubeCore.UpdateService(ctx, config, namespace, name, &newSpec)
}

func (uc *KubeVirtUseCase) DeleteVirtualMachineService(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeCore.DeleteService(ctx, config, namespace, name)
}
