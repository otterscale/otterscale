<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { type Network, NetworkService } from '$lib/api/network/v1/network_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table/index';
	import { Statistics } from './statistics';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');

	const networks = writable<Network[]>([]);
	let isLoaded = $state(false);

	const networkClient = createClient(NetworkService, transport);

	async function fetch() {
		try {
			const response = await networkClient.listNetworks({});
			networks.set(response.networks);
		} catch (error) {
			console.error('Failed to fetch networks:', error);
		}
	}

	const reloadManager = new ReloadManager(fetch);

	onMount(async () => {
		await fetch();
		isLoaded = true;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<main class="space-y-4 py-4">
	{#if isLoaded}
		<Statistics networks={$networks} />
		<DataTable {networks} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
