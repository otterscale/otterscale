<script lang="ts">
	import { page } from '$app/state';
	import { ManagementMachine } from '$lib/components/otterscale/index';
	import Icon from '@iconify/svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { Nexus, type Machine, type Network } from '$gen/api/nexus/v1/nexus_pb';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const machineStore = writable<Machine>();
	const machineIsLoading = writable(true);
	async function fetchMachine() {
		try {
			const response = await client.getMachine({
				id: page.params.system_id
			});
			machineStore.set(response);
		} catch (error) {
			console.error('Error fetching machine:', error);
		} finally {
			machineIsLoading.set(false);
		}
	}

	let mounted = false;
	onMount(async () => {
		try {
			await fetchMachine();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

<main class="h-[calc(100vh_-_theme(spacing.16))]">
	{#if mounted}
		<ManagementMachine machine={$machineStore} />
	{:else}
		<div class="flex h-full w-full items-center justify-center gap-2 text-sm text-muted-foreground">
			<Icon icon="ph:spinner" class="size-8 animate-spin" />
			Loading...
		</div>
	{/if}
</main>
