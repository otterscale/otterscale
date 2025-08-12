<script lang="ts" module>
	import { StorageService, type Image } from '$lib/api/storage/v1/storage_pb';
	import { DataTable as DataTableLoading } from '$lib/components/custom/loading';
	import * as Reloader from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table';
	import { activeScope } from '$lib/stores';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	const selectedScopeUuid = $activeScope ? $activeScope.uuid : '';

	const images = $state(writable([] as Image[]));
	const reloadManager = new Reloader.ReloadManager(() => {
		storageClient
			.listImages({ scopeUuid: selectedScopeUuid, facilityName: 'ceph-mon' })
			.then((response) => {
				images.set(response.images);
			});
	});

	let isMounted = $state(false);
	onMount(() => {
		storageClient
			.listImages({ scopeUuid: selectedScopeUuid, facilityName: 'ceph-mon' })
			.then((response) => {
				images.set(response.images);
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
	{#if !isMounted}
		<DataTableLoading />
	{:else}
		<Reloader.Root {reloadManager} />
		<DataTable {selectedScopeUuid} {images} />
	{/if}
</main>
