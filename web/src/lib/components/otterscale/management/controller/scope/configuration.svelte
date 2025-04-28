<!-- <script lang="ts">
	import Markdown from 'svelte-exmarkdown';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { Nexus, type Configuration } from '$gen/api/nexus/v1/nexus_pb';
	import type { Plugin } from 'svelte-exmarkdown';
	import rehypeHighlight from 'rehype-highlight';
	import yaml from 'highlight.js/lib/languages/yaml';
	import 'highlight.js/styles/github.css';
	import Icon from '@iconify/svelte';
	import {
		Drawer,
		DrawerContent,
		DrawerDescription,
		DrawerHeader,
		DrawerTitle
	} from '$lib/components/ui/drawer';

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);
	const plugins: Plugin[] = [
		{
			rehypePlugin: [rehypeHighlight, { ignoreMissing: true, languages: { yaml } }]
		}
	];

	const configurationStore = writable<string>();
	const configurationIsLoading = writable(true);
	async function fetchConfiguration() {
		try {
			const response = await client.getModelConfig({
				uuid: model_uuid
			});
			configurationStore.set(response.configYaml);
		} catch (error) {
			console.error('Error fetching configuration:', error);
		} finally {
			configurationIsLoading.set(false);
		}
	}

	let {
		model_uuid,
		open = $bindable()
	}: {
		model_uuid: string;
		open: boolean;
	} = $props();

	let mounted = $state(false);
	onMount(async () => {
		try {
			await fetchConfiguration();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

{#if mounted}
	<Drawer {open} onOpenChange={(value) => (open = value)}>
		<DrawerContent class="inset-x-auto inset-y-0 right-0 w-3/5 px-3">
			<DrawerHeader>
				<DrawerTitle class="text-center">Configuration</DrawerTitle>
			</DrawerHeader>
			<div class="markdown-body overflow-y-auto px-3">
				{#if $configurationStore.length > 0}
					<Markdown {plugins} md={'```yaml\n' + $configurationStore + '```'} />
				{:else}
					<p class="text-muted-foreground">No configuration available</p>
				{/if}
			</div>
		</DrawerContent>
	</Drawer>
{:else}
	<div class="flex h-full w-full items-center justify-center gap-2 text-sm text-muted-foreground">
		<Icon icon="ph:spinner" class="size-8 animate-spin" />
		Loading...
	</div>
{/if} -->
