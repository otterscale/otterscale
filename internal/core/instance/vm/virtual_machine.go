package vm

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kvcorev1 "kubevirt.io/api/core/v1"

	"github.com/otterscale/otterscale/internal/core/application/service"
	"github.com/otterscale/otterscale/internal/core/instance/vmi"
	"github.com/otterscale/otterscale/internal/core/machine"
)

const (
	VirtualMachineDiskBusVirtio = kvcorev1.DiskBusVirtio
	VirtualMachineDiskBusSATA   = kvcorev1.DiskBusSATA
	VirtualMachineDiskBusSCSI   = kvcorev1.DiskBusSCSI
	VirtualMachineDiskBusUSB    = kvcorev1.DiskBusUSB
)

const nameLabel = "otterscale.com/virtual-machine.name"

const (
	groupName = "kubevirt.io"
	kind      = "VirtualMachine"
)

type (
	// VirtualMachine represents a KubeVirt VirtualMachine resource.
	VirtualMachine = kvcorev1.VirtualMachine

	// VirtualMachineDisk represents a KubeVirt Disk resource.
	VirtualMachineDisk = kvcorev1.Disk

	// VirtualMachineVolume represents a KubeVirt Volume resource.
	VirtualMachineVolume = kvcorev1.Volume

	// VirtualMachineVolumeSource represents a KubeVirt VolumeSource resource.
	VirtualMachineVolumeSource = kvcorev1.VolumeSource

	// VirtualMachineDiskBus represents a KubeVirt DiskBus resource.
	VirtualMachineDiskBus = kvcorev1.DiskBus
)

type VirtualMachineData struct {
	*VirtualMachine
	Instance  *vmi.VirtualMachineInstance
	Machine   *machine.Machine
	Clones    []VirtualMachineClone
	Snapshots []VirtualMachineSnapshot
	Restores  []VirtualMachineRestore
	Services  []service.Service
}

type VirtualMachineRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]VirtualMachine, error)
	Get(ctx context.Context, scope, namespace, name string) (*VirtualMachine, error)
	Create(ctx context.Context, scope, namespace string, vm *VirtualMachine) (*VirtualMachine, error)
	Update(ctx context.Context, scope, namespace string, vm *VirtualMachine) (*VirtualMachine, error)
	Delete(ctx context.Context, scope, namespace, name string) error
	Start(ctx context.Context, scope, namespace, name string) error
	Stop(ctx context.Context, scope, namespace, name string) error
	Restart(ctx context.Context, scope, namespace, name string) error
}

type UseCase struct {
	virtualMachine         VirtualMachineRepo
	virtualMachineClone    VirtualMachineCloneRepo
	virtualMachineRestore  VirtualMachineRestoreRepo
	virtualMachineSnapshot VirtualMachineSnapshotRepo

	machine                machine.MachineRepo
	service                service.ServiceRepo
	virtualMachineInstance vmi.VirtualMachineInstanceRepo
}

func NewUseCase(virtualMachine VirtualMachineRepo, virtualMachineClone VirtualMachineCloneRepo, virtualMachineRestore VirtualMachineRestoreRepo, virtualMachineSnapshot VirtualMachineSnapshotRepo, machine machine.MachineRepo, service service.ServiceRepo, virtualMachineInstance vmi.VirtualMachineInstanceRepo) *UseCase {
	return &UseCase{
		virtualMachine:         virtualMachine,
		virtualMachineClone:    virtualMachineClone,
		virtualMachineRestore:  virtualMachineRestore,
		virtualMachineSnapshot: virtualMachineSnapshot,
		service:                service,
		machine:                machine,
		virtualMachineInstance: virtualMachineInstance,
	}
}

func (uc *UseCase) ListVirtualMachines(ctx context.Context, scope, namespace string) ([]VirtualMachineData, error) {
	var (
		virtualMachines []VirtualMachine
		instances       []vmi.VirtualMachineInstance
		machines        []machine.Machine
		clones          []VirtualMachineClone
		snapshots       []VirtualMachineSnapshot
		restores        []VirtualMachineRestore
		services        []service.Service
	)

	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		v, err := uc.virtualMachine.List(egctx, scope, namespace, "")
		if err == nil {
			virtualMachines = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.virtualMachineInstance.List(egctx, scope, namespace, "")
		if err == nil {
			instances = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.machine.List(egctx)
		if err == nil {
			machines = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.virtualMachineClone.List(egctx, scope, namespace, "")
		if err == nil {
			clones = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.virtualMachineSnapshot.List(egctx, scope, namespace, "")
		if err == nil {
			snapshots = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.virtualMachineRestore.List(egctx, scope, namespace, "")
		if err == nil {
			restores = v
		}
		return err
	})

	eg.Go(func() error {
		selector := nameLabel

		v, err := uc.service.List(egctx, scope, namespace, selector)
		if err == nil {
			services = v
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return uc.combineVirtualMachines(virtualMachines, instances, machines, clones, snapshots, restores, services), nil
}

func (uc *UseCase) GetVirtualMachine(ctx context.Context, scope, namespace, name string) (*VirtualMachineData, error) {
	var (
		virtualMachine *VirtualMachine
		instance       *vmi.VirtualMachineInstance
	)

	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		v, err := uc.virtualMachine.Get(egctx, scope, namespace, name)
		if err == nil {
			virtualMachine = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.virtualMachineInstance.Get(egctx, scope, namespace, name)
		if err == nil {
			instance = v
		}
		if k8serrors.IsNotFound(err) {
			return nil
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	var (
		machines  []machine.Machine
		clones    []VirtualMachineClone
		snapshots []VirtualMachineSnapshot
		restores  []VirtualMachineRestore
		services  []service.Service
	)

	eg, egctx = errgroup.WithContext(ctx)

	eg.Go(func() error {
		v, err := uc.machine.List(egctx)
		if err == nil {
			machines = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.virtualMachineClone.List(egctx, scope, namespace, name)
		if err == nil {
			clones = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.virtualMachineSnapshot.List(egctx, scope, namespace, name)
		if err == nil {
			snapshots = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.virtualMachineRestore.List(egctx, scope, namespace, name)
		if err == nil {
			restores = v
		}
		return err
	})

	eg.Go(func() error {
		selector := nameLabel + "=" + name

		v, err := uc.service.List(egctx, scope, namespace, selector)
		if err == nil {
			services = v
		}
		return err
	})

	return uc.combineVirtualMachine(namespace, name, virtualMachine, instance, machines, clones, snapshots, restores, services), nil
}

func (uc *UseCase) CreateVirtualMachine(ctx context.Context, scope, namespace, name, instanceType, bootDataVolume, startupScript string) (*VirtualMachineData, error) {
	virtualMachine, err := uc.virtualMachine.Create(ctx, scope, namespace, uc.buildVirtualMachine(namespace, name, instanceType, bootDataVolume, startupScript))
	if err != nil {
		return nil, err
	}

	return &VirtualMachineData{
		VirtualMachine: virtualMachine,
	}, nil
}

func (uc *UseCase) DeleteVirtualMachine(ctx context.Context, scope, namespace, name string) error {
	// Get related services before deleting the virtual machine
	selector := nameLabel + "=" + name

	services, err := uc.service.List(ctx, scope, namespace, selector)
	if err != nil {
		return err
	}

	// Get related snapshots before deleting the virtual machine
	snapshots, err := uc.virtualMachineSnapshot.List(ctx, scope, namespace, name)
	if err != nil {
		return err
	}

	// Delete the virtual machine first
	if err := uc.virtualMachine.Delete(ctx, scope, namespace, name); err != nil {
		return err
	}

	// Delete related services
	for i := range services {
		_ = uc.service.Delete(ctx, scope, namespace, services[i].Name)
	}

	// Delete related snapshots
	for i := range snapshots {
		_ = uc.virtualMachineSnapshot.Delete(ctx, scope, namespace, snapshots[i].Name)
	}

	return nil
}

func (uc *UseCase) AttachVirtualMachineDisk(ctx context.Context, scope, namespace, name, dvName string) (disk *VirtualMachineDisk, volume *VirtualMachineVolume, err error) {
	vm, err := uc.virtualMachine.Get(ctx, scope, namespace, name)
	if err != nil {
		return nil, nil, err
	}

	// check if disk already attached
	volumes := vm.Spec.Template.Spec.Volumes

	for i := range volumes {
		if volumes[i].Name == dvName {
			return nil, nil, fmt.Errorf("disk %s is already attached to virtual machine %s/%s", dvName, namespace, name)
		}
	}

	// attach disk
	vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks, kvcorev1.Disk{
		Name: dvName,
		DiskDevice: kvcorev1.DiskDevice{
			Disk: &kvcorev1.DiskTarget{
				Bus: kvcorev1.DiskBusVirtio,
			},
		},
	})

	// add volume
	vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, kvcorev1.Volume{
		Name: dvName,
		VolumeSource: kvcorev1.VolumeSource{
			DataVolume: &kvcorev1.DataVolumeSource{
				Name: dvName,
			},
		},
	})

	// update virtual machine
	newVM, err := uc.virtualMachine.Update(ctx, scope, namespace, vm)
	if err != nil {
		return nil, nil, err
	}

	// get disk and volume to verify
	disks := newVM.Spec.Template.Spec.Domain.Devices.Disks

	for i := range disks {
		if disks[i].Name == dvName {
			disk = &disks[i]
			break
		}
	}

	volumes = newVM.Spec.Template.Spec.Volumes

	for i := range volumes {
		if volumes[i].Name == dvName {
			volume = &volumes[i]
			break
		}
	}

	return disk, volume, nil
}

func (uc *UseCase) DetachVirtualMachineDisk(ctx context.Context, scope, namespace, name, dvName string) error {
	vm, err := uc.virtualMachine.Get(ctx, scope, namespace, name)
	if err != nil {
		return err
	}

	// check if disk is attached
	found := false
	volumes := vm.Spec.Template.Spec.Volumes

	for i := range vm.Spec.Template.Spec.Volumes {
		if volumes[i].Name == dvName {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("disk %s is not attached to virtual machine %s/%s", dvName, namespace, name)
	}

	// detach disk
	newDisks := make([]kvcorev1.Disk, 0, len(vm.Spec.Template.Spec.Domain.Devices.Disks)-1)
	disks := vm.Spec.Template.Spec.Domain.Devices.Disks

	for i := range disks {
		if disks[i].Name != dvName {
			newDisks = append(newDisks, disks[i])
		}
	}

	vm.Spec.Template.Spec.Domain.Devices.Disks = newDisks

	// remove volume
	newVolumes := make([]kvcorev1.Volume, 0, len(vm.Spec.Template.Spec.Volumes)-1)
	volumes = vm.Spec.Template.Spec.Volumes

	for i := range volumes {
		if volumes[i].Name != dvName {
			newVolumes = append(newVolumes, volumes[i])
		}
	}

	vm.Spec.Template.Spec.Volumes = newVolumes

	// update virtual machine
	if _, err := uc.virtualMachine.Update(ctx, scope, namespace, vm); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) StartVirtualMachine(ctx context.Context, scope, namespace, name string) error {
	return uc.virtualMachine.Start(ctx, scope, namespace, name)
}

func (uc *UseCase) StopVirtualMachine(ctx context.Context, scope, namespace, name string) error {
	return uc.virtualMachine.Stop(ctx, scope, namespace, name)
}

func (uc *UseCase) RestartVirtualMachine(ctx context.Context, scope, namespace, name string) error {
	return uc.virtualMachine.Restart(ctx, scope, namespace, name)
}

func (uc *UseCase) filterClones(namespace, name string, clones []VirtualMachineClone) []VirtualMachineClone {
	ret := []VirtualMachineClone{}

	for i := range clones {
		if val, ok := clones[i].GetLabels()[nameLabel]; clones[i].Namespace == namespace && ok && val == name {
			ret = append(ret, clones[i])
		}
	}

	return ret
}

func (uc *UseCase) filterSnapshots(namespace, name string, snapshots []VirtualMachineSnapshot) []VirtualMachineSnapshot {
	ret := []VirtualMachineSnapshot{}

	for i := range snapshots {
		if val, ok := snapshots[i].GetLabels()[nameLabel]; snapshots[i].Namespace == namespace && ok && val == name {
			ret = append(ret, snapshots[i])
		}
	}

	return ret
}

func (uc *UseCase) filterRestores(namespace, name string, restores []VirtualMachineRestore) []VirtualMachineRestore {
	ret := []VirtualMachineRestore{}

	for i := range restores {
		if val, ok := restores[i].GetLabels()[nameLabel]; restores[i].Namespace == namespace && ok && val == name {
			ret = append(ret, restores[i])
		}
	}

	return ret
}

func (uc *UseCase) filterServices(namespace, name string, services []service.Service) []service.Service {
	ret := []service.Service{}

	for i := range services {
		if val, ok := services[i].GetLabels()[nameLabel]; services[i].Namespace == namespace && ok && val == name {
			ret = append(ret, services[i])
		}
	}

	return ret
}

func (uc *UseCase) combineVirtualMachines(virtualMachines []VirtualMachine, instances []vmi.VirtualMachineInstance, machines []machine.Machine, clones []VirtualMachineClone, snapshots []VirtualMachineSnapshot, restores []VirtualMachineRestore, services []service.Service) []VirtualMachineData {
	machineMap := map[string]*machine.Machine{}

	for i := range machines {
		m := machines[i]
		machineMap[m.Hostname] = &m
	}

	vmiMap := map[string]*vmi.VirtualMachineInstance{}

	for i := range instances {
		vmi := instances[i]
		vmiMap[vmi.Namespace+"/"+vmi.Name] = &vmi
	}

	ret := make([]VirtualMachineData, len(virtualMachines))

	for i := range virtualMachines {
		vm := virtualMachines[i]

		var machine *machine.Machine
		instance, ok := vmiMap[vm.Namespace+"/"+vm.Name]
		if ok {
			if nodeName := instance.Status.NodeName; nodeName != "" {
				machine = machineMap[nodeName]
			}
		}

		ret[i] = VirtualMachineData{
			VirtualMachine: &virtualMachines[i],
			Instance:       instance,
			Machine:        machine,
			Clones:         uc.filterClones(virtualMachines[i].Namespace, virtualMachines[i].Name, clones),
			Snapshots:      uc.filterSnapshots(virtualMachines[i].Namespace, virtualMachines[i].Name, snapshots),
			Restores:       uc.filterRestores(virtualMachines[i].Namespace, virtualMachines[i].Name, restores),
			Services:       uc.filterServices(virtualMachines[i].Namespace, virtualMachines[i].Name, services),
		}
	}

	return ret
}

func (uc *UseCase) combineVirtualMachine(namespace, name string, virtualMachine *VirtualMachine, instance *vmi.VirtualMachineInstance, machines []machine.Machine, clones []VirtualMachineClone, snapshots []VirtualMachineSnapshot, restores []VirtualMachineRestore, services []service.Service) *VirtualMachineData {
	machineMap := map[string]*machine.Machine{}

	for i := range machines {
		m := machines[i]
		machineMap[m.Hostname] = &m
	}

	var machine *machine.Machine
	if instance != nil {
		if nodeName := instance.Status.NodeName; nodeName != "" {
			machine = machineMap[nodeName]
		}
	}

	return &VirtualMachineData{
		VirtualMachine: virtualMachine,
		Instance:       instance,
		Machine:        machine,
		Clones:         uc.filterClones(namespace, name, clones),
		Snapshots:      uc.filterSnapshots(namespace, name, snapshots),
		Restores:       uc.filterRestores(namespace, name, restores),
		Services:       uc.filterServices(namespace, name, services),
	}
}

func (uc *UseCase) buildVirtualMachine(namespace, name, instanceType, bootDataVolume, startupScript string) *VirtualMachine {
	var (
		runStrategy   = kvcorev1.RunStrategyHalted
		enabled       = true
		bootOrder     = uint(1)
		osDisk        = "os-disk"
		cloudInitDisk = "cloud-init-disk"
		nic1          = "nic1"
	)

	virtualMachine := &kvcorev1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Annotations: map[string]string{
				"kubevirt.io/allow-pod-bridge-network-live-migration": "true",
			},
		},
		Spec: kvcorev1.VirtualMachineSpec{
			RunStrategy: &runStrategy,
			Instancetype: &kvcorev1.InstancetypeMatcher{
				Name: instanceType,
			},
			Template: &kvcorev1.VirtualMachineInstanceTemplateSpec{
				Spec: kvcorev1.VirtualMachineInstanceSpec{
					Domain: kvcorev1.DomainSpec{
						Devices: kvcorev1.Devices{
							Disks: []kvcorev1.Disk{
								{
									Name: osDisk,
									DiskDevice: kvcorev1.DiskDevice{
										Disk: &kvcorev1.DiskTarget{
											Bus: kvcorev1.DiskBusVirtio,
										},
									},
									BootOrder: &bootOrder,
								},
							},
							Interfaces: []kvcorev1.Interface{
								{
									Name: nic1,
									InterfaceBindingMethod: kvcorev1.InterfaceBindingMethod{
										Bridge: &kvcorev1.InterfaceBridge{},
									},
								},
							},
							TPM: &kvcorev1.TPMDevice{
								Enabled: &enabled,
							},
						},
					},
					Networks: []kvcorev1.Network{
						{
							Name: nic1,
							NetworkSource: kvcorev1.NetworkSource{
								Pod: &kvcorev1.PodNetwork{},
							},
						},
					},
					Volumes: []kvcorev1.Volume{
						{
							Name: osDisk,
							VolumeSource: kvcorev1.VolumeSource{
								DataVolume: &kvcorev1.DataVolumeSource{
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
		virtualMachine.Spec.Template.Spec.Domain.Devices.Disks = append(virtualMachine.Spec.Template.Spec.Domain.Devices.Disks, kvcorev1.Disk{
			Name: cloudInitDisk,
			DiskDevice: kvcorev1.DiskDevice{
				Disk: &kvcorev1.DiskTarget{
					Bus: kvcorev1.DiskBusVirtio,
				},
			},
		})
		virtualMachine.Spec.Template.Spec.Volumes = append(virtualMachine.Spec.Template.Spec.Volumes, kvcorev1.Volume{
			Name: cloudInitDisk,
			VolumeSource: kvcorev1.VolumeSource{
				CloudInitNoCloud: &kvcorev1.CloudInitNoCloudSource{
					UserData: startupScript,
				},
			},
		})
	}

	return virtualMachine
}
