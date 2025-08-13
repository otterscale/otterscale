<script lang="ts" module>
	import { NetworkService, type Network } from '$lib/api/network/v1/network_pb';
	import { DataTable as DataTableLoading } from '$lib/components/custom/loading';
	import * as Reloader from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table/index';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const networkClient = createClient(NetworkService, transport);

	const networks = writable<Network[]>([]);
	const reloadManager = new Reloader.ReloadManager(() => {
		networkClient.listNetworks({}).then((response) => {
			networks.set(response.networks);
		});
	});

	let isMounted = $state(false);
	onMount(() => {
		networkClient
			.listNetworks({})
			.then((response) => {
				networks.set(response.networks);
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
		<DataTable {networks} />
	{/if}
</main>
