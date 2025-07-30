<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { homePath, machinesPath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [homePath], current: machinesPath });

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);

	const machinesStore = writable<Machine[]>([]);
	async function fetchMachines() {
		try {
			const response = await machineClient.listMachines({});
			machinesStore.set(response.machines);
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	let mounted = false;
	onMount(async () => {
		try {
			await fetchMachines();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
		mounted = true;
	});
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}

{#if mounted}
	{#each $machinesStore as machine}
		<p>{machine.id}</p>
	{/each}
{/if}
