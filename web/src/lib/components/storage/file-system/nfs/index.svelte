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
	const storageClient = createClient(StorageService, transport);

	const subvolumes = writable([] as Subvolume[]);
	async function fetch() {
		try {
			const response = await storageClient.listSubvolumes({
				scope: scope,
				volumeName: volume,
				groupName: selectedSubvolumeGroupName
			});
			subvolumes.set(response.subvolumes);
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	}
	const reloadManager = new ReloadManager(fetch);

	let isMounted = $state(false);
	onMount(async () => {
		await fetch();
		isMounted = true;
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
