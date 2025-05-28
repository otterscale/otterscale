<script lang="ts">
	import { PageLoading } from '$lib/components/otterscale/ui/index';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { ManagementNetworks } from '$lib/components/otterscale';
	import { Nexus, type Machine, type Network } from '$gen/api/nexus/v1/nexus_pb';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	const transport: Transport = getContext('transport');
	const client = createClient(Nexus, transport);

	const networksStore = writable<Network[]>([]);
	const networksLoading = writable(true);
	async function fetchNetworks() {
		try {
			const response = await client.listNetworks({});
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
			const response = await client.listMachines({});
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
	<ManagementNetworks networks={$networksStore} />
{:else}
	<PageLoading />
{/if}
