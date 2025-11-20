<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { type ObjectStorageDaemon, StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table';
</script>

<script lang="ts">
	let {
		scope
	}: {
		scope: string;
	} = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);
	const reloadManager = new ReloadManager(fetch, false);

	const objectStorageDaemons = writable([] as ObjectStorageDaemon[]);
	async function fetch() {
		try {
			const response = await storageClient.listObjectStorageDaemons({ scope: scope });
			objectStorageDaemons.set(response.objectStorageDaemons);
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	}

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
		<DataTable {objectStorageDaemons} {scope} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
