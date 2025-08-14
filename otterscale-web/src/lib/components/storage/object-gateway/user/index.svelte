<script lang="ts" module>
	import { StorageService, type User } from '$lib/api/storage/v1/storage_pb';
	import * as Loading from '$lib/components/custom/loading';
	import * as Reloader from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, type Snippet } from 'svelte';
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
	{#if isMounted}
		<div class="flex items-center justify-between gap-2">
			{@render trigger()}
			<Reloader.Root {reloadManager} />
		</div>
		<DataTable {selectedScopeUuid} {selectedFacility} {users} />
	{:else}
		<Loading.DataTables.Table />
	{/if}
</main>
