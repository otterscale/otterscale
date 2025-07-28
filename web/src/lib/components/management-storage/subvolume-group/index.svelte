<script lang="ts" module>
	import { StorageService, type SubvolumeGroup } from '$gen/api/storage/v1/storage_pb';
	import { DataTable as DataTableLoading } from '$lib/components/custom/loading';
	import * as Reloader from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table';
	import Pickers from './pickers.svelte';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	let selectedScope = $state('b62d195e-3905-4960-85ee-7673f71eb21e');
	let selectedFacility = $state('ceph-mon');
	let selectedVolume = $state('');

	const subvolumeGroups = $state(writable([] as SubvolumeGroup[]));
	const reloadManager = new Reloader.ReloadManager(() => {
		storageClient
			.listSubvolumeGroups({
				scopeUuid: selectedScope,
				facilityName: selectedFacility,
				volumeName: selectedVolume
			})
			.then((response) => {
				subvolumeGroups.set(response.subvolumeGroups);
			});
	});

	let isMounted = $state(false);
	onMount(() => {
		storageClient
			.listSubvolumeGroups({
				scopeUuid: selectedScope,
				facilityName: selectedFacility,
				volumeName: selectedVolume
			})
			.then((response) => {
				subvolumeGroups.set(response.subvolumeGroups);
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
	<Pickers bind:selectedScope bind:selectedFacility bind:selectedVolume />
	<Reloader.Root {reloadManager} />
	{#if !isMounted}
		<DataTableLoading />
	{:else}
		<DataTable {selectedScope} {selectedFacility} {selectedVolume} {subvolumeGroups} />
	{/if}
</main>
