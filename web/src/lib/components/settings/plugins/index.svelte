<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { writable, type Writable } from 'svelte/store';

	import Node from './plugin-node.svelte';
	import Thumbnail from './plugin-thumbnail.svelte';
	import { getAccordionValue } from './utils.svelte';

	import { OrchestratorService, type Plugin } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import * as Accordion from '$lib/components/ui/accordion/index.js';
</script>

<script lang="ts">
	let { scope, facility }: { scope: string; facility: string } = $props();

	const transport: Transport = getContext('transport');
	const orchestratorClient = createClient(OrchestratorService, transport);

	const instancePlugins: Writable<Plugin[]> = writable([]);
	const modelPlugins: Writable<Plugin[]> = writable([]);
	const storagePlugins: Writable<Plugin[]> = writable([]);
	const generalPlugins: Writable<Plugin[]> = writable([]);

	orchestratorClient
		.listInstancePlugins({ scope: scope, facility: facility })
		.then((respoonse) => {
			instancePlugins.set(respoonse.plugins);
		})
		.catch((error) => {
			console.error('Failed to fetch instance plugins:', error);
		});
	orchestratorClient
		.listModelPlugins({ scope: scope, facility: facility })
		.then((respoonse) => {
			modelPlugins.set(respoonse.plugins);
		})
		.catch((error) => {
			console.error('Failed to fetch model plugins:', error);
		});
	orchestratorClient
		.listStoragePlugins({ scope: scope, facility: facility })
		.then((respoonse) => {
			storagePlugins.set(respoonse.plugins);
		})
		.catch((error) => {
			console.error('Failed to fetch storage plugins:', error);
		});
	orchestratorClient
		.listGeneralPlugins({ scope: scope, facility: facility })
		.then((respoonse) => {
			generalPlugins.set(respoonse.plugins);
		})
		.catch((error) => {
			console.error('Failed to fetch general plugins:', error);
		});
</script>

<Accordion.Root
	type="multiple"
	class="group bg-card text-card-foreground w-full overflow-hidden rounded-lg border transition-all duration-300 **:data-[slot='accordion-trigger']:p-6"
	value={getAccordionValue()}
>
	{#if $instancePlugins.filter((instancePlugin) => !instancePlugin.latest).length == 0}
		<Accordion.Item value="instance">
			<Accordion.Trigger>
				<Thumbnail {scope} {facility} pluginsBundle="instance" plugins={instancePlugins} />
			</Accordion.Trigger>
			<Accordion.Content>
				{#each $instancePlugins as instancePlugin, index}
					<Node {scope} {facility} plugin={instancePlugin} alignment={index % 2 ? 'right' : 'left'} />
				{/each}
			</Accordion.Content>
		</Accordion.Item>
	{/if}

	{#if $modelPlugins.filter((modelPlugin) => !modelPlugin.latest).length == 0}
		<Accordion.Item value="model">
			<Accordion.Trigger>
				<Thumbnail {scope} {facility} pluginsBundle="model" plugins={modelPlugins} />
			</Accordion.Trigger>
			<Accordion.Content>
				{#each $modelPlugins as modelPlugin, index}
					<Node {scope} {facility} plugin={modelPlugin} alignment={index % 2 ? 'right' : 'left'} />
				{/each}
			</Accordion.Content>
		</Accordion.Item>
	{/if}

	{#if $storagePlugins.filter((storagePlugin) => !storagePlugin.latest).length == 0}
		<Accordion.Item value="storage">
			<Accordion.Trigger>
				<Thumbnail {scope} {facility} pluginsBundle="storage" plugins={storagePlugins} />
			</Accordion.Trigger>
			<Accordion.Content>
				{#each $storagePlugins as storagePlugin, index}
					<Node {scope} {facility} plugin={storagePlugin} alignment={index % 2 ? 'right' : 'left'} />
				{/each}
			</Accordion.Content>
		</Accordion.Item>
	{/if}

	{#if $generalPlugins.filter((generalPlugin) => !generalPlugin.latest).length == 0}
		<Accordion.Item value="general">
			<Accordion.Trigger>
				<Thumbnail {scope} {facility} pluginsBundle="general" plugins={generalPlugins} />
			</Accordion.Trigger>
			<Accordion.Content>
				{#each $generalPlugins as generalPlugin, index}
					<Node {scope} {facility} plugin={generalPlugin} alignment={index % 2 ? 'right' : 'left'} />
				{/each}
			</Accordion.Content>
		</Accordion.Item>
	{/if}
</Accordion.Root>
