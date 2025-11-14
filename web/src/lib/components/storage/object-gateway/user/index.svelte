<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, type Snippet } from 'svelte';
	import { writable } from 'svelte/store';

	import { StorageService, type User } from '$lib/api/storage/v1/storage_pb';
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

	const users = $state(writable([] as User[]));
	let isMounted = $state(false);

	const storageClient = createClient(StorageService, transport);
	const reloadManager = new ReloadManager(() => {
		storageClient.listUsers({ scope: scope }).then((response) => {
			users.set(response.users);
		});
	});

	onMount(() => {
		storageClient
			.listUsers({ scope: scope })
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

<main class="space-y-4 py-4">
	{#if isMounted}
		{@render trigger()}
		<DataTable {users} {scope} {reloadManager} />
	{:else}
		<Loading.DataTables.Table />
	{/if}
</main>
