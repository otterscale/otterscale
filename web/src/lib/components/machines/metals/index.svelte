<script lang="ts" module>
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table/index';
	import { Statistics } from './statistics';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);

	const machines = writable<Machine[]>([]);

	let isMounted = $state(false);

	const reloadManager = new ReloadManager(() => {
		machineClient.listMachines({}).then((response) => {
			machines.set(response.machines);
		});
	});
	setContext('reloadManager', reloadManager);

	onMount(() => {
		machineClient
			.listMachines({})
			.then((response) => {
				machines.set(response.machines);
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
		<Statistics machines={$machines} />
		<DataTable {machines} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
