package core

import (
    "context"

    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
    "google.golang.org/protobuf/types/known/timestamppb"
)

const vmLabelKey = "kubevirt.io/vm"

/*type KubeVirtVMServiceRepo interface {
    CreateVMService(ctx context.Context, config *rest.Config, namespace, name string, spec *KubeVirtVMService) (*KubeVirtVMService, error)
    GetVMService(ctx context.Context, config *rest.Config, namespace, name string) (*KubeVirtVMService, error)
    ListVMServices(ctx context.Context, config *rest.Config, namespace, name string) ([]KubeVirtVMService, error)
    DeleteVMService(ctx context.Context, config *rest.Config, namespace, name string) error
    UpdateVMService(ctx context.Context, config *rest.Config, namespace, name string, spec *KubeVirtVMService) (*KubeVirtVMService, error)
}*/

func (uc *KubeVirtUseCase) CreateVMService(ctx context.Context, uuid, facility, namespace, name string, vmservice KubeVirtVMService) (*KubeVirtVMService, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	selector := map[string]string{}
	for k, v := range vmservice.Selector {
		selector[k] = v
	}

	if vm, ok := selector[vmLabelKey]; !ok || vm == "" {
		selector[vmLabelKey] = name
	}
	
	ports := make([]corev1.ServicePort, 0, len(vmservice.Ports))
	svcType := corev1.ServiceTypeClusterIP
	for _, p := range vmservice.Ports {
		sp := corev1.ServicePort{
			Name:       p.Name,
			Protocol:   corev1.ProtocolTCP,
			Port:       p.Port,
			TargetPort: intstr.FromInt(int(p.TargetPort)),
		}
		if p.NodePort > 0 {
			svcType = corev1.ServiceTypeNodePort
			sp.NodePort = p.NodePort
		}
		ports = append(ports, sp)
	}

	spec := &corev1.ServiceSpec{
		Type:     svcType,
		Selector: selector,
		Ports:    ports,
	}

	created, err := uc.kubeCore.CreateService(ctx, config, namespace, name, spec)
	if err != nil {
		return nil, err
	}
	return fromK8sService(created), nil
}

func (uc *KubeVirtUseCase) GetVMService(ctx context.Context,uuid, facility, namespace, name string) (*KubeVirtVMService, error) {
    config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
    if err != nil {
        return nil, err
    }
    vmservice, err := uc.kubeCore.GetService(ctx, config, namespace, name)
    if err != nil {
        return nil, err
    }
    return fromK8sService(vmservice), nil
}

func (uc *KubeVirtUseCase) ListVMServices(ctx context.Context, uuid, facility, namespace string) ([]KubeVirtVMService, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	vmservices, err := uc.kubeCore.ListServices(ctx, config, namespace)
	if err != nil {
		return nil, err
	}
	ret := make([]KubeVirtVMService, 0, len(vmservices))
	for i := range vmservices {
		if vmservices[i].Spec.Selector != nil {
			if _, ok := vmservices[i].Spec.Selector[vmLabelKey]; ok {
				ret = append(ret, *fromK8sService(&vmservices[i]))
			}
		}
	}
	return ret, nil
}

func (uc *KubeVirtUseCase) UpdateVMService(ctx context.Context, uuid, facility, namespace, name string, vmservice KubeVirtVMService) (*KubeVirtVMService, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	current_vmservice, err := uc.kubeCore.GetService(ctx, config, namespace, name)
	if err != nil {
		return nil, err
	}

	newSpec := current_vmservice.Spec

	if newSpec.Selector == nil {
		newSpec.Selector = map[string]string{}
	}
	for k, v := range vmservice.Selector {
		newSpec.Selector[k] = v
	}
	if v := newSpec.Selector[vmLabelKey]; v == "" {
		newSpec.Selector[vmLabelKey] = name
	}

	newSpec.Ports = newSpec.Ports[:0]
	newType := corev1.ServiceTypeClusterIP
	for _, p := range vmservice.Ports {
		sp := corev1.ServicePort{
			Name:       p.Name,
			Protocol:   corev1.ProtocolTCP,
			Port:       p.Port,
			TargetPort: intstr.FromInt(int(p.TargetPort)),
		}
		if p.NodePort > 0 {
			newType = corev1.ServiceTypeNodePort
			sp.NodePort = p.NodePort
		}
		newSpec.Ports = append(newSpec.Ports, sp)
	}
	newSpec.Type = newType
	if newSpec.Type != corev1.ServiceTypeNodePort && len(newSpec.Ports) > 0 {
		newSpec.Ports[0].NodePort = 0
	}

	updated, err := uc.kubeCore.UpdateService(ctx, config, namespace, name, &newSpec)
	if err != nil {
		return nil, err
	}
	return fromK8sService(updated), nil
}

func (uc *KubeVirtUseCase) DeleteVMService(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeCore.DeleteService(ctx, config, namespace, name)
}

func fromK8sService(svc *corev1.Service) *KubeVirtVMService {
	ret := &KubeVirtVMService{
		Metadata: Metadata{
			Name:        svc.Name,
			Namespace:   svc.Namespace,
			Labels:      svc.Labels,
			Annotations: svc.Annotations,
			CreatedAt:   timestamppb.New(svc.CreationTimestamp.Time),
		},
		Type:     string(svc.Spec.Type),
		Selector: svc.Spec.Selector,
	}
	ret.Ports = make([]KubeVirtVMServicePort, 0, len(svc.Spec.Ports))
	for _, p := range svc.Spec.Ports {
		tp := int32(0)
		if p.TargetPort.Type == intstr.Int {
			tp = int32(p.TargetPort.IntValue())
		}
		ret.Ports = append(ret.Ports, KubeVirtVMServicePort{
			Name:       p.Name,
			Port:       p.Port,
			NodePort:   p.NodePort,
			TargetPort: tp,
		})
	}
	return ret
}
