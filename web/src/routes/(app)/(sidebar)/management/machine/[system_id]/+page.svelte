<script lang="ts">
	import { page } from '$app/state';
	import { ManagementMachine } from '$lib/components/otterscale/index';
	import Icon from '@iconify/svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { MachineService, type Machine } from '$gen/api/machine/v1/machine_pb';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	const transport: Transport = getContext('transport');
	const client = createClient(MachineService, transport);

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
