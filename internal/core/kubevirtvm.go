package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/rest"
	virtCorev1 "kubevirt.io/api/core/v1"
)

var kubevirt = "kubevirt.io"

const (
	TYPEDATAVOLUME            = "datavolume"
	TYPEPERSISTENTVOLUMECLAIM = "persistentvolumeclaim"
	TYPECONFIGMAP             = "configmap"
	TYPESECRET                = "secret"
	TYPEHOSTDISK              = "hostdisk"
	TYPECLOUDINITNOCLOUD      = "cloudinitnocloud"
	TYPECONTAINERDISK         = "containerdisk"
	cloudInitName             = "cloudinit-no-cloud"
)

type KubeVirtVMRepo interface {
	// VirtualMachines
	CreateVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string, labels, annotations map[string]string, spec *VirtualMachineSpec) (*VirtualMachine, error)
	GetVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachine, error)
	ListVirtualMachines(ctx context.Context, config *rest.Config, namespace string) ([]VirtualMachine, error)
	UpdateVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string, vm *VirtualMachine) (*VirtualMachine, error)
	DeleteVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	MigrateVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	StartVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	RestartVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	StopVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	CreateVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string, annotations, labels map[string]string, spec *VirtualMachineCloneSpec) (*VirtualMachineClone, error)
	GetVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineClone, error)
	ListVirtualMachineClones(ctx context.Context, config *rest.Config, namespace string) ([]VirtualMachineClone, error)
	DeleteVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string) error
	CreateVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string, annotations, labels map[string]string, spec *VirtualMachineSnapshotSpec) (*VirtualMachineSnapshot, error)
	GetVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineSnapshot, error)
	ListVirtualMachineSnapshots(ctx context.Context, config *rest.Config, namespace string) ([]VirtualMachineSnapshot, error)
	DeleteVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string) error
	CreateVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string, annotations, labels map[string]string, spec *VirtualMachineRestoreSpec) (*VirtualMachineRestore, error)
	GetVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineRestore, error)
	ListVirtualMachineRestores(ctx context.Context, config *rest.Config, namespace string) ([]VirtualMachineRestore, error)
	DeleteVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string) error
	// CreateVirtualMachineMigrate(ctx context.Context, config *rest.Config, namespace, name string, annotations, labels map[string]string, sp) (*VirtualMachineMigrate, error)
	GetVirtualMachineMigrate(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineInstanceMigration, error)
	ListVirtualMachineMigrates(ctx context.Context, config *rest.Config, namespace string) ([]VirtualMachineInstanceMigration, error)
	DeleteVirtualMachineMigrate(ctx context.Context, config *rest.Config, namespace, name string) error
	// VirtualMachine Instances
	GetVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineInstance, error)
	ListVirtualMachineInstances(ctx context.Context, config *rest.Config, namespace string) ([]VirtualMachineInstance, error)
	UpdateVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string, spec *VirtualMachineInstanceSpec) (*VirtualMachineInstance, error)
	DeleteVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) error
	MigrateVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string, annotations, labels map[string]string, spec *VirtualMachineInstanceMigrationSpec) (*VirtualMachineInstanceMigration, error)
	PauseVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) error
	UnpauseVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) error
}

func (uc *KubeVirtUseCase) CreateVirtualMachine(ctx context.Context, uuid, facility, namespace, name, network, script string, labels, annotations map[string]string, resources VirtualMachineResources, disks []DiskDevice) (*VirtualMachine, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	// generate disk and volume
	var vmDisks []virtCorev1.Disk
	var vmVolumes []virtCorev1.Volume

	if labels != nil {
		labels["kubevirt.io/vm"] = name
	} else {
		labels = map[string]string{
			"kubevirt.io/vm": name,
		}
	}

	if annotations != nil {
		annotations["kubevirt.otterscale.io/network"] = network
		annotations["kubevirt.otterscale.io/script"] = script
		annotations["kubevirt.io/allow-pod-bridge-network-live-migration"] = "true"
	} else {
		annotations = map[string]string{
			"kubevirt.otterscale.io/network":                      network,
			"otterscale/virtualmachine/script":                    script,
			"kubevirt.io/allow-pod-bridge-network-live-migration": "true",
		}
	}

	for _, d := range disks {
		var bus virtCorev1.DiskBus
		switch strings.ToLower(d.Bus) {
		case "sata":
			bus = virtCorev1.DiskBusSATA
		case "scsi":
			bus = virtCorev1.DiskBusSCSI
		case "virtio":
			bus = virtCorev1.DiskBusVirtio
		default:
			// default using virtio
			bus = virtCorev1.DiskBusVirtio
		}
		vmDisks = append(vmDisks, virtCorev1.Disk{
			Name: d.Name,
			DiskDevice: virtCorev1.DiskDevice{
				Disk: &virtCorev1.DiskTarget{
					Bus: bus,
				},
			},
		})

		// Create VolumeSource according to disktype
		switch strings.ToLower(d.DiskType) {
		case TYPEDATAVOLUME: // protobuf enum => DATAVOLUME
			vmVolumes = append(vmVolumes, virtCorev1.Volume{
				Name: d.Name,
				VolumeSource: virtCorev1.VolumeSource{
					DataVolume: &virtCorev1.DataVolumeSource{
						Name: d.Name,
					},
				},
			})
		case TYPEPERSISTENTVOLUMECLAIM:
			vmVolumes = append(vmVolumes, virtCorev1.Volume{
				Name: d.Name,
				VolumeSource: virtCorev1.VolumeSource{
					PersistentVolumeClaim: &virtCorev1.PersistentVolumeClaimVolumeSource{
						PersistentVolumeClaimVolumeSource: corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: d.Name,
						},
					},
				},
			})
		case TYPECONFIGMAP:
			vmVolumes = append(vmVolumes, virtCorev1.Volume{
				Name: d.Name,
				VolumeSource: virtCorev1.VolumeSource{
					ConfigMap: &virtCorev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: d.Name,
						},
					},
				},
			})
		case TYPESECRET:
			vmVolumes = append(vmVolumes, virtCorev1.Volume{
				Name: d.Name,
				VolumeSource: virtCorev1.VolumeSource{
					Secret: &virtCorev1.SecretVolumeSource{
						SecretName: d.Name,
					},
				},
			})
		case TYPECLOUDINITNOCLOUD:
			vmVolumes = append(vmVolumes, virtCorev1.Volume{
				Name: d.Name,
				VolumeSource: virtCorev1.VolumeSource{
					CloudInitNoCloud: &virtCorev1.CloudInitNoCloudSource{
						UserDataBase64: d.Data,
					},
				},
			})
		case TYPECONTAINERDISK:
			vmVolumes = append(vmVolumes, virtCorev1.Volume{
				Name: d.Name,
				VolumeSource: virtCorev1.VolumeSource{
					ContainerDisk: &virtCorev1.ContainerDiskSource{
						Image: d.Data,
					},
				},
			})
		default:
			// [TODO] Not Support
		}
	}

	// if script is given，add CloudInitNoCloud volume & disk
	if strings.TrimSpace(script) != "" {
		vmVolumes = append(vmVolumes, virtCorev1.Volume{
			Name: cloudInitName,
			VolumeSource: virtCorev1.VolumeSource{
				CloudInitNoCloud: &virtCorev1.CloudInitNoCloudSource{
					UserData: script,
				},
			},
		})

		vmDisks = append(vmDisks, virtCorev1.Disk{
			Name: cloudInitName,
			DiskDevice: virtCorev1.DiskDevice{
				Disk: &virtCorev1.DiskTarget{
					Bus: virtCorev1.DiskBusVirtio,
				},
			},
		})
	}

	secureBoot := false
	hpet := false
	var retries uint32 = 8191
	strategy := virtCorev1.RunStrategyManual

	spec := &VirtualMachineSpec{
		RunStrategy: &strategy,
		Template: &virtCorev1.VirtualMachineInstanceTemplateSpec{
			Spec: virtCorev1.VirtualMachineInstanceSpec{
				Domain: virtCorev1.DomainSpec{
					// Clock
					Clock: &virtCorev1.Clock{
						Timer: &virtCorev1.Timer{
							HPET:   &virtCorev1.HPETTimer{Enabled: &hpet},
							Hyperv: &virtCorev1.HypervTimer{},
							PIT: &virtCorev1.PITTimer{
								TickPolicy: virtCorev1.PITTickPolicyDelay,
							},
							RTC: &virtCorev1.RTCTimer{
								TickPolicy: virtCorev1.RTCTickPolicyCatchup,
							},
						},
						ClockOffset: virtCorev1.ClockOffset{
							UTC: &virtCorev1.ClockOffsetUTC{},
						},
					},
					Devices: virtCorev1.Devices{
						Disks: vmDisks,
						Interfaces: []virtCorev1.Interface{
							{
								Name:  network,
								Model: "e1000",
								InterfaceBindingMethod: virtCorev1.InterfaceBindingMethod{
									Masquerade: &virtCorev1.InterfaceMasquerade{},
								},
							},
						},
						TPM: &virtCorev1.TPMDevice{},
					},
					Features: &virtCorev1.Features{
						ACPI: virtCorev1.FeatureState{},
						APIC: &virtCorev1.FeatureAPIC{},
						Hyperv: &virtCorev1.FeatureHyperv{
							Relaxed:   &virtCorev1.FeatureState{},
							Spinlocks: &virtCorev1.FeatureSpinlocks{Retries: &retries},
							VAPIC:     &virtCorev1.FeatureState{},
						},
					},
					Firmware: &virtCorev1.Firmware{
						Bootloader: &virtCorev1.Bootloader{
							EFI: &virtCorev1.EFI{
								SecureBoot: &secureBoot,
							},
						},
					},
				},
				Networks: []virtCorev1.Network{
					{
						Name: "default",
						NetworkSource: virtCorev1.NetworkSource{
							Pod: &virtCorev1.PodNetwork{},
						},
					},
				},
				Volumes: vmVolumes,
			},
		},
	}

	// If instancetype is given, use instancetype
	if resources.InstanceName != "" {
		spec.Instancetype = &virtCorev1.InstancetypeMatcher{
			Name: resources.InstanceName,
		}
	} else {
		spec.Template.Spec.Domain.CPU = &virtCorev1.CPU{
			Cores: uint32(resources.CPUcores),
		}

		memoryStr := fmt.Sprintf("%dB", resources.MemoryBytes)

		const MiB = 1024 * 1024
		memoryMiB := resources.MemoryBytes / MiB
		memoryStr = fmt.Sprintf("%dMi", memoryMiB)

		spec.Template.Spec.Domain.Resources = virtCorev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(memoryStr),
			},
		}
	}

	return uc.kubeVirtVM.CreateVirtualMachine(ctx, config, namespace, name, labels, annotations, spec)
}

func (uc *KubeVirtUseCase) GetVirtualMachine(ctx context.Context, uuid, facility, namespace, name string) (*VirtualMachine, *VirtualMachineInstance, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, nil, err
	}

	vm, err := uc.kubeVirtVM.GetVirtualMachine(ctx, config, name, namespace)
	if err != nil {
		return nil, nil, err
	}
	vmi, err := uc.kubeVirtVM.GetVirtualMachineInstance(ctx, config, name, namespace)
	if err != nil {
		return nil, nil, err
	}
	return vm, vmi, err
}

func (uc *KubeVirtUseCase) ListVirtualMachines(ctx context.Context, uuid, facility, namespace string) ([]VirtualMachine, []VirtualMachineInstance, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, nil, err
	}

	vms, err := uc.kubeVirtVM.ListVirtualMachines(ctx, config, namespace)
	if err != nil {
		return nil, nil, err
	}

	vmis, err := uc.kubeVirtVM.ListVirtualMachineInstances(ctx, config, namespace)
	if err != nil {
		return nil, nil, err
	}
	return vms, vmis, err
}

func (uc *KubeVirtUseCase) UpdateVirtualMachine(ctx context.Context, uuid, facility, namespace, name, networkName, startupScript string, labels, annotations map[string]string, disks []DiskDevice) (*VirtualMachine, *VirtualMachineInstance, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, nil, err
	}

	oldVM, err := uc.kubeVirtVM.GetVirtualMachine(ctx, config, namespace, name)
	if err != nil {
		return nil, nil, err
	}

	if labels == nil {
		labels = map[string]string{}
	}
	oldVM.SetLabels(labels)

	if annotations == nil {
		annotations = map[string]string{}
	}
	oldVM.SetAnnotations(annotations)

	var vmDisks []virtCorev1.Disk
	var vmVolumes []virtCorev1.Volume

	for _, d := range disks {
		// --- Disk Bus ---
		var bus virtCorev1.DiskBus
		switch strings.ToLower(d.Bus) {
		case "sata":
			bus = virtCorev1.DiskBusSATA
		case "scsi":
			bus = virtCorev1.DiskBusSCSI
		case "virtio":
			bus = virtCorev1.DiskBusVirtio
		default:
			bus = virtCorev1.DiskBusVirtio
		}

		vmDisks = append(vmDisks, virtCorev1.Disk{
			Name: d.Name,
			DiskDevice: virtCorev1.DiskDevice{
				Disk: &virtCorev1.DiskTarget{Bus: bus},
			},
		})

		switch strings.ToLower(d.DiskType) {
		case TYPEDATAVOLUME:
			vmVolumes = append(vmVolumes, virtCorev1.Volume{
				Name: d.Name,
				VolumeSource: virtCorev1.VolumeSource{
					DataVolume: &virtCorev1.DataVolumeSource{
						Name: d.Name,
					},
				},
			})
		case TYPEPERSISTENTVOLUMECLAIM:
			vmVolumes = append(vmVolumes, virtCorev1.Volume{
				Name: d.Name,
				VolumeSource: virtCorev1.VolumeSource{
					PersistentVolumeClaim: &virtCorev1.PersistentVolumeClaimVolumeSource{
						PersistentVolumeClaimVolumeSource: corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: d.Name,
						},
					},
				},
			})
		case TYPECONFIGMAP:
			vmVolumes = append(vmVolumes, virtCorev1.Volume{
				Name: d.Name,
				VolumeSource: virtCorev1.VolumeSource{
					ConfigMap: &virtCorev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: d.Name,
						},
					},
				},
			})
		case TYPESECRET:
			vmVolumes = append(vmVolumes, virtCorev1.Volume{
				Name: d.Name,
				VolumeSource: virtCorev1.VolumeSource{
					Secret: &virtCorev1.SecretVolumeSource{
						SecretName: d.Name,
					},
				},
			})
		case TYPECLOUDINITNOCLOUD:
			vmVolumes = append(vmVolumes, virtCorev1.Volume{
				Name: d.Name,
				VolumeSource: virtCorev1.VolumeSource{
					CloudInitNoCloud: &virtCorev1.CloudInitNoCloudSource{
						UserDataBase64: d.Data,
					},
				},
			})
		case TYPECONTAINERDISK:
			vmVolumes = append(vmVolumes, virtCorev1.Volume{
				Name: d.Name,
				VolumeSource: virtCorev1.VolumeSource{
					ContainerDisk: &virtCorev1.ContainerDiskSource{
						Image: d.Data,
					},
				},
			})
		}
	}

	if strings.TrimSpace(startupScript) != "" {
		cloudVol := virtCorev1.Volume{
			Name: cloudInitName,
			VolumeSource: virtCorev1.VolumeSource{
				CloudInitNoCloud: &virtCorev1.CloudInitNoCloudSource{UserData: startupScript},
			},
		}
		cloudDisk := virtCorev1.Disk{
			Name: cloudInitName,
			DiskDevice: virtCorev1.DiskDevice{
				Disk: &virtCorev1.DiskTarget{Bus: virtCorev1.DiskBusVirtio},
			},
		}
		vmVolumes = append(vmVolumes, cloudVol)
		vmDisks = append(vmDisks, cloudDisk)
	}

	oldVM.Spec.Template.Spec.Domain.Devices.Disks = vmDisks
	oldVM.Spec.Template.Spec.Volumes = vmVolumes

	updatedVM, err := uc.kubeVirtVM.UpdateVirtualMachine(ctx, config, namespace, name, oldVM)
	if err != nil {
		return nil, nil, err
	}

	vmi, err := uc.kubeVirtVM.GetVirtualMachineInstance(ctx, config, namespace, name)
	if err != nil {
		return updatedVM, nil, err
	}
	return updatedVM, vmi, nil
}

func (uc *KubeVirtUseCase) DeleteVirtualMachine(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtVM.DeleteVirtualMachine(ctx, config, namespace, name)
}

// Virtual Machine Control Operations
func (uc *KubeVirtUseCase) StartVirtualMachine(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtVM.StartVirtualMachine(ctx, config, namespace, name)
}

func (uc *KubeVirtUseCase) StopVirtualMachine(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtVM.StopVirtualMachine(ctx, config, namespace, name)
}

func (uc *KubeVirtUseCase) PauseVirtualMachine(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtVM.PauseVirtualMachineInstance(ctx, config, namespace, name)
}

func (uc *KubeVirtUseCase) UnpauseVirtualMachine(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtVM.UnpauseVirtualMachineInstance(ctx, config, namespace, name)
}

// Virtual Machine Advanced Operations
func (uc *KubeVirtUseCase) CloneVirtualMachine(ctx context.Context, uuid, facility, targetNamespace, targetName, sourceNamespace, sourceName, description string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	labels := map[string]string{}
	annotations := map[string]string{
		"otterscale.io/clone-description": description,
	}

	spec := &VirtualMachineCloneSpec{
		Source: &corev1.TypedLocalObjectReference{
			APIGroup: &kubevirt,
			Kind:     "VirtualMachine",
			Name:     sourceName,
		},
		Target: &corev1.TypedLocalObjectReference{
			APIGroup: &kubevirt,
			Kind:     "VirtualMachine",
			Name:     targetName,
		},
	}
	_, err = uc.kubeVirtVM.CreateVirtualMachineClone(ctx, config, targetNamespace, targetName, annotations, labels, spec)
	return err
}

func (uc *KubeVirtUseCase) SnapshotVirtualMachine(ctx context.Context, uuid, facility, namespace, name, snapshotName, description string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	labels := map[string]string{}
	annotations := map[string]string{
		"otterscale.io/snapshot-description": description,
	}

	spec := &VirtualMachineSnapshotSpec{
		Source: corev1.TypedLocalObjectReference{
			APIGroup: &kubevirt,
			Kind:     "VirtualMachine",
			Name:     name,
		},
		FailureDeadline: &metav1.Duration{
			Duration: 5 * time.Minute,
		},
	}
	_, err = uc.kubeVirtVM.CreateVirtualMachineSnapshot(ctx, config, namespace, snapshotName, annotations, labels, spec)

	return err
}

func (uc *KubeVirtUseCase) RestoreVirtualMachine(ctx context.Context, uuid, facility, namespace, name, snapshotName, description string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	labels := map[string]string{}
	annotations := map[string]string{
		"otterscale.io/restore-description": description,
	}

	spec := &VirtualMachineRestoreSpec{
		Target: corev1.TypedLocalObjectReference{
			APIGroup: &kubevirt,
			Kind:     "VirtualMachine",
			Name:     name,
		},
		VirtualMachineSnapshotName: snapshotName,
	}
	_, err = uc.kubeVirtVM.CreateVirtualMachineRestore(ctx, config, namespace, name, annotations, labels, spec)

	return err
}

func (uc *KubeVirtUseCase) MigrateVirtualMachine(ctx context.Context, uuid, facility, namespace, name, targetNode, description string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	return uc.kubeVirtVM.MigrateVirtualMachine(ctx, config, namespace, name)
}

// GetVirtualMachineClone retrieves a virtual machine clone
func (uc *KubeVirtUseCase) GetVirtualMachineClone(ctx context.Context, uuid, facility, namespace, name string) (*VirtualMachineClone, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	clone, err := uc.kubeVirtVM.GetVirtualMachineClone(ctx, config, namespace, name)
	if err != nil {
		return nil, err
	}

	return clone, nil
}

func (uc *KubeVirtUseCase) ListVirtualMachineClones(ctx context.Context, uuid, facility, namespace string) ([]VirtualMachineClone, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	clones, err := uc.kubeVirtVM.ListVirtualMachineClones(ctx, config, namespace)
	if err != nil {
		return nil, err
	}

	return clones, nil
}

// DeleteVirtualMachineClone deletes a virtual machine clone
func (uc *KubeVirtUseCase) DeleteVirtualMachineClone(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	return uc.kubeVirtVM.DeleteVirtualMachineClone(ctx, config, namespace, name)
}

// GetVirtualMachineSnapshot retrieves a virtual machine snapshot
func (uc *KubeVirtUseCase) GetVirtualMachineSnapshot(ctx context.Context, uuid, facility, namespace, name string) (*VirtualMachineSnapshot, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	snapshot, err := uc.kubeVirtVM.GetVirtualMachineSnapshot(ctx, config, namespace, name)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

func (uc *KubeVirtUseCase) ListVirtualMachineSnapshots(ctx context.Context, uuid, facility, namespace string) ([]VirtualMachineSnapshot, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	ops, err := uc.kubeVirtVM.ListVirtualMachineSnapshots(ctx, config, namespace)
	if err != nil {
		return nil, err
	}

	return ops, nil
}

// DeleteVirtualMachineSnapshot deletes a virtual machine snapshot
func (uc *KubeVirtUseCase) DeleteVirtualMachineSnapshot(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	return uc.kubeVirtVM.DeleteVirtualMachineSnapshot(ctx, config, namespace, name)
}

// GetVirtualMachineRestore retrieves a virtual machine restore operation
func (uc *KubeVirtUseCase) GetVirtualMachineRestore(ctx context.Context, uuid, facility, namespace, name string) (*VirtualMachineRestore, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	restore, err := uc.kubeVirtVM.GetVirtualMachineRestore(ctx, config, namespace, name)
	if err != nil {
		return nil, err
	}

	return restore, nil
}

func (uc *KubeVirtUseCase) ListVirtualMachineRestores(ctx context.Context, uuid, facility, namespace string) ([]VirtualMachineRestore, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	restores, err := uc.kubeVirtVM.ListVirtualMachineRestores(ctx, config, namespace)
	if err != nil {
		return nil, err
	}

	return restores, nil
}

// DeleteVirtualMachineRestore deletes a virtual machine restore operation
func (uc *KubeVirtUseCase) DeleteVirtualMachineRestore(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	return uc.kubeVirtVM.DeleteVirtualMachineRestore(ctx, config, namespace, name)
}

// GetVirtualMachineMigrate retrieves a virtual machine Migrate operation
func (uc *KubeVirtUseCase) GetVirtualMachineMigrate(ctx context.Context, uuid, facility, namespace, name string) (*VirtualMachineInstanceMigration, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	Migrate, err := uc.kubeVirtVM.GetVirtualMachineMigrate(ctx, config, namespace, name)
	if err != nil {
		return nil, err
	}

	return Migrate, nil
}

func (uc *KubeVirtUseCase) ListVirtualMachineMigrates(ctx context.Context, uuid, facility, namespace string) ([]VirtualMachineInstanceMigration, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	Migrates, err := uc.kubeVirtVM.ListVirtualMachineMigrates(ctx, config, namespace)
	if err != nil {
		return nil, err
	}

	return Migrates, nil
}

// DeleteVirtualMachineMigrate deletes a virtual machine Migrate operation
func (uc *KubeVirtUseCase) DeleteVirtualMachineMigrate(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	return uc.kubeVirtVM.DeleteVirtualMachineMigrate(ctx, config, namespace, name)
}
