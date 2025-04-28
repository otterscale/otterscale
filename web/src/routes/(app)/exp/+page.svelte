<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { StackService, type Machine, type Application } from '$gen/api/stack/v1/stack_pb';
	import { getContext, onMount } from 'svelte';
	import Icon from '@iconify/svelte';
	import { get, writable } from 'svelte/store';

	import type { Plugin } from 'svelte-exmarkdown';
	import Markdown from 'svelte-exmarkdown';

	import rehypeHighlight from 'rehype-highlight';

	import yaml from 'highlight.js/lib/languages/yaml';
	import 'highlight.js/styles/github.css';

	const plugins: Plugin[] = [
		{
			rehypePlugin: [rehypeHighlight, { ignoreMissing: true, languages: { yaml } }]
		}
	];

	// Get the transport out of context
	const transport: Transport = getContext('transport');
	const client = createClient(StackService, transport);

	// Create a store for machines
	const machinesStore = writable<Machine[]>([]);
	const applicationStore = writable<Application>();
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

	async function getApplication() {
		try {
			const application = await client.getApplication({
				modelUuid: '5dcb4647-0618-461d-85ca-02660fbd53d4',
				name: 'mysql'
			});
			applicationStore.set(application);
		} catch (error) {
			console.error('Error fetching machines:', error);
		} finally {
			isLoading.set(false);
		}
	}

	onMount(() => {
		fetchMachines();
		getApplication();
	});
</script>

{#if $isLoading}
	<div class="flex h-full w-full items-center justify-center gap-2 text-sm text-muted-foreground">
		<Icon icon="ph:spinner" class="size-8 animate-spin" />
		Loading...
	</div>
{:else}
	<div class="markdown-body">
		{#if $applicationStore?.configYaml}
			<Markdown {plugins} md={'```yaml\n' + $applicationStore.configYaml + '```'} />
		{:else}
			<p class="text-muted-foreground">No configuration available</p>
		{/if}
	</div>
	{#each $machinesStore as machine}
		<p>{machine.fqdn}: {machine.tags}</p>
	{/each}
{/if}
