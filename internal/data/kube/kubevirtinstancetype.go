package kube

import (
	"context"
	"math"

	oscore "github.com/openhdc/otterscale/internal/core"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	kubevirtv1 "kubevirt.io/api/instancetype/v1alpha1"
)

func resourcePtr(bytes int64) *resource.Quantity {
	q := resource.NewQuantity(bytes, resource.BinarySI)
	return q
}

type virtInstancetype struct {
	kubevirt *kubevirt
}

func NewVirtInstancetype(kube *Kube, kubevirt *kubevirt) oscore.KubeVirtInstancetypeRepo {
	return &virtInstancetype{
		kubevirt: kubevirt,
	}
}

var _ oscore.KubeVirtInstancetypeRepo = (*virtInstancetype)(nil)

func (r *virtInstancetype) CreateInstancetype(ctx context.Context, config *rest.Config, instancetype oscore.Instancetype) (*oscore.Instancetype, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	obj := &kubevirtv1.VirtualMachineClusterInstancetype{
		ObjectMeta: metav1.ObjectMeta{
			Name:        instancetype.Metadata.Name,
			Labels:      instancetype.Metadata.Labels,
			Annotations: instancetype.Metadata.Annotations,
		},
		Spec: kubevirtv1.VirtualMachineInstancetypeSpec{
			CPU: kubevirtv1.CPUInstancetype{
				Guest: uint32(math.Round(float64(instancetype.CpuCores))),
			},
			Memory: kubevirtv1.MemoryInstancetype{
				Guest: *resourcePtr(instancetype.MemoryBytes),
			},
		},
	}

	created, err := virtClient.VirtualMachineClusterInstancetype().Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return &oscore.Instancetype{
		Metadata: oscore.Metadata{
			Name:        created.Name,
			Labels:      created.Labels,
			Annotations: created.Annotations,
			// CreatedAt/UpdatedAt 可視需要補充
		},
		CpuCores:    float32(created.Spec.CPU.Guest),
		MemoryBytes: created.Spec.Memory.Guest.Value(),
	}, nil
}

func (r *virtInstancetype) GetInstancetype(ctx context.Context, config *rest.Config, name string) (*oscore.Instancetype, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	obj, err := virtClient.VirtualMachineClusterInstancetype().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &oscore.Instancetype{
		Metadata: oscore.Metadata{
			Name:        obj.Name,
			Labels:      obj.Labels,
			Annotations: obj.Annotations,
		},
		CpuCores:    float32(obj.Spec.CPU.Guest),
		MemoryBytes: obj.Spec.Memory.Guest.Value(),
	}, nil
}

func (r *virtInstancetype) ListInstancetypes(ctx context.Context, config *rest.Config) ([]oscore.Instancetype, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	virtClient.VirtualMachineInstancetype().List()
	list, err := virtClient.VirtualMachineInstancetype().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []oscore.Instancetype
	for _, obj := range list.Items {
		result = append(result, oscore.Instancetype{
			Metadata: oscore.Metadata{
				Name:        obj.Name,
				Labels:      obj.Labels,
				Annotations: obj.Annotations,
			},
			CpuCores:    float32(obj.Spec.CPU.Guest),
			MemoryBytes: obj.Spec.Memory.Guest.Value(),
		})
	}
	return result, nil
}

func (r *virtInstancetype) DeleteInstancetype(ctx context.Context, config *rest.Config, name string) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}
	return virtClient.VirtualMachineClusterInstancetype().Delete(ctx, name, metav1.DeleteOptions{})
}
