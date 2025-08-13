<script lang="ts" module>
	import { StorageService, type User } from '$lib/api/storage/v1/storage_pb';
	import { DataTable as DataTableLoading } from '$lib/components/custom/loading';
	import * as Reloader from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table';
</script>

<script lang="ts">
	let {
		selectedScopeUuid = $bindable(),
		selectedFacility = $bindable()
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
	} = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	const users = $state(writable([] as User[]));
	const reloadManager = new Reloader.ReloadManager(() => {
		storageClient
			.listUsers({ scopeUuid: selectedScopeUuid, facilityName: selectedFacility })
			.then((response) => {
				users.set(response.users);
			});
	});

	let isMounted = $state(false);
	onMount(() => {
		storageClient
			.listUsers({ scopeUuid: selectedScopeUuid, facilityName: selectedFacility })
			.then((response) => {
				users.set(response.users);
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
		<DataTable {selectedScopeUuid} {selectedFacility} {users} />
	{/if}
</main>
