<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, type Snippet } from 'svelte';
	import { writable } from 'svelte/store';

	import { type Bucket, StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table';
</script>

<script lang="ts">
	let {
		scope,
		trigger
	}: {
		scope: string;
		trigger: Snippet;
	} = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	const buckets = writable([] as Bucket[]);
	let serviceUri = $state('');
	async function fetch() {
		try {
			const response = await storageClient.listBuckets({ scope: scope });
			buckets.set(response.buckets);
			serviceUri = response.serviceUri;
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
		{@render trigger()}
		<DataTable {buckets} {scope} {serviceUri} {reloadManager} />
	{:else}
		<Loading.DataTables.Table />
	{/if}
</main>
