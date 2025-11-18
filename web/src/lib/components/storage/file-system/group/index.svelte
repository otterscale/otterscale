<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, type Snippet } from 'svelte';
	import { writable } from 'svelte/store';

	import { StorageService, type SubvolumeGroup } from '$lib/api/storage/v1/storage_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table';
</script>

<script lang="ts">
	let {
		scope,
		volume,
		trigger
	}: {
		scope: string;
		volume: string;
		trigger: Snippet;
	} = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	const subvolumeGroups = writable([] as SubvolumeGroup[]);
	async function fetch() {
		storageClient
			.listSubvolumeGroups({
				scope: scope,
				volumeName: volume
			})
			.then((response) => {
				subvolumeGroups.set(response.subvolumeGroups);
			})
			.catch((error) => {
				console.error('Error fetching subvolume groups:', error);
			});
	}
	const reloadManager = new ReloadManager(fetch, false);

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
		{@render trigger()}
		<DataTable {subvolumeGroups} {scope} {volume} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
