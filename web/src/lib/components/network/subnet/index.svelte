<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';

	import { type Network,NetworkService } from '$lib/api/network/v1/network_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table/index';
	import { Statistics } from './statistics';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');

	const networks = writable<Network[]>([]);
	let isMounted = $state(false);

	const networkClient = createClient(NetworkService, transport);
	const reloadManager = new ReloadManager(() => {
		networkClient.listNetworks({}).then((response) => {
			networks.set(response.networks);
		});
	});
	setContext('reloadManager', reloadManager);

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

<main class="space-y-4 py-4">
	{#if isMounted}
		<Statistics networks={$networks} />
		<DataTable {networks} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
