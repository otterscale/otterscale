<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { type Image, StorageService } from '$lib/api/storage/v1/storage_pb';
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

	const images = $state(writable([] as Image[]));
	let isMounted = $state(false);

	const storageClient = createClient(StorageService, transport);
	const reloadManager = new ReloadManager(() => {
		storageClient.listImages({ scope: scope }).then((response) => {
			images.set(response.images);
		});
	});

	onMount(() => {
		storageClient
			.listImages({ scope: scope })
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
		<DataTable {images} {scope} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
