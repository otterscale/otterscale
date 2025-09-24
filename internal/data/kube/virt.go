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

const VirtualMachineKind = "VirtualMachine"

type virt struct {
	kube *Kube
}

func NewVirt(kube *Kube) oscore.KubeVirtRepo {
	return &virt{
		kube: kube,
	}
}

var _ oscore.KubeVirtRepo = (*virt)(nil)

func (r *virt) ListVirtualMachines(ctx context.Context, config *rest.Config, namespace string) ([]oscore.VirtualMachine, error) {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	list, err := clientset.KubevirtV1().VirtualMachines(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *virt) GetVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.VirtualMachine, error) {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	return clientset.KubevirtV1().VirtualMachines(namespace).Get(ctx, name, opts)
}

func (r *virt) CreateVirtualMachine(ctx context.Context, config *rest.Config, namespace, name, instanceType, bootDataVolume, startupScript string) (*oscore.VirtualMachine, error) {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return nil, err
	}
	var (
		runStrategy   = virtv1.RunStrategyHalted
		enabled       = true
		bootOrder     = uint(1)
		osDisk        = "os-disk"
		cloudInitDisk = "cloud-init-disk"
		nic1          = "nic1"
	)
	virtualMachine := &virtv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Annotations: map[string]string{
				"kubevirt.io/allow-pod-bridge-network-live-migration": "true",
			},
		},
		Spec: virtv1.VirtualMachineSpec{
			RunStrategy: &runStrategy,
			Instancetype: &virtv1.InstancetypeMatcher{
				Name: instanceType,
			},
			Template: &virtv1.VirtualMachineInstanceTemplateSpec{
				Spec: virtv1.VirtualMachineInstanceSpec{
					Domain: virtv1.DomainSpec{
						Devices: virtv1.Devices{
							Disks: []virtv1.Disk{
								{
									Name: osDisk,
									DiskDevice: virtv1.DiskDevice{
										Disk: &virtv1.DiskTarget{
											Bus: virtv1.DiskBusVirtio,
										},
									},
									BootOrder: &bootOrder,
								},
							},
							Interfaces: []virtv1.Interface{
								{
									Name: nic1,
									InterfaceBindingMethod: virtv1.InterfaceBindingMethod{
										Bridge: &virtv1.InterfaceBridge{},
									},
								},
							},
							TPM: &virtv1.TPMDevice{
								Enabled: &enabled,
							},
						},
					},
					Networks: []virtv1.Network{
						{
							Name: nic1,
							NetworkSource: virtv1.NetworkSource{
								Pod: &virtv1.PodNetwork{},
							},
						},
					},
					Volumes: []virtv1.Volume{
						{
							Name: osDisk,
							VolumeSource: virtv1.VolumeSource{
								DataVolume: &virtv1.DataVolumeSource{
									Name: bootDataVolume,
								},
							},
						},
					},
				},
			},
		},
	}
	if startupScript != "" {
		virtualMachine.Spec.Template.Spec.Domain.Devices.Disks = append(virtualMachine.Spec.Template.Spec.Domain.Devices.Disks, virtv1.Disk{
			Name: cloudInitDisk,
			DiskDevice: virtv1.DiskDevice{
				Disk: &virtv1.DiskTarget{
					Bus: virtv1.DiskBusVirtio,
				},
			},
		})
		virtualMachine.Spec.Template.Spec.Volumes = append(virtualMachine.Spec.Template.Spec.Volumes, virtv1.Volume{
			Name: cloudInitDisk,
			VolumeSource: virtv1.VolumeSource{
				CloudInitNoCloud: &virtv1.CloudInitNoCloudSource{
					UserData: startupScript,
				},
			},
		})
	}
	opts := metav1.CreateOptions{}
	return clientset.KubevirtV1().VirtualMachines(namespace).Create(ctx, virtualMachine, opts)
}

func (r *virt) UpdateVirtualMachine(ctx context.Context, config *rest.Config, namespace string, virtualMachine *oscore.VirtualMachine) (*oscore.VirtualMachine, error) {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.UpdateOptions{}
	return clientset.KubevirtV1().VirtualMachines(namespace).Update(ctx, virtualMachine, opts)
}

func (r *virt) DeleteVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return err
	}
	opts := metav1.DeleteOptions{}
	return clientset.KubevirtV1().VirtualMachines(namespace).Delete(ctx, name, opts)
}

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

func (r *virt) ListInstances(ctx context.Context, config *rest.Config, namespace string) ([]oscore.VirtualMachineInstance, error) {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	list, err := clientset.KubevirtV1().VirtualMachineInstances(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *virt) GetInstance(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.VirtualMachineInstance, error) {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	return clientset.KubevirtV1().VirtualMachineInstances(namespace).Get(ctx, name, opts)
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
