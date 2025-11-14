<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onDestroy } from 'svelte';
	import { SvelteMap } from 'svelte/reactivity';

	import type { VirtualMachine } from '$lib/api/instance/v1/instance_pb';
	import type { DataVolume } from '$lib/api/instance/v1/instance_pb';
	import {
		InstanceService,
		VirtualMachine_Disk_Volume_Source_Type
	} from '$lib/api/instance/v1/instance_pb';
	import type { EnhancedDisk } from '$lib/components/compute/virtual-machine/units/type';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Sheet from '$lib/components/ui/sheet';
	import { currentKubernetes } from '$lib/stores';

	import { DataTable } from './data-table';
</script>

<script lang="ts">
	let {
		virtualMachine
	}: {
		virtualMachine: VirtualMachine;
	} = $props();

	// Context dependencies
	const transport: Transport = getContext('transport');
	const virtualMachineClient = createClient(InstanceService, transport);

	let enhancedDisks: EnhancedDisk[] = $state([]);
	let isLoaded = $state(false);

	async function loadEnhancedDisks() {
		try {
			// Get data volumes
			const dataVolumesResponse = await virtualMachineClient.listDataVolumes({
				scope: scope,
				,
				namespace: virtualMachine.namespace,
				bootImage: false // Set to true if you only want boot images
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
						vmName: virtualMachine.name,
						namespace: virtualMachine.namespace
					};
				} else {
					// Return original disk without DataVolume properties
					return {
						...disk,
						bootImage: undefined,
						sizeBytes: undefined,
						phase: undefined,
						vmName: virtualMachine.name,
						namespace: virtualMachine.namespace
					};
				}
			});
			isLoaded = true;
		} catch (error) {
			console.error('Failed to load enhanced disks:', error);
		}
	}

	// Create ReloadManager for automatic reloading (only when data is loaded)
	const reloadManager = new ReloadManager(() => {
		if (isLoaded) {
			loadEnhancedDisks();
		}
	});

	function handleSheetOpen() {
		if (!isLoaded) {
			loadEnhancedDisks();
			reloadManager.start();
		}
	}

	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<div class="flex items-center justify-end gap-1">
	{virtualMachine.disks.length}
	<Sheet.Root>
		<Sheet.Trigger onclick={handleSheetOpen}>
			<Icon icon="ph:arrow-square-out" />
		</Sheet.Trigger>
		<Sheet.Content class="min-w-[70vw] p-4">
			<DataTable {virtualMachine} {enhancedDisks} {reloadManager} />
		</Sheet.Content>
	</Sheet.Root>
</div>
