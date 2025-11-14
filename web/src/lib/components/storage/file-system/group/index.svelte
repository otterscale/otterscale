<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, setContext, type Snippet } from 'svelte';
	import { writable } from 'svelte/store';

	import { StorageService, type SubvolumeGroup } from '$lib/api/storage/v1/storage_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table';
	import * as store from './utils.svelte';
</script>

<script lang="ts">
	let {
		selectedScope = $bindable(),
		selectedFacility = $bindable(),
		selectedVolume = $bindable(),
		trigger
	}: {
		selectedScope: string;
		selectedFacility: string;
		selectedVolume: string;
		trigger: Snippet;
	} = $props();

	const transport: Transport = getContext('transport');

	let isMounted = $state(false);

	const subvolumeGroups = $state(writable([] as SubvolumeGroup[]));
	const groupStore: store.GroupStore = store.createGroupStore();
	const storageClient = createClient(StorageService, transport);
	const reloadManager = new ReloadManager(() => {
		storageClient
			.listSubvolumeGroups({
				scope: selectedScope,
				facility: selectedFacility,
				volumeName: selectedVolume
			})
			.then((response) => {
				subvolumeGroups.set(response.subvolumeGroups);
			});
	});
	setContext('groupStore', groupStore);

	onMount(() => {
		groupStore.selectedScope.set(selectedScope);
		groupStore.selectedFacility.set(selectedFacility);
		groupStore.selectedVolumeName.set(selectedVolume);

		storageClient
			.listSubvolumeGroups({
				scope: selectedScope,
				facility: selectedFacility,
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

<main class="space-y-4 py-4">
	{#if isMounted}
		{@render trigger()}
		<DataTable {subvolumeGroups} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
