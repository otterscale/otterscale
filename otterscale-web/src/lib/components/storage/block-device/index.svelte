<script lang="ts">
	import { StorageService, type Image } from '$lib/api/storage/v1/storage_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { Reloader, ReloadManager } from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table';

	let {
		selectedScopeUuid = $bindable(),
		selectedFacility = $bindable()
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
	} = $props();

	const transport: Transport = getContext('transport');

	const images = $state(writable([] as Image[]));
	let isMounted = $state(false);

	const storageClient = createClient(StorageService, transport);
	const reloadManager = new ReloadManager(() => {
		storageClient
			.listImages({ scopeUuid: selectedScopeUuid, facilityName: selectedFacility })
			.then((response) => {
				images.set(response.images);
			});
	});
	setContext('reloadManager', reloadManager);

	onMount(() => {
		storageClient
			.listImages({ scopeUuid: selectedScopeUuid, facilityName: selectedFacility })
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

<main class="space-y-4 py-4">
	{#if isMounted}
		<Reloader {reloadManager} />
		<DataTable {images} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
