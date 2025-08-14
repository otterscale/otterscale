<script lang="ts" module>
	import { StorageService, type SubvolumeGroup } from '$lib/api/storage/v1/storage_pb';
	import { DataTable as DataTableLoading } from '$lib/components/custom/loading';
	import * as Reloader from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, type Snippet } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table';
</script>

<script lang="ts">
	let {
		selectedScopeUuid = $bindable(),
		selectedFacility = $bindable(),
		selectedVolume = $bindable(),
		trigger
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		selectedVolume: string;
		trigger: Snippet;
	} = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	const subvolumeGroups = $state(writable([] as SubvolumeGroup[]));
	const reloadManager = new Reloader.ReloadManager(() => {
		storageClient
			.listSubvolumeGroups({
				scopeUuid: selectedScopeUuid,
				facilityName: selectedFacility,
				volumeName: selectedVolume
			})
			.then((response) => {
				subvolumeGroups.set(response.subvolumeGroups);
			});
	});

	let isMounted = $state(false);
	onMount(() => {
		storageClient
			.listSubvolumeGroups({
				scopeUuid: selectedScopeUuid,
				facilityName: selectedFacility,
				volumeName: selectedVolume
			})
			.then((response) => {
				subvolumeGroups.set(response.subvolumeGroups);
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
			<Reloader.Root {reloadManager} />
		</div>
		<DataTable {selectedScopeUuid} {selectedFacility} {selectedVolume} {subvolumeGroups} />
	{:else}
		<DataTableLoading />
	{/if}
</main>
