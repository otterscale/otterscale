<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as Alert from '$lib/components/ui/alert';
	import { m } from '$lib/paraglide/messages';
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

<Alert.Root variant="default">
	<Icon icon="ph:airplane-takeoff" />
	<Alert.Title>{m.migrating()}</Alert.Title>
	<Alert.Description>{m.migrating_description()}</Alert.Description>
</Alert.Root>

<div class="pointer-events-none fixed inset-0 flex flex-col items-center justify-center">
	<Icon icon="ph:barricade" class="text-9xl" />
	{m.current_version({ version: import.meta.env.PACKAGE_VERSION })}
</div>
