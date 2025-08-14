<script lang="ts" module>
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Loading from '$lib/components/custom/loading';
	import * as Reloader from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table/index';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);

	const machines = writable<Machine[]>([]);
	const reloadManager = new Reloader.ReloadManager(() => {
		machineClient.listMachines({}).then((response) => {
			machines.set(response.machines);
		});
	});

	let isMounted = $state(false);
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
	{#if !isMounted}
		<Loading.DataTable />
	{:else}
		<Reloader.Root {reloadManager} />
		<DataTable {machines} />
	{/if}
</main>
