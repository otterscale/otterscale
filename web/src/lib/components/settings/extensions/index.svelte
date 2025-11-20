<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { type Writable, writable } from 'svelte/store';

	import { type Extension, OrchestratorService } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import * as Accordion from '$lib/components/ui/accordion/index.js';

	import Node from './extension-node.svelte';
	import Thumbnail from './extension-thumbnail.svelte';
	import { getAccordionValue } from './utils.svelte';
</script>

<script lang="ts">
	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const orchestratorClient = createClient(OrchestratorService, transport);

	const instanceExtensions: Writable<Extension[]> = writable([]);
	let isInstanceExtensionsLoaded = $state(false);
	const modelExtensions: Writable<Extension[]> = writable([]);
	let isModelExtensionsLoaded = $state(false);
	const storageExtensions: Writable<Extension[]> = writable([]);
	let isStorageExtensionsLoaded = $state(false);
	const generalExtensions: Writable<Extension[]> = writable([]);
	let isGeneralExtensionsLoaded = $state(false);

	async function fetchInstanceExtensions() {
		try {
			const response = await orchestratorClient.listInstanceExtensions({ scope: scope });
			instanceExtensions.set(response.Extensions);
			isInstanceExtensionsLoaded = true;
		} catch (error) {
			console.error('Failed to fetch instance extensions:', error);
		}
	}

	async function fetchModelExtensions() {
		try {
			const response = await orchestratorClient.listModelExtensions({ scope: scope });
			modelExtensions.set(response.Extensions);
			isModelExtensionsLoaded = true;
		} catch (error) {
			console.error('Failed to fetch model extensions:', error);
		}
	}

	async function fetchStorageExtensions() {
		try {
			const response = await orchestratorClient.listStorageExtensions({ scope: scope });
			storageExtensions.set(response.Extensions);
			isStorageExtensionsLoaded = true;
		} catch (error) {
			console.error('Failed to fetch storage extensions:', error);
		}
	}

	async function fetchGeneralExtensions() {
		try {
			const response = await orchestratorClient.listGeneralExtensions({ scope: scope });
			generalExtensions.set(response.Extensions);
			isGeneralExtensionsLoaded = true;
		} catch (error) {
			console.error('Failed to fetch general extensions:', error);
		}
	}

	onMount(() => {
		Promise.all([
			fetchInstanceExtensions(),
			fetchModelExtensions(),
			fetchStorageExtensions(),
			fetchGeneralExtensions()
		]);
	});
</script>

<Accordion.Root
	type="multiple"
	class="group w-full overflow-hidden rounded-lg border bg-card text-card-foreground transition-all duration-300 **:data-[slot='accordion-trigger']:p-6"
	value={getAccordionValue()}
>
	{#if isInstanceExtensionsLoaded && $instanceExtensions.filter((instanceExtension) => !instanceExtension.latest).length == 0}
		<Accordion.Item value="instance">
			<Accordion.Trigger>
				<Thumbnail
					{scope}
					extensionsBundle="instance"
					extensions={instanceExtensions}
					updator={fetchInstanceExtensions}
				/>
			</Accordion.Trigger>
			<Accordion.Content>
				{#each $instanceExtensions as instanceExtension, index (index)}
					<Node extension={instanceExtension} alignment={index % 2 ? 'right' : 'left'} />
				{/each}
			</Accordion.Content>
		</Accordion.Item>
	{/if}

	{#if isModelExtensionsLoaded && $modelExtensions.filter((modelExtension) => !modelExtension.latest).length == 0}
		<Accordion.Item value="model">
			<Accordion.Trigger>
				<Thumbnail
					{scope}
					extensionsBundle="model"
					extensions={modelExtensions}
					updator={fetchModelExtensions}
				/>
			</Accordion.Trigger>
			<Accordion.Content>
				{#each $modelExtensions as modelExtension, index (index)}
					<Node extension={modelExtension} alignment={index % 2 ? 'right' : 'left'} />
				{/each}
			</Accordion.Content>
		</Accordion.Item>
	{/if}

	{#if isGeneralExtensionsLoaded && isStorageExtensionsLoaded && $storageExtensions.filter((storageExtension) => !storageExtension.latest).length == 0}
		<Accordion.Item value="storage">
			<Accordion.Trigger>
				<Thumbnail
					{scope}
					extensionsBundle="storage"
					extensions={storageExtensions}
					updator={fetchStorageExtensions}
				/>
			</Accordion.Trigger>
			<Accordion.Content>
				{#each $storageExtensions as storageExtension, index (index)}
					<Node extension={storageExtension} alignment={index % 2 ? 'right' : 'left'} />
				{/each}
			</Accordion.Content>
		</Accordion.Item>
	{/if}

	{#if $generalExtensions.filter((generalExtension) => !generalExtension.latest).length == 0}
		<Accordion.Item value="general">
			<Accordion.Trigger>
				<Thumbnail
					{scope}
					extensionsBundle="general"
					extensions={generalExtensions}
					updator={fetchGeneralExtensions}
				/>
			</Accordion.Trigger>
			<Accordion.Content>
				{#each $generalExtensions as generalExtension, index (index)}
					<Node extension={generalExtension} alignment={index % 2 ? 'right' : 'left'} />
				{/each}
			</Accordion.Content>
		</Accordion.Item>
	{/if}
</Accordion.Root>
