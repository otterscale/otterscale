<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { writable, type Writable } from 'svelte/store';

	import { getAccordionValue } from './utils.svelte';

	import { OrchestratorService, type Plugin } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import Node from '$lib/components/settings/plugins/plugin-node.svelte';
	import Thumbnail from '$lib/components/settings/plugins/plugin-thumbnail.svelte';
	import * as Accordion from '$lib/components/ui/accordion/index.js';
</script>

<script lang="ts">
	let { scope, facility }: { scope: string; facility: string } = $props();

	const transport: Transport = getContext('transport');
	const orchestratorClient = createClient(OrchestratorService, transport);

	const modelPlugins: Writable<Plugin[]> = writable([]);
	const generalPlugins: Writable<Plugin[]> = writable([]);

	orchestratorClient
		.listModelPlugins({ scope: scope, facility: facility })
		.then((respoonse) => {
			modelPlugins.set(respoonse.plugins);
		})
		.catch((error) => {
			console.error('Failed to fetch model plugins:', error);
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
	class="group bg-card text-card-foreground w-full overflow-hidden rounded-lg border transition-all duration-300"
	value={getAccordionValue()}
>
	<Accordion.Item value="model" class="p-6">
		<Accordion.Trigger>
			<Thumbnail {scope} {facility} pluginsBundle="model" plugins={$modelPlugins} />
		</Accordion.Trigger>
		<Accordion.Content class="flex flex-col gap-4 text-balance">
			{#each $modelPlugins as modelPlugin, index}
				<div class="relative mx-auto hidden w-full space-y-12 pt-10 md:block">
					<Node {scope} {facility} plugin={modelPlugin} alignment={index % 2 ? 'right' : 'left'} />
				</div>
			{/each}
		</Accordion.Content>
	</Accordion.Item>

	<Accordion.Item value="general" class="p-6">
		<Accordion.Trigger>
			<Thumbnail {scope} {facility} pluginsBundle="general" plugins={$generalPlugins} />
		</Accordion.Trigger>
		<Accordion.Content class="flex flex-col gap-4 text-balance">
			{#each $generalPlugins as generalPlugin, index}
				<div class="relative mx-auto hidden w-full space-y-12 pt-10 md:block">
					<Node {scope} {facility} plugin={generalPlugin} alignment={index % 2 ? 'right' : 'left'} />
				</div>
			{/each}
		</Accordion.Content>
	</Accordion.Item>
</Accordion.Root>
