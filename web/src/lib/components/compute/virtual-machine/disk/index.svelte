<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { SvelteMap } from 'svelte/reactivity';

	import { DataTable } from './data-table';

	import {
		VirtualMachineService,
		VirtualMachine_Disk_Volume_Source_Type,
	} from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import type { VirtualMachine } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import type { DataVolume } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import type { EnhancedDisk } from '$lib/components/compute/virtual-machine/units/type';
	import * as Sheet from '$lib/components/ui/sheet';
</script>

<script lang="ts">
	let {
		virtualMachine,
	}: {
		virtualMachine: VirtualMachine;
	} = $props();

	// Context dependencies
	const transport: Transport = getContext('transport');
	const virtualMachineClient = createClient(VirtualMachineService, transport);

	let enhancedDisks: EnhancedDisk[] = $state([]);

	async function loadEnhancedDisks() {
		try {
			// Get data volumes
			const dataVolumesResponse = await virtualMachineClient.listDataVolumes({
				scopeUuid: '', // Add your scope UUID
				facilityName: '', // Add your facility name
				namespace: virtualMachine.namespace,
				bootImage: false, // Set to true if you only want boot images
			});

			// Create a map of data volumes by name for quick lookup
			const dataVolumeMap = new SvelteMap<string, DataVolume>();
			dataVolumesResponse.dataVolumes.forEach((dv) => {
				dataVolumeMap.set(dv.name, dv);
			});

			// Combine disk data with data volume information
			enhancedDisks = virtualMachine.disks.map((disk) => {
				// Only map if volume name exists and source type is DATA_VOLUME
				if (
					disk.volume?.name &&
					disk.volume?.source?.type === VirtualMachine_Disk_Volume_Source_Type.DATA_VOLUME
				) {
					const dataVolumeName = disk.volume.source.data;
					const dataVolume = dataVolumeName ? dataVolumeMap.get(dataVolumeName) : undefined;

					return {
						...disk, // Spread all original disk properties
						bootImage: dataVolume?.bootImage,
						sizeBytes: dataVolume?.sizeBytes,
						phase: dataVolume?.phase,
					};
				} else {
					// Return original disk without DataVolume properties
					return {
						...disk,
						bootImage: undefined,
						sizeBytes: undefined,
						phase: undefined,
					};
				}
			});
		} catch (error) {
			console.error('Failed to load enhanced disks:', error);
		}
	}

	onMount(() => {
		loadEnhancedDisks();
	});
</script>

<div class="flex items-center justify-end gap-1">
	{enhancedDisks.length}
	<Sheet.Root>
		<Sheet.Trigger>
			<Icon icon="ph:arrow-square-out" />
		</Sheet.Trigger>
		<Sheet.Content class="min-w-[70vw] p-4">
			<DataTable {enhancedDisks} />
		</Sheet.Content>
	</Sheet.Root>
</div>
