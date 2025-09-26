import { writable, type Writable } from 'svelte/store';

import {
	VirtualMachine_Disk_Volume_Source_Type,
	VirtualMachine_Disk_Bus,
	// DataVolumeSource_Type,
} from '$lib/api/virtual_machine/v1/virtual_machine_pb';
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
		value: VirtualMachine_Disk_Volume_Source_Type.DATA_VOLUME,
		label: 'Data Volume',
		icon: 'ph:database',
	},
	{
		value: VirtualMachine_Disk_Volume_Source_Type.CLOUD_INIT_NO_CLOUD,
		label: 'Cloud Init No Cloud',
		icon: 'ph:cloud',
	},
]);

export const busTypes: Writable<SingleSelect.OptionType[]> = writable([
	{
		value: VirtualMachine_Disk_Bus.VIRTIO,
		label: 'VirtIO',
		icon: 'ph:cpu',
	},
	{
		value: VirtualMachine_Disk_Bus.SATA,
		label: 'SATA',
		icon: 'ph:hard-drive',
	},
	{
		value: VirtualMachine_Disk_Bus.SCSI,
		label: 'SCSI',
		icon: 'ph:hard-drives',
	},
	{
		value: VirtualMachine_Disk_Bus.USB,
		label: 'USB',
		icon: 'ph:usb',
	},
]);

// export const dataVolumeSourceTypes: Writable<SingleSelect.OptionType[]> = writable([
// 	{
// 		value: DataVolumeSource_Type.HTTP,
// 		label: 'HTTP',
// 		icon: 'ph:globe',
// 	},
// 	{
// 		value: DataVolumeSource_Type.BLANK,
// 		label: 'Blank',
// 		icon: 'ph:file',
// 	},
// 	{
// 		value: DataVolumeSource_Type.PVC,
// 		label: 'PVC',
// 		icon: 'ph:hard-drive',
// 	},
// ]);
