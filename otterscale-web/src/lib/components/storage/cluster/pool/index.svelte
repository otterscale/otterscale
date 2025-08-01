<script lang="ts" module>
	import { StorageService, type Pool } from '$lib/api/storage/v1/storage_pb';
	import { DataTable as DataTableLoading } from '$lib/components/custom/loading';
	import * as Reloader from '$lib/components/custom/reloader';
	import { activeScope } from '$lib/stores';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	const selectedScopeUuid = $activeScope.uuid;

	const pools = writable([] as Pool[]);
	const reloadManager = new Reloader.ReloadManager(() => {
		storageClient
			.listPools({ scopeUuid: selectedScopeUuid, facilityName: 'ceph-mon' })
			.then((response) => {
				pools.set(response.pools);
			});
	});

	let isMounted = $state(false);
	onMount(() => {
		storageClient
			.listPools({ scopeUuid: selectedScopeUuid, facilityName: 'ceph-mon' })
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

<main class="space-y-4">
	<Reloader.Root {reloadManager} />
	{#if !isMounted}
		<DataTableLoading />
	{:else}
		<DataTable {selectedScopeUuid} {pools} />
	{/if}
</main>
