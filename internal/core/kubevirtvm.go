package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
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
	CreateVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string, labels map[string]string, spec *VirtualMachineCloneSpec) (*VirtualMachineClone, error)
	GetVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineClone, error)
	ListVirtualMachineClones(ctx context.Context, config *rest.Config, namespace string) ([]VirtualMachineClone, error)
	ListVirtualMachineClonesByVM(ctx context.Context, config *rest.Config, namespace, name string) ([]VirtualMachineClone, error)
	DeleteVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string) error
	CreateVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string, annotations, labels map[string]string, spec *VirtualMachineSnapshotSpec) (*VirtualMachineSnapshot, error)
	GetVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineSnapshot, error)
	ListVirtualMachineSnapshots(ctx context.Context, config *rest.Config, namespace string) ([]VirtualMachineSnapshot, error)
	ListVirtualMachineSnapshotsByVM(ctx context.Context, config *rest.Config, namespace, vmName string) ([]VirtualMachineSnapshot, error)
	DeleteVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string) error
	CreateVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string, labels map[string]string, spec *VirtualMachineRestoreSpec) (*VirtualMachineRestore, error)
	GetVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineRestore, error)
	ListVirtualMachineRestores(ctx context.Context, config *rest.Config, namespace string) ([]VirtualMachineRestore, error)
	ListVirtualMachineRestoresByVM(ctx context.Context, config *rest.Config, namespace, name string) ([]VirtualMachineRestore, error)
	DeleteVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string) error
	GetVirtualMachineMigrate(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineInstanceMigration, error)
	ListVirtualMachineMigrates(ctx context.Context, config *rest.Config, namespace string) ([]VirtualMachineInstanceMigration, error)
	ListVirtualMachineMigratesByVM(ctx context.Context, config *rest.Config, namespace, name string) ([]VirtualMachineInstanceMigration, error)
	DeleteVirtualMachineMigrate(ctx context.Context, config *rest.Config, namespace, name string) error
	// VirtualMachine Instances
	GetVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineInstance, error)
	ListVirtualMachineInstances(ctx context.Context, config *rest.Config, namespace string) ([]VirtualMachineInstance, error)
	UpdateVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string, spec *VirtualMachineInstanceSpec) (*VirtualMachineInstance, error)
	DeleteVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) error
	MigrateVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string, labels map[string]string, spec *VirtualMachineInstanceMigrationSpec) (*VirtualMachineInstanceMigration, error)
	PauseVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) error
	UnpauseVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) error
}

func (uc *KubeVirtUseCase) ListNamespaces(ctx context.Context, uuid, facility string) ([]corev1.Namespace, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeCore.ListNamespaces(ctx, config)
}

func (uc *KubeVirtUseCase) ListBootablePersistentVolumeClaims(ctx context.Context, uuid, facility, namespace string) ([]PersistentVolumeClaim, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	labelSelector := &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"otterscale.io/is_bootable": "true",
		},
	}
	selector, _ := metav1.LabelSelectorAsSelector(labelSelector)

	return uc.kubeCore.ListPersistentVolumeClaimsByLabel(ctx, config, namespace, selector.String())
}

func (uc *KubeVirtUseCase) CreateVirtualMachine(ctx context.Context, uuid, facility, namespace, name, network, script string, labels map[string]string, resources VirtualMachineResources, disks []DiskDevice, dataVolumeSources map[string]*DataVolumeInfo) (*VirtualMachine, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	if labels != nil {
		labels["otterscale.io/virtualmachine"] = name
	} else {
		labels = map[string]string{
			"otterscale.io/virtualmachine": name,
		}
	}

	annotations := map[string]string{
		"kubevirt.io/allow-pod-bridge-network-live-migration": "true",
	}

	// Create DataVolumes first if provided
	if dataVolumeSources != nil {
		for diskName, dvInfo := range dataVolumeSources {
			// Use disk name as DataVolume name if not specified
			dvName := dvInfo.Name
			if dvName == "" {
				dvName = diskName
			}

			// Create the DataVolume
			_, err := uc.kubeVirtDV.CreateDataVolume(
				ctx,
				config,
				namespace,
				dvName,
				dvInfo.SourceType,
				dvInfo.Source,
				name,
				dvInfo.SizeBytes,
				dvInfo.IsBootable,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create DataVolume %s: %w", dvName, err)
			}
		}
	}

	vmDisks, vmVolumes := buildDisksAndVolumes(disks, script)

	spec := buildVMSpec(resources, vmDisks, vmVolumes, network)

	return uc.kubeVirtVM.CreateVirtualMachine(ctx, config, namespace, name, labels, annotations, spec)
}

func volumeFromDisk(d DiskDevice) virtCorev1.Volume {
	switch strings.ToLower(d.DiskType) {
	case TYPEDATAVOLUME:
		return virtCorev1.Volume{
			Name: d.Name,
			VolumeSource: virtCorev1.VolumeSource{
				DataVolume: &virtCorev1.DataVolumeSource{Name: d.Name},
			},
		}
	case TYPEPERSISTENTVOLUMECLAIM:
		return virtCorev1.Volume{
			Name: d.Name,
			VolumeSource: virtCorev1.VolumeSource{
				PersistentVolumeClaim: &virtCorev1.PersistentVolumeClaimVolumeSource{
					PersistentVolumeClaimVolumeSource: corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: d.Name,
					},
				},
			},
		}
	case TYPECONFIGMAP:
		return virtCorev1.Volume{
			Name: d.Name,
			VolumeSource: virtCorev1.VolumeSource{
				ConfigMap: &virtCorev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: d.Name},
				},
			},
		}
	case TYPESECRET:
		return virtCorev1.Volume{
			Name: d.Name,
			VolumeSource: virtCorev1.VolumeSource{
				Secret: &virtCorev1.SecretVolumeSource{SecretName: d.Name},
			},
		}
	default:
		return virtCorev1.Volume{}
	}
}

func buildVMSpec(resources VirtualMachineResources, disks []virtCorev1.Disk, volumes []virtCorev1.Volume, network string) *VirtualMachineSpec {
	const (
		miB = 1024 * 1024
	)

	secureBoot := false
	hpet := false
	var retries uint32 = 8191
	strategy := virtCorev1.RunStrategyManual

	if network == "" {
		network = "default"
	}

	spec := &VirtualMachineSpec{
		RunStrategy: &strategy,
		Template: &virtCorev1.VirtualMachineInstanceTemplateSpec{
			Spec: virtCorev1.VirtualMachineInstanceSpec{
				Domain: virtCorev1.DomainSpec{
					Clock: &virtCorev1.Clock{
						Timer: &virtCorev1.Timer{
							HPET:   &virtCorev1.HPETTimer{Enabled: &hpet},
							Hyperv: &virtCorev1.HypervTimer{},
							PIT:    &virtCorev1.PITTimer{TickPolicy: virtCorev1.PITTickPolicyDelay},
							RTC:    &virtCorev1.RTCTimer{TickPolicy: virtCorev1.RTCTickPolicyCatchup},
						},
						ClockOffset: virtCorev1.ClockOffset{
							UTC: &virtCorev1.ClockOffsetUTC{},
						},
					},
					Devices: virtCorev1.Devices{
						Disks: disks,
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
							EFI: &virtCorev1.EFI{SecureBoot: &secureBoot},
						},
					},
				},
				Networks: []virtCorev1.Network{
					{
						Name: network,
						NetworkSource: virtCorev1.NetworkSource{
							Pod: &virtCorev1.PodNetwork{},
						},
					},
				},
				Volumes: volumes,
			},
		},
	}

	if resources.InstanceName != "" {
		spec.Instancetype = &virtCorev1.InstancetypeMatcher{Name: resources.InstanceName}
	} else {
		// CPU
		spec.Template.Spec.Domain.CPU = &virtCorev1.CPU{Cores: resources.CPUcores}
		// Memory
		mib := resources.MemoryBytes / miB
		memoryStr := fmt.Sprintf("%dMi", mib)
		spec.Template.Spec.Domain.Resources = virtCorev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(memoryStr),
			},
		}
	}
	return spec
}

func (uc *KubeVirtUseCase) ListVirtualMachines(ctx context.Context, uuid, facility, namespace string) ([]VirtualMachineInfo, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	vms, err := uc.kubeVirtVM.ListVirtualMachines(ctx, config, namespace)
	if err != nil {
		return nil, err
	}

	vmis, err := uc.kubeVirtVM.ListVirtualMachineInstances(ctx, config, namespace)
	if err != nil {
		return nil, err
	}

	machines, err := uc.machine.List(ctx)
	if err != nil {
		machines = []Machine{}
	}

	nodeToSystemID := make(map[string]string)
	for _, machine := range machines {
		nodeName := machine.FQDN
		if nodeName != "" {
			nodeToSystemID[nodeName] = machine.SystemID
		}
		if machine.Hostname != "" {
			nodeToSystemID[machine.Hostname] = machine.SystemID
		}
	}

	vmiMap := make(map[string]*VirtualMachineInstance)
	for i := range vmis {
		vmiMap[vmis[i].Name] = &vmis[i]
	}

	vmsInfo := make([]VirtualMachineInfo, 0, len(vms))

	for i := range vms {
		vm := vms[i]

		labelSelector := &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"vm.kubevirt.io/name": vm.Name,
			},
		}
		selector, _ := metav1.LabelSelectorAsSelector(labelSelector)

		// Pods
		podList, _ := uc.kubeCore.ListPodsByLabel(ctx, config, namespace, selector.String())
		labelSelector = &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"otterscale.io/virtualmachine": vm.Name,
			},
		}
		selector, _ = metav1.LabelSelectorAsSelector(labelSelector)
		// Services
		serviceList, _ := uc.kubeCore.ListServicesByOptions(ctx, config, namespace, selector.String(), "")

		dvNameSet := make(map[string]bool)
		if vm.Spec.Template != nil {
			for _, volume := range vm.Spec.Template.Spec.Volumes {
				if volume.DataVolume != nil && volume.DataVolume.Name != "" {
					dvNameSet[volume.DataVolume.Name] = true
				}
			}
		}

		// Datavolumes
		dvList, _ := uc.kubeVirtDV.ListDataVolumesByOptions(ctx, config, namespace, selector.String(), "")

		var dataVolumes []*DataVolume

		for j := range dvList {
			dataVolumes = append(dataVolumes, &dvList[j])
			dvNameSet[dvList[j].Name] = false
		}

		for dvName, needFetch := range dvNameSet {
			if !needFetch {
				continue
			}
			fieldSelector := fmt.Sprintf("metadata.name=%s", dvName)
			specificDvList, err := uc.kubeVirtDV.ListDataVolumesByOptions(ctx, config, namespace, "", fieldSelector)
			if err == nil && len(specificDvList) > 0 {
				dataVolumes = append(dataVolumes, &specificDvList[0])
			}
		}

		var pod *VirtualMachinePod
		if len(podList) > 0 {
			pod = &podList[0]
		}

		var service *VirtualMachineService
		if len(serviceList) > 0 {
			service = &serviceList[0]
		}

		// MAAS ID
		systemID := ""
		if pod != nil && pod.Spec.NodeName != "" {
			if id, ok := nodeToSystemID[pod.Spec.NodeName]; ok {
				systemID = id
			} else {
				shortName := pod.Spec.NodeName
				if idx := strings.Index(shortName, "."); idx > 0 {
					shortName = shortName[:idx]
				}
				if id, ok := nodeToSystemID[shortName]; ok {
					systemID = id
				}
			}
		}

		vmsInfo = append(vmsInfo, VirtualMachineInfo{
			VM:         &vms[i],
			VMI:        vmiMap[vm.Name],
			Pod:        pod,
			Service:    service,
			DataVolume: dataVolumes,
			SystemID:   systemID,
		})
	}

	return vmsInfo, nil
}

func (uc *KubeVirtUseCase) UpdateVirtualMachine(ctx context.Context, uuid, facility, namespace, name, networkName string, labels map[string]string, disks []DiskDevice) (vm *VirtualMachine, vmi *VirtualMachineInstance, err error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, nil, err
	}

	oldVM, err := uc.kubeVirtVM.GetVirtualMachine(ctx, config, namespace, name)
	if err != nil {
		return nil, nil, err
	}
	oldVM.SetLabels(ensureLabels(labels))

	if disks != nil {
		vmDisks, vmVolumes := buildDisksAndVolumes(disks, "")
		oldVM.Spec.Template.Spec.Domain.Devices.Disks = vmDisks
		oldVM.Spec.Template.Spec.Volumes = vmVolumes
	}

	updatedVM, err := uc.kubeVirtVM.UpdateVirtualMachine(ctx, config, namespace, name, oldVM)
	if err != nil {
		return nil, nil, err
	}

	vmi, err = uc.kubeVirtVM.GetVirtualMachineInstance(ctx, config, namespace, name)
	if err != nil {
		return updatedVM, nil, err
	}
	return updatedVM, vmi, nil
}

func ensureLabels(labels map[string]string) map[string]string {
	if labels == nil {
		labels = map[string]string{}
	}
	return labels
}

func buildDisksAndVolumes(disks []DiskDevice, script string) (vmDisks []virtCorev1.Disk, vmVolumes []virtCorev1.Volume) {
	for _, d := range disks {
		// ---------- Disk ----------
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

		vol := volumeFromDisk(d)
		if vol.Name == "" {
			continue
		}
		vmVolumes = append(vmVolumes, vol)
	}

	if strings.TrimSpace(script) != "" {
		vmVolumes = append(vmVolumes, virtCorev1.Volume{
			Name: cloudInitName,
			VolumeSource: virtCorev1.VolumeSource{
				CloudInitNoCloud: &virtCorev1.CloudInitNoCloudSource{UserData: script},
			},
		})
		vmDisks = append(vmDisks, virtCorev1.Disk{
			Name: cloudInitName,
			DiskDevice: virtCorev1.DiskDevice{
				Disk: &virtCorev1.DiskTarget{Bus: virtCorev1.DiskBusVirtio},
			},
		})
	}
	return vmDisks, vmVolumes
}

func (uc *KubeVirtUseCase) DeleteVirtualMachine(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	snapshots, err := uc.kubeVirtVM.ListVirtualMachineSnapshotsByVM(ctx, config, namespace, name)
	if err != nil {
		return fmt.Errorf("list snapshots failed: %w", err)
	}
	for i := range snapshots {
		if err := uc.kubeVirtVM.DeleteVirtualMachineSnapshot(ctx, config, namespace, snapshots[i].Name); err != nil && !apierrors.IsNotFound(err) {
			return fmt.Errorf("delete snapshot %s failed: %w", snapshots[i].Name, err)
		}
	}

	clones, err := uc.kubeVirtVM.ListVirtualMachineClonesByVM(ctx, config, namespace, name)
	if err != nil {
		return fmt.Errorf("list clones failed: %w", err)
	}
	for i := range clones {
		if err := uc.kubeVirtVM.DeleteVirtualMachineClone(ctx, config, namespace, clones[i].Name); err != nil && !apierrors.IsNotFound(err) {
			return fmt.Errorf("delete clone %s failed: %w", clones[i].Name, err)
		}
	}

	restores, err := uc.kubeVirtVM.ListVirtualMachineRestoresByVM(ctx, config, namespace, name)
	if err != nil {
		return fmt.Errorf("list restores failed: %w", err)
	}
	for i := range restores {
		if err := uc.kubeVirtVM.DeleteVirtualMachineRestore(ctx, config, namespace, restores[i].Name); err != nil && !apierrors.IsNotFound(err) {
			return fmt.Errorf("delete restore %s failed: %w", restores[i].Name, err)
		}
	}

	migrations, err := uc.kubeVirtVM.ListVirtualMachineMigratesByVM(ctx, config, namespace, name)
	if err != nil {
		return fmt.Errorf("list migrations failed: %w", err)
	}
	for i := range migrations {
		if err := uc.kubeVirtVM.DeleteVirtualMachineMigrate(ctx, config, namespace, migrations[i].Name); err != nil && !apierrors.IsNotFound(err) {
			return fmt.Errorf("delete migration %s failed: %w", migrations[i].Name, err)
		}
	}

	vm, err := uc.kubeVirtVM.GetVirtualMachine(ctx, config, namespace, name)
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("get vm failed: %w", err)
	}

	dataVolumeNames := make(map[string]bool)

	if err == nil && vm != nil && vm.Spec.Template != nil {
		for _, volume := range vm.Spec.Template.Spec.Volumes {
			if volume.DataVolume != nil && volume.DataVolume.Name != "" {
				dataVolumeNames[volume.DataVolume.Name] = true
			}
		}
	}

	labelSelector := &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"otterscale.io/virtualmachine": name,
		},
	}
	selector, _ := metav1.LabelSelectorAsSelector(labelSelector)

	// Services
	serviceList, err := uc.kubeCore.ListServicesByOptions(ctx, config, namespace, selector.String(), "")
	if err != nil {
		return fmt.Errorf("list services failed: %w", err)
	}
	for i := range serviceList {
		if err := uc.kubeCore.DeleteService(ctx, config, namespace, serviceList[i].Name); err != nil && !apierrors.IsNotFound(err) {
			return fmt.Errorf("delete service %s failed: %w", serviceList[i].Name, err)
		}
	}

	datavolumes, err := uc.kubeVirtDV.ListDataVolumesByOptions(ctx, config, namespace, selector.String(), "")
	if err != nil {
		return fmt.Errorf("list datavolumes by label failed: %w", err)
	}

	for i := range datavolumes {
		dataVolumeNames[datavolumes[i].Name] = true
	}

	for dvName := range dataVolumeNames {
		if err := uc.kubeVirtDV.DeleteDataVolume(ctx, config, namespace, dvName); err != nil && !apierrors.IsNotFound(err) {
			return fmt.Errorf("delete datavolume %s failed: %w", dvName, err)
		}
	}

	return uc.kubeVirtVM.DeleteVirtualMachine(ctx, config, namespace, name)
}

func (uc *KubeVirtUseCase) CreateVirtualMachineDisk(ctx context.Context, uuid, facility, namespace, vmName string, disk DiskDevice, dvInfo *DataVolumeInfo) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	if strings.ToLower(disk.DiskType) == TYPEDATAVOLUME && dvInfo != nil {
		dvName := dvInfo.Name
		if dvName == "" {
			dvName = disk.Name
		}
		_, err := uc.kubeVirtDV.CreateDataVolume(
			ctx,
			config,
			namespace,
			dvName,
			dvInfo.SourceType,
			dvInfo.Source,
			vmName,
			dvInfo.SizeBytes,
			dvInfo.IsBootable,
		)
		if err != nil {
			return err
		}
	}

	vm, err := uc.kubeVirtVM.GetVirtualMachine(ctx, config, namespace, vmName)
	if err != nil {
		return err
	}

	found := false
	for _, d := range vm.Spec.Template.Spec.Domain.Devices.Disks {
		if d.Name == disk.Name {
			found = true
			break
		}
	}
	if !found {
		var bus virtCorev1.DiskBus
		switch strings.ToLower(disk.Bus) {
		case "sata":
			bus = virtCorev1.DiskBusSATA
		case "scsi":
			bus = virtCorev1.DiskBusSCSI
		case "virtio":
			bus = virtCorev1.DiskBusVirtio
		default:
			bus = virtCorev1.DiskBusVirtio
		}
		newDisk := virtCorev1.Disk{
			Name: disk.Name,
			DiskDevice: virtCorev1.DiskDevice{
				Disk: &virtCorev1.DiskTarget{Bus: bus},
			},
		}
		vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks, newDisk)

		// append volume
		vol := volumeFromDisk(disk)
		if vol.Name != "" {
			vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, vol)
		}
	}

	// 4. 更新 VM
	_, err = uc.kubeVirtVM.UpdateVirtualMachine(ctx, config, namespace, vmName, vm)
	if err != nil {
		return err
	}
	return nil
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
func (uc *KubeVirtUseCase) CloneVirtualMachine(ctx context.Context, uuid, facility, targetNamespace, targetName, sourceNamespace, sourceName string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	labels := map[string]string{}

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
	_, err = uc.kubeVirtVM.CreateVirtualMachineClone(ctx, config, targetNamespace, targetName, labels, spec)
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

func (uc *KubeVirtUseCase) RestoreVirtualMachine(ctx context.Context, uuid, facility, namespace, name, snapshotName string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	labels := map[string]string{}

	spec := &VirtualMachineRestoreSpec{
		Target: corev1.TypedLocalObjectReference{
			APIGroup: &kubevirt,
			Kind:     "VirtualMachine",
			Name:     name,
		},
		VirtualMachineSnapshotName: snapshotName,
	}
	_, err = uc.kubeVirtVM.CreateVirtualMachineRestore(ctx, config, namespace, name, labels, spec)

	return err
}

func (uc *KubeVirtUseCase) MigrateVirtualMachine(ctx context.Context, uuid, facility, namespace, name, targetNode string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	labels := map[string]string{}

	spec := &VirtualMachineInstanceMigrationSpec{
		VMIName: name,
	}

	if targetNode != "" {
		spec.AddedNodeSelector = map[string]string{
			"kubernetes.io/hostname": targetNode,
		}
	}

	_, err = uc.kubeVirtVM.MigrateVirtualMachineInstance(ctx, config, namespace, name, labels, spec)

	return err
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

func (uc *KubeVirtUseCase) ListVirtualMachineSnapshots(ctx context.Context, uuid, facility, namespace, vmName string) ([]VirtualMachineSnapshot, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	ops, err := uc.kubeVirtVM.ListVirtualMachineSnapshotsByVM(ctx, config, namespace, vmName)
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
