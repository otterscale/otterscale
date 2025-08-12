<script lang="ts" module>
	import { StorageService, type OSD } from '$lib/api/storage/v1/storage_pb';
	import { DataTable as DataTableLoading } from '$lib/components/custom/loading';
	import * as Reloader from '$lib/components/custom/reloader';
	import { activeScope } from '$lib/stores';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table';

	const selectedFacility = 'ceph-mon';
</script>

<script lang="ts">
	const selectedScopeUuid = $activeScope ? $activeScope.uuid : '';

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	const objectStorageDaemons = $state(writable([] as OSD[]));
	const reloadManager = new Reloader.ReloadManager(() => {
		storageClient
			.listOSDs({ scopeUuid: selectedScopeUuid, facilityName: selectedFacility })
			.then((response) => {
				objectStorageDaemons.set(response.osds);
			});
	});

	let isMounted = $state(false);
	onMount(() => {
		storageClient
			.listOSDs({ scopeUuid: selectedScopeUuid, facilityName: selectedFacility })
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

<main class="space-y-4">
	<Reloader.Root {reloadManager} />
	{#if !isMounted}
		<DataTableLoading />
	{:else}
		<DataTable {selectedScopeUuid} {selectedFacility} {objectStorageDaemons} />
	{/if}
</main>
