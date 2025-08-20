package kube

import (
	"context"
	"math"

	oscore "github.com/openhdc/otterscale/internal/core"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	kubevirtv1 "kubevirt.io/api/instancetype/v1beta1"
)

func resourcePtr(bytes int64) *resource.Quantity {
	q := resource.NewQuantity(bytes, resource.BinarySI)
	return q
}

type virtInstanceType struct {
	kubevirt *kubevirt
}

func NewVirtInstanceType(kube *Kube, kubevirt *kubevirt) oscore.KubeVirtInstanceTypeRepo {
	return &virtInstanceType{
		kubevirt: kubevirt,
	}
}

var _ oscore.KubeVirtInstanceTypeRepo = (*virtInstanceType)(nil)

func (r *virtInstanceType) CreateInstanceType(ctx context.Context, config *rest.Config, InstanceType oscore.InstanceType) (*oscore.InstanceType, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	obj := &kubevirtv1.VirtualMachineClusterInstancetype{
		ObjectMeta: metav1.ObjectMeta{
			Name:        InstanceType.Metadata.Name,
			Labels:      InstanceType.Metadata.Labels,
			Annotations: InstanceType.Metadata.Annotations,
		},
		Spec: kubevirtv1.VirtualMachineInstancetypeSpec{
			CPU: kubevirtv1.CPUInstancetype{
				Guest: uint32(math.Round(float64(InstanceType.CPUCores))),
			},
			Memory: kubevirtv1.MemoryInstancetype{
				Guest: *resourcePtr(InstanceType.MemoryBytes),
			},
		},
	}

	opts := metav1.CreateOptions{}
	created, err := virtClient.VirtualMachineClusterInstancetype().Create(ctx, obj, opts)
	if err != nil {
		return nil, err
	}

	return &oscore.InstanceType{
		Metadata: oscore.Metadata{
			Name:        created.Name,
			Labels:      created.Labels,
			Annotations: created.Annotations,
		},
		CPUCores:    float32(created.Spec.CPU.Guest),
		MemoryBytes: created.Spec.Memory.Guest.Value(),
	}, nil
}

func (r *virtInstanceType) GetInstanceType(ctx context.Context, config *rest.Config, name string) (*oscore.InstanceType, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	obj, err := virtClient.VirtualMachineClusterInstancetype().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &oscore.InstanceType{
		Metadata: oscore.Metadata{
			Name:        obj.Name,
			Labels:      obj.Labels,
			Annotations: obj.Annotations,
		},
		CPUCores:    float32(obj.Spec.CPU.Guest),
		MemoryBytes: obj.Spec.Memory.Guest.Value(),
	}, nil
}

func (r *virtInstanceType) ListInstanceTypes(ctx context.Context, config *rest.Config) ([]oscore.InstanceType, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	list, err := virtClient.VirtualMachineClusterInstancetype().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []oscore.InstanceType
	for _, obj := range list.Items {
		result = append(result, oscore.InstanceType{
			Metadata: oscore.Metadata{
				Name:        obj.Name,
				Labels:      obj.Labels,
				Annotations: obj.Annotations,
			},
			CPUCores:    float32(obj.Spec.CPU.Guest),
			MemoryBytes: obj.Spec.Memory.Guest.Value(),
		})
	}
	return result, nil
}

func (r *virtInstanceType) DeleteInstanceType(ctx context.Context, config *rest.Config, name string) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}
	return virtClient.VirtualMachineClusterInstancetype().Delete(ctx, name, metav1.DeleteOptions{})
}
