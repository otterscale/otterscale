<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { type SMBShare, StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table';
</script>

<script lang="ts">
	let {
		scope,
		namespace
	}: {
		scope: string;
		namespace: string;
	} = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	const smbShares = writable([] as SMBShare[]);
	async function fetch() {
		storageClient
			.listSMBShares({ scope: scope })
			.then((response) => {
				smbShares.set(response.smbShares);
			})
			.catch((error) => {
				console.error('Error reloading SMB shares:', error);
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
		<DataTable {smbShares} {scope} {namespace} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
