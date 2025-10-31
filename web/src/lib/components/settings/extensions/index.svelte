<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable, type Writable } from 'svelte/store';

	import Node from './extension-node.svelte';
	import Thumbnail from './extension-thumbnail.svelte';
	import { getAccordionValue } from './utils.svelte';

	import { OrchestratorService, type Extension } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import * as Accordion from '$lib/components/ui/accordion/index.js';
</script>

<script lang="ts">
	let { scope, facility }: { scope: string; facility: string } = $props();

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

	onMount(async () => {
		orchestratorClient
			.listInstanceExtensions({ scope: scope, facility: facility })
			.then((respoonse) => {
				instanceExtensions.set(respoonse.Extensions);
				isInstanceExtensionsLoaded = true;
			})
			.catch((error) => {
				console.error('Failed to fetch instance extensions:', error);
			});
		orchestratorClient
			.listModelExtensions({ scope: scope, facility: facility })
			.then((respoonse) => {
				modelExtensions.set(respoonse.Extensions);
				isModelExtensionsLoaded = true;
			})
			.catch((error) => {
				console.error('Failed to fetch model extensions:', error);
			});
		orchestratorClient
			.listStorageExtensions({ scope: scope, facility: facility })
			.then((respoonse) => {
				storageExtensions.set(respoonse.Extensions);
				isStorageExtensionsLoaded = true;
			})
			.catch((error) => {
				console.error('Failed to fetch storage extensions:', error);
			});
		orchestratorClient
			.listGeneralExtensions({ scope: scope, facility: facility })
			.then((respoonse) => {
				generalExtensions.set(respoonse.Extensions);
				isGeneralExtensionsLoaded = true;
			})
			.catch((error) => {
				console.error('Failed to fetch general extensions:', error);
			});
	});
</script>

<Accordion.Root
	type="multiple"
	class="group bg-card text-card-foreground w-full overflow-hidden rounded-lg border transition-all duration-300 **:data-[slot='accordion-trigger']:p-6"
	value={getAccordionValue()}
>
	{#if isInstanceExtensionsLoaded && $instanceExtensions.filter((instanceExtension) => !instanceExtension.latest).length == 0}
		<Accordion.Item value="instance">
			<Accordion.Trigger>
				<Thumbnail
					{scope}
					{facility}
					extensionsBundle="instance"
					extensions={instanceExtensions}
					updator={() => {
						orchestratorClient
							.listInstanceExtensions({ scope: scope, facility: facility })
							.then((respoonse) => {
								instanceExtensions.set(respoonse.Extensions);
							})
							.catch((error) => {
								console.error('Failed to fetch instance extensions:', error);
							});
					}}
				/>
			</Accordion.Trigger>
			<Accordion.Content>
				{#each $instanceExtensions as instanceExtension, index}
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
					{facility}
					extensionsBundle="model"
					extensions={modelExtensions}
					updator={() => {
						orchestratorClient
							.listModelExtensions({ scope: scope, facility: facility })
							.then((respoonse) => {
								modelExtensions.set(respoonse.Extensions);
							})
							.catch((error) => {
								console.error('Failed to fetch model extensions:', error);
							});
					}}
				/>
			</Accordion.Trigger>
			<Accordion.Content>
				{#each $modelExtensions as modelExtension, index}
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
					{facility}
					extensionsBundle="storage"
					extensions={storageExtensions}
					updator={() => {
						orchestratorClient
							.listStorageExtensions({ scope: scope, facility: facility })
							.then((respoonse) => {
								storageExtensions.set(respoonse.Extensions);
							})
							.catch((error) => {
								console.error('Failed to fetch storage extensions:', error);
							});
					}}
				/>
			</Accordion.Trigger>
			<Accordion.Content>
				{#each $storageExtensions as storageExtension, index}
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
					{facility}
					extensionsBundle="general"
					extensions={generalExtensions}
					updator={() => {
						orchestratorClient
							.listGeneralExtensions({ scope: scope, facility: facility })
							.then((respoonse) => {
								generalExtensions.set(respoonse.Extensions);
							})
							.catch((error) => {
								console.error('Failed to fetch general extensions:', error);
							});
					}}
				/>
			</Accordion.Trigger>
			<Accordion.Content>
				{#each $generalExtensions as generalExtension, index}
					<Node extension={generalExtension} alignment={index % 2 ? 'right' : 'left'} />
				{/each}
			</Accordion.Content>
		</Accordion.Item>
	{/if}
</Accordion.Root>
