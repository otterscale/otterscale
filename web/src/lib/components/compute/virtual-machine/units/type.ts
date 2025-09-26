import type { VirtualMachine_Disk } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
import { m } from '$lib/paraglide/messages';

export type EnhancedDisk = VirtualMachine_Disk & {
	phase?: string;
	bootImage?: boolean;
	sizeBytes?: bigint;
};

export interface StatusInfo {
	icon: string;
	color: string;
	text: string;
}

export function getStatusInfo(status: string): StatusInfo {
	switch (status) {
		case 'Running':
			return {
				icon: 'ph:power',
				color: 'text-green-600',
				text: m.running(),
			};
		case 'Stopped':
			return {
				icon: 'ph:power',
				color: 'text-gray-600',
				text: m.stopped(),
			};
		case 'Paused':
			return {
				icon: 'ph:pause-circle',
				color: 'text-yellow-600',
				text: m.paused(),
			};
		case 'Starting':
			return {
				icon: 'ph:arrow-clockwise',
				color: 'text-blue-500 animate-spin',
				text: m.starting(),
			};
		case 'Stopping':
			return {
				icon: 'ph:stop-circle',
				color: 'text-orange-500 animate-pulse',
				text: m.stopping(),
			};
		case 'Provisioning':
			return {
				icon: 'ph:hourglass-high',
				color: 'text-blue-500',
				text: m.provisioning(),
			};
		case 'Migrating':
			return {
				icon: 'ph:arrows-left-right',
				color: 'text-purple-500 animate-pulse',
				text: m.migrating(),
			};
		case 'CrashLoopBackOff':
			return {
				icon: 'ph:warning-octagon',
				color: 'text-red-600 animate-pulse',
				text: m.crash_loop_back_off(),
			};
		case 'ErrorUnschedulable':
			return {
				icon: 'ph:prohibit',
				color: 'text-red-500',
				text: m.error_unschedulable(),
			};
		case 'ErrImagePull':
			return {
				icon: 'ph:download-simple',
				color: 'text-red-500',
				text: m.err_image_pull(),
			};
		case 'ImagePullBackOff':
			return {
				icon: 'ph:download-simple',
				color: 'text-red-500 animate-pulse',
				text: m.image_pull_back_off(),
			};
		case 'ErrorPvcNotFound':
			return {
				icon: 'ph:hard-drive',
				color: 'text-red-500',
				text: m.error_pvc_not_found(),
			};
		case 'DataVolumeError':
			return {
				icon: 'ph:database',
				color: 'text-red-500',
				text: m.data_volume_error(),
			};
		case 'WaitingForVolumeBinding':
			return {
				icon: 'ph:link',
				color: 'text-yellow-500 animate-pulse',
				text: m.waiting_for_volume_binding(),
			};
		case 'WaitingForReceiver':
			return {
				icon: 'ph:clock',
				color: 'text-blue-500 animate-pulse',
				text: m.waiting_for_receiver(),
			};
		case 'Unknown':
		default:
			return {
				icon: 'ph:warning-circle-fill',
				color: 'text-amber-500',
				text: m.unknown(),
			};
	}
}

export function getVolumeStatusInfo(phase: string): StatusInfo {
	switch (phase) {
		case 'Bound':
			return {
				icon: 'ph:check-circle',
				color: 'text-green-600',
				text: m.bound(),
			};
		case 'Pending':
			return {
				icon: 'ph:clock',
				color: 'text-yellow-500 animate-pulse',
				text: m.pending(),
			};
		case 'Available':
			return {
				icon: 'ph:check',
				color: 'text-blue-500',
				text: m.available(),
			};
		case 'Released':
			return {
				icon: 'ph:arrow-circle-up',
				color: 'text-gray-500',
				text: m.released(),
			};
		case 'Failed':
			return {
				icon: 'ph:x-circle',
				color: 'text-red-500',
				text: m.failed(),
			};
		case 'Terminating':
			return {
				icon: 'ph:trash',
				color: 'text-red-500 animate-pulse',
				text: m.terminating(),
			};
		case 'Attaching':
			return {
				icon: 'ph:link',
				color: 'text-blue-500 animate-spin',
				text: m.attaching(),
			};
		case 'Detaching':
			return {
				icon: 'ph:link-break',
				color: 'text-orange-500 animate-pulse',
				text: m.detaching(),
			};
		case 'Resizing':
			return {
				icon: 'ph:arrows-out',
				color: 'text-purple-500 animate-pulse',
				text: m.resizing(),
			};
		case 'Ready':
			return {
				icon: 'ph:check-circle-fill',
				color: 'text-green-500',
				text: m.ready(),
			};
		case 'NotReady':
			return {
				icon: 'ph:warning-circle',
				color: 'text-yellow-600',
				text: m.not_ready(),
			};
		case 'Lost':
			return {
				icon: 'ph:question-mark',
				color: 'text-red-600',
				text: m.lost(),
			};
		case 'Unknown':
		default:
			return {
				icon: 'ph:warning-circle-fill',
				color: 'text-amber-500',
				text: m.unknown(),
			};
	}
}

export function getInstancePhaseInfo(phase: string): StatusInfo {
	switch (phase) {
		case 'VmPhaseUnset':
			return {
				icon: 'ph:circle-dashed',
				color: 'text-gray-400',
				text: m.unset(),
			};
		case 'Pending':
			return {
				icon: 'ph:clock',
				color: 'text-yellow-500 animate-pulse',
				text: m.pending(),
			};
		case 'Scheduling':
			return {
				icon: 'ph:calendar-check',
				color: 'text-blue-500 animate-pulse',
				text: m.scheduling(),
			};
		case 'Scheduled':
			return {
				icon: 'ph:check-square',
				color: 'text-blue-600',
				text: m.scheduled(),
			};
		case 'Running':
			return {
				icon: 'ph:play-circle',
				color: 'text-green-600',
				text: m.running(),
			};
		case 'Succeeded':
			return {
				icon: 'ph:check-circle-fill',
				color: 'text-green-500',
				text: m.succeeded(),
			};
		case 'Failed':
			return {
				icon: 'ph:x-circle-fill',
				color: 'text-red-500',
				text: m.failed(),
			};
		case 'WaitingForSync':
			return {
				icon: 'ph:arrows-clockwise',
				color: 'text-purple-500 animate-spin',
				text: m.waiting_for_sync(),
			};
		case 'Unknown':
		default:
			return {
				icon: 'ph:warning-circle-fill',
				color: 'text-amber-500',
				text: m.unknown(),
			};
	}
}
