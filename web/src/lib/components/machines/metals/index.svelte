<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { on } from 'events';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { type Machine, MachineService } from '$lib/api/machine/v1/machine_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table/index';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);

	const machines = writable<Machine[]>([]);
	// async function fetch() {
	// 	machineClient
	// 		.listMachines({})
	// 		.then((response) => {
	// 			machines.set(response.machines);
	// 			// isMachinesLoaded = true;
	// 		})
	// 		.catch((error) => {
	// 			console.error('Error during initial data load:', error);
	// 		});
	// }
	async function fetch() {
		const response = await machineClient.listMachines({});
		machines.set(response.machines);
	}

	// async function fetch2() {
	// 	const response = await machineClient.listMachines({});
	// 	machines.set(response.machines);
	// }
	let isMounted = $state(false);

	const reloadManager = new ReloadManager(fetch, false);

	// const isMounted = $derived(isMachinesLoaded);
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
		<DataTable {machines} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
