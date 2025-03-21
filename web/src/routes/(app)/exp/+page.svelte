<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { StackService, type Machine } from '$gen/api/stack/v1/stack_pb';
	import { getContext, onMount } from 'svelte';
	import Icon from '@iconify/svelte';
	import { writable } from 'svelte/store';

	// Get the transport out of context
	const transport: Transport = getContext('transport');
	const client = createClient(StackService, transport);

	// Create a store for machines
	const machinesStore = writable<Machine[]>([]);
	const isLoading = writable(true);

	// Extract this to a separate function for better organization
	async function fetchMachines() {
		try {
			const response = await client.listMachines({
				pageSize: 10
			});
			machinesStore.set(response.machines);
		} catch (error) {
			console.error('Error fetching machines:', error);
		} finally {
			isLoading.set(false);
		}
	}

	onMount(() => {
		fetchMachines();
	});
</script>

{#if $isLoading}
	<div class="text-muted-foreground flex h-full w-full items-center justify-center gap-2 text-sm">
		<Icon icon="ph:spinner" class="size-8 animate-spin" />
		Loading...
	</div>
{:else}
	{#each $machinesStore as machine}
		<p>{machine.fqdn}: {machine.tags}</p>
	{/each}
{/if}
