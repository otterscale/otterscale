<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { type Machine, MachineService } from '$lib/api/machine/v1/machine_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table/index';
	import { Statistics } from './statistics';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);

	const machines = writable<Machine[]>([]);
	async function fetch() {
		machineClient
			.listMachines({})
			.then((response) => {
				machines.set(response.machines);
				isMounted = true;
			})
			.catch((error) => {
				console.error('Error during initial data load:', error);
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
		<Statistics machines={$machines} />
		<DataTable {machines} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
