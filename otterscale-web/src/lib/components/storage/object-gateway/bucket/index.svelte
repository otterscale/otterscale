<script lang="ts" module>
	import { StorageService, type Bucket } from '$lib/api/storage/v1/storage_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { Reloader, ReloadManager } from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, setContext, type Snippet } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table';
</script>

<script lang="ts">
	let {
		selectedScopeUuid = $bindable(),
		selectedFacility = $bindable(),
		trigger
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		trigger: Snippet;
	} = $props();

	const transport: Transport = getContext('transport');

	const buckets = $state(writable([] as Bucket[]));
	let isMounted = $state(false);

	const storageClient = createClient(StorageService, transport);
	const reloadManager = new ReloadManager(() => {
		storageClient
			.listBuckets({ scopeUuid: selectedScopeUuid, facilityName: selectedFacility })
			.then((response) => {
				buckets.set(response.buckets);
			});
	});
	setContext('reloadManager', reloadManager);

	onMount(() => {
		storageClient
			.listBuckets({ scopeUuid: selectedScopeUuid, facilityName: selectedFacility })
			.then((response) => {
				buckets.set(response.buckets);
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
			<Reloader {reloadManager} />
		</div>
		<DataTable {buckets} />
	{:else}
		<Loading.DataTables.Table />
	{/if}
</main>
