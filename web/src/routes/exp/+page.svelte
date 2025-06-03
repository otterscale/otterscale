<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import Icon from '@iconify/svelte';
	import { get, writable } from 'svelte/store';

	import type { Plugin } from 'svelte-exmarkdown';
	import Markdown from 'svelte-exmarkdown';

	import rehypeHighlight from 'rehype-highlight';

	import yaml from 'highlight.js/lib/languages/yaml';
	import 'highlight.js/styles/github.css';
	import { MachineService, type Machine } from '$gen/api/machine/v1/machine_pb';

	const plugins: Plugin[] = [
		{
			rehypePlugin: [rehypeHighlight, { ignoreMissing: true, languages: { yaml } }]
		}
	];

	// Get the transport out of context
	const transport: Transport = getContext('transport');
	const client = createClient(MachineService, transport);

	// Create a store for machines
	const machinesStore = writable<Machine[]>([]);
	const isLoading = writable(true);

	// Extract this to a separate function for better organization
	async function fetchMachines() {
		try {
			const response = await client.listMachines({});
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
