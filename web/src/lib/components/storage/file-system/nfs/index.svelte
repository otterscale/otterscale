<script lang="ts" module>
	import { StorageService, type Subvolume } from '$lib/api/storage/v1/storage_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager, Reloader } from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, setContext, type Snippet } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table';
	import Pickers from './pickers.svelte';
	import { createNFSStore, type NFSStore } from './utils.svelte';
</script>

<script lang="ts">
	let {
		selectedScopeUuid = $bindable(),
		selectedFacility = $bindable(),
		selectedVolume = $bindable(),
		selectedSubvolumeGroupName = $bindable(),
		trigger,
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		selectedVolume: string;
		selectedSubvolumeGroupName: string;
		trigger: Snippet;
	} = $props();

	const transport: Transport = getContext('transport');

	let isMounted = $state(false);
	const subvolumes = $state(writable([] as Subvolume[]));
	const nfsStore: NFSStore = createNFSStore();
	const storageClient = createClient(StorageService, transport);
	const reloadManager = new ReloadManager(() => {
		storageClient
			.listSubvolumes({
				scopeUuid: selectedScopeUuid,
				facilityName: selectedFacility,
				volumeName: selectedVolume,
				groupName: selectedSubvolumeGroupName,
			})
			.then((response) => {
				subvolumes.set(response.subvolumes);
			});
	});
	setContext('nfsStore', nfsStore);
	setContext('reloadManager', reloadManager);
	onMount(() => {
		nfsStore.selectedScopeUuid.set(selectedScopeUuid);
		nfsStore.selectedFacilityName.set(selectedFacility);
		nfsStore.selectedVolumeName.set(selectedVolume);
		nfsStore.selectedSubvolumeGroupName.set(selectedSubvolumeGroupName);

		storageClient
			.listSubvolumes({
				scopeUuid: selectedScopeUuid,
				facilityName: selectedFacility,
				volumeName: selectedVolume,
				groupName: selectedSubvolumeGroupName,
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

<main class="space-y-4 py-4">
	{#if isMounted}
		<div class="flex items-center justify-between gap-2">
			{@render trigger()}
			<Pickers {selectedScopeUuid} {selectedFacility} {selectedVolume} bind:selectedSubvolumeGroupName />
		</div>
		<DataTable {subvolumes} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
