<script lang="ts" module>
	import { StorageService, type Pool } from '$lib/api/storage/v1/storage_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table';
</script>

<script lang="ts">
	let {
		selectedScopeUuid = $bindable(),
		selectedFacility = $bindable(),
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
	} = $props();

	const transport: Transport = getContext('transport');

	const storageClient = createClient(StorageService, transport);

	const pools = writable([] as Pool[]);

	const reloadManager = new ReloadManager(() => {
		storageClient.listPools({ scopeUuid: selectedScopeUuid, facilityName: selectedFacility }).then((response) => {
			pools.set(response.pools);
		});
	});
	setContext('reloadManager', reloadManager);

	let isMounted = $state(false);
	onMount(() => {
		storageClient
			.listPools({ scopeUuid: selectedScopeUuid, facilityName: selectedFacility })
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
		<DataTable {pools} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
