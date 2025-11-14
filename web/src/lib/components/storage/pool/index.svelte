<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { type Pool, StorageService } from '$lib/api/storage/v1/storage_pb';
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

	const pools = writable([] as Pool[]);

	const reloadManager = new ReloadManager(() => {
		storageClient.listPools({ scope: scope }).then((response) => {
			pools.set(response.pools);
		});
	});

	let isMounted = $state(false);
	onMount(() => {
		storageClient
			.listPools({ scope: scope })
			.then((response) => {
				pools.set(response.pools);
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
		<DataTable {pools} {scope} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
