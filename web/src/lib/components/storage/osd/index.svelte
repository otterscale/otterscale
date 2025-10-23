<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';

	import { DataTable } from './data-table';

	import { StorageService, type OSD } from '$lib/api/storage/v1/storage_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';
</script>

<script lang="ts">
	let {
		selectedScope = $bindable(),
		selectedFacility = $bindable(),
	}: {
		selectedScope: string;
		selectedFacility: string;
	} = $props();

	const transport: Transport = getContext('transport');

	const storageClient = createClient(StorageService, transport);

	const objectStorageDaemons = $state(writable([] as OSD[]));

	const reloadManager = new ReloadManager(() => {
		storageClient.listOSDs({ scope: selectedScope, facility: selectedFacility }).then((response) => {
			objectStorageDaemons.set(response.osds);
		});
	});
	setContext('reloadManager', reloadManager);

	let isMounted = $state(false);
	onMount(() => {
		storageClient
			.listOSDs({ scope: selectedScope, facility: selectedFacility })
			.then((response) => {
				objectStorageDaemons.set(response.osds);
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
		<DataTable {objectStorageDaemons} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
