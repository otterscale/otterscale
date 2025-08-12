<script lang="ts" module>
	import { StorageService, type Subvolume } from '$lib/api/storage/v1/storage_pb';
	import { DataTable as DataTableLoading } from '$lib/components/custom/loading';
	import * as Reloader from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table';
	import Pickers from './pickers.svelte';
</script>

<script lang="ts">
	let {
		selectedScopeUuid = $bindable(),
		selectedFacility = $bindable(),
		selectedVolume = $bindable(),
		selectedSubvolumeGroup = $bindable()
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		selectedVolume: string;
		selectedSubvolumeGroup: string;
	} = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	const subvolumes = $state(writable([] as Subvolume[]));
	const reloadManager = new Reloader.ReloadManager(() => {
		storageClient
			.listSubvolumes({
				scopeUuid: selectedScopeUuid,
				facilityName: selectedFacility,
				volumeName: selectedVolume,
				groupName: selectedSubvolumeGroup
			})
			.then((response) => {
				subvolumes.set(response.subvolumes);
			});
	});

	let isMounted = $state(false);
	onMount(() => {
		storageClient
			.listSubvolumes({
				scopeUuid: selectedScopeUuid,
				facilityName: selectedFacility,
				volumeName: selectedVolume,
				groupName: selectedSubvolumeGroup
			})
			.then((response) => {
				subvolumes.set(response.subvolumes);
				isMounted = true;
			})
			.catch((error) => {
				console.error('Error during initial data load:', error);
			});

		reloadManager.start();
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<main class="space-y-4">
	{#if !isMounted}
		<DataTableLoading />
	{:else}
		<Pickers {selectedScopeUuid} {selectedFacility} {selectedVolume} bind:selectedSubvolumeGroup />
		<Reloader.Root {reloadManager} />
		<DataTable
			{selectedScopeUuid}
			{selectedFacility}
			{selectedVolume}
			{selectedSubvolumeGroup}
			{subvolumes}
		/>
	{/if}
</main>
