<script lang="ts">
	import { PageLoading } from '$lib/components/otterscale/ui/index';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { ManagementMachines } from '$lib/components/otterscale';
	import { NetworkService, type Network } from '$gen/api/network/v1/network_pb';
	import { MachineService, type Machine } from '$gen/api/machine/v1/machine_pb';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	const transport: Transport = getContext('transport');
	const networkClient = createClient(NetworkService, transport);
	const machineClient = createClient(MachineService, transport);

	const networksStore = writable<Network[]>([]);
	const networksLoading = writable(true);
	async function fetchNetworks() {
		try {
			const response = await networkClient.listNetworks({});
			networksStore.set(response.networks);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			networksLoading.set(false);
		}
	}

	const machinesStore = writable<Machine[]>([]);
	const machinesLoading = writable(true);
	async function fetchMachines() {
		try {
			const response = await machineClient.listMachines({});
			machinesStore.set(response.machines);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			machinesLoading.set(false);
		}
	}

	let mounted = false;
	onMount(async () => {
		try {
			await fetchNetworks();
			await fetchMachines();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

{#if mounted}
	<ManagementMachines machines={$machinesStore} />
{:else}
	<PageLoading />
{/if}
