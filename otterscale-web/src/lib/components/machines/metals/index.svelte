<script lang="ts" module>
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Loading from '$lib/components/custom/loading';
	import * as Reloader from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table/index';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');

	const machines = writable<Machine[]>([]);
	let isMounted = $state(false);

	const machineClient = createClient(MachineService, transport);
	const reloadManager = new Reloader.ReloadManager(() => {
		machineClient.listMachines({}).then((response) => {
			machines.set(response.machines);
		});
	});
	setContext('ReloadManager', reloadManager);

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

<main class="space-y-4">
	{#if isMounted}
		<Reloader.Root {reloadManager} />
		<DataTable {machines} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
