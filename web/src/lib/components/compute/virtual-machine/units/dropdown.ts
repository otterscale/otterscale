import { writable, type Writable } from 'svelte/store';

import {
	VirtualMachineDisk_type,
	VirtualMachineDisk_bus,
	DataVolumeSource_Type,
} from '$lib/api/kubevirt/v1/kubevirt_pb';
import type { Single as SingleSelect } from '$lib/components/custom/select';

// ==================== Dropdown Options ====================

export const resourcesCase: Writable<SingleSelect.OptionType[]> = writable([
	{
		value: 'instancetypeName',
		label: 'InstanceType',
		icon: 'ph:layout',
	},
	{
		value: 'custom',
		label: 'Custom',
		icon: 'ph:gear',
	},
]);

export const diskTypes: Writable<SingleSelect.OptionType[]> = writable([
	{
		value: VirtualMachineDisk_type.DATAVOLUME,
		label: 'Data Volume',
		icon: 'ph:database',
	},
	{
		value: VirtualMachineDisk_type.PERSISTENTVOLUMECLAIM,
		label: 'Persistent Volume Claim',
		icon: 'ph:hard-drive',
	},
	{
		value: VirtualMachineDisk_type.CONFIGMAP,
		label: 'Config Map',
		icon: 'ph:file-text',
	},
	{
		value: VirtualMachineDisk_type.SECRET,
		label: 'Secret',
		icon: 'ph:lock',
	},
	{
		value: VirtualMachineDisk_type.CLOUDINITNOCLOUD,
		label: 'Cloud Init No Cloud',
		icon: 'ph:cloud',
	},
]);

export const busTypes: Writable<SingleSelect.OptionType[]> = writable([
	{
		value: VirtualMachineDisk_bus.VIRTIO,
		label: 'VirtIO',
		icon: 'ph:cpu',
	},
	{
		value: VirtualMachineDisk_bus.SATA,
		label: 'SATA',
		icon: 'ph:hard-drive',
	},
	{
		value: VirtualMachineDisk_bus.SCSI,
		label: 'SCSI',
		icon: 'ph:hard-drives',
	},
]);

export const dataVolumeSourceTypes: Writable<SingleSelect.OptionType[]> = writable([
	{
		value: DataVolumeSource_Type.HTTP,
		label: 'HTTP',
		icon: 'ph:globe',
	},
	{
		value: DataVolumeSource_Type.BLANK,
		label: 'Blank',
		icon: 'ph:file',
	},
	// {
	// 	value: DataVolumeSource_Type.REGISTRY,
	// 	label: 'Registry',
	// 	icon: 'ph:package',
	// },
	// {
	// 	value: DataVolumeSource_Type.UPLOAD,
	// 	label: 'Upload',
	// 	icon: 'ph:upload',
	// },
	// {
	// 	value: DataVolumeSource_Type.S3,
	// 	label: 'S3',
	// 	icon: 'ph:cloud',
	// },
	// {
	// 	value: DataVolumeSource_Type.VDDK,
	// 	label: 'VDDK',
	// 	icon: 'ph:hard-drives',
	// },
	{
		value: DataVolumeSource_Type.PVC,
		label: 'PVC',
		icon: 'ph:hard-drive',
	},
]);
