<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, setContext, type Snippet } from 'svelte';
	import { writable } from 'svelte/store';

	import { StorageService, type Subvolume } from '$lib/api/storage/v1/storage_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table';
	import Pickers from './pickers.svelte';
	import { createNFSStore, type NFSStore } from './utils.svelte';
</script>

<script lang="ts">
	let {
		selectedScope = $bindable(),
		selectedFacility = $bindable(),
		selectedVolume = $bindable(),
		selectedSubvolumeGroupName = $bindable(),
		trigger
	}: {
		selectedScope: string;
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
				scope: selectedScope,
				facility: selectedFacility,
				volumeName: selectedVolume,
				groupName: selectedSubvolumeGroupName
			})
			.then((response) => {
				subvolumes.set(response.subvolumes);
			});
	});
	setContext('nfsStore', nfsStore);
	setContext('reloadManager', reloadManager);
	onMount(() => {
		nfsStore.selectedScope.set(selectedScope);
		nfsStore.selectedFacility.set(selectedFacility);
		nfsStore.selectedVolumeName.set(selectedVolume);
		nfsStore.selectedSubvolumeGroupName.set(selectedSubvolumeGroupName);

		storageClient
			.listSubvolumes({
				scope: selectedScope,
				facility: selectedFacility,
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

<main class="space-y-4 py-4">
	{#if isMounted}
		<div class="flex items-center justify-between gap-2">
			{@render trigger()}
			<Pickers
				{selectedScope}
				{selectedFacility}
				{selectedVolume}
				bind:selectedSubvolumeGroupName
			/>
		</div>
		<DataTable {subvolumes} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
