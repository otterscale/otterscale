<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, type Snippet } from 'svelte';
	import { writable } from 'svelte/store';

	import { StorageService, type Subvolume } from '$lib/api/storage/v1/storage_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table';
	import Pickers from './pickers.svelte';
</script>

<script lang="ts">
	let {
		scope,
		volume,
		selectedSubvolumeGroupName = $bindable(),
		trigger
	}: {
		scope: string;
		volume: string;
		selectedSubvolumeGroupName: string;
		trigger: Snippet;
	} = $props();

	const transport: Transport = getContext('transport');

	let isMounted = $state(false);
	const subvolumes = $state(writable([] as Subvolume[]));
	const storageClient = createClient(StorageService, transport);
	const reloadManager = new ReloadManager(() => {
		storageClient
			.listSubvolumes({
				scope: scope,
				volumeName: volume,
				groupName: selectedSubvolumeGroupName
			})
			.then((response) => {
				subvolumes.set(response.subvolumes);
			});
	});

	onMount(() => {
		storageClient
			.listSubvolumes({
				scope: scope,
				volumeName: volume,
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
			<Pickers {scope} {volume} bind:selectedSubvolumeGroupName />
		</div>
		<DataTable {subvolumes} {scope} {volume} group={selectedSubvolumeGroupName} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
