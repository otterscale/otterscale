<script lang="ts" module>
	import { StorageService, type Subvolume } from '$lib/api/storage/v1/storage_pb';
	import { DataTable as DataTableLoading } from '$lib/components/custom/loading';
	import * as Reloader from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, type Snippet } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table';
	import Pickers from './pickers.svelte';
</script>

<script lang="ts">
	let {
		selectedScopeUuid = $bindable(),
		selectedFacility = $bindable(),
		selectedVolume = $bindable(),
		selectedSubvolumeGroupName = $bindable(),
		trigger
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		selectedVolume: string;
		selectedSubvolumeGroupName: string;
		trigger: Snippet;
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
				groupName: selectedSubvolumeGroupName
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
				groupName: selectedSubvolumeGroupName
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
	{#if isMounted}
		<div class="flex items-center justify-between gap-2">
			{@render trigger()}
			<div class="flex items-center justify-end gap-2">
				<Pickers
					{selectedScopeUuid}
					{selectedFacility}
					{selectedVolume}
					bind:selectedSubvolumeGroupName
				/>
				<Reloader.Root {reloadManager} />
			</div>
		</div>
		<DataTable
			{selectedScopeUuid}
			{selectedFacility}
			{selectedVolume}
			{selectedSubvolumeGroupName}
			{subvolumes}
		/>
	{:else}
		<DataTableLoading />
	{/if}
</main>
