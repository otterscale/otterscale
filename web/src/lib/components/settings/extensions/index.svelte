<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { type Writable, writable } from 'svelte/store';

	import {
		type Extension,
		Extension_Type,
		OrchestratorService
	} from '$lib/api/orchestrator/v1/orchestrator_pb';
	import * as Accordion from '$lib/components/ui/accordion/index.js';

	import Node from './extension-node.svelte';
	import Thumbnail from './extension-thumbnail.svelte';
	import { getAccordionValue } from './utils.svelte';
</script>

<script lang="ts">
	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const orchestratorClient = createClient(OrchestratorService, transport);

	const metricsExtensions: Writable<Extension[]> = writable([]);
	let isMetricsExtensionsLoaded = $state(false);
	const serviceMeshExtensions: Writable<Extension[]> = writable([]);
	let isServiceMeshExtensionsLoaded = $state(false);
	const modelExtensions: Writable<Extension[]> = writable([]);
	let isModelExtensionsLoaded = $state(false);
	const registryExtensions: Writable<Extension[]> = writable([]);
	let isRegistryExtensionsLoaded = $state(false);
	const instanceExtensions: Writable<Extension[]> = writable([]);
	let isInstanceExtensionsLoaded = $state(false);
	const storageExtensions: Writable<Extension[]> = writable([]);
	let isStorageExtensionsLoaded = $state(false);

	async function fetchMetricsExtensions() {
		try {
			const response = await orchestratorClient.listExtensions({
				scope: scope,
				type: Extension_Type.METRICS
			});
			metricsExtensions.set(response.Extensions);
			isMetricsExtensionsLoaded = true;
		} catch (error) {
			console.error('Failed to fetch metrics extensions:', error);
		}
	}

	async function fetchServiceMeshExtensions() {
		try {
			const response = await orchestratorClient.listExtensions({
				scope: scope,
				type: Extension_Type.SERVICE_MESH
			});
			serviceMeshExtensions.set(response.Extensions);
			isServiceMeshExtensionsLoaded = true;
		} catch (error) {
			console.error('Failed to fetch service mesh extensions:', error);
		}
	}

	async function fetchRegistryExtensions() {
		try {
			const response = await orchestratorClient.listExtensions({
				scope: scope,
				type: Extension_Type.REGISTRY
			});
			registryExtensions.set(response.Extensions);
			isRegistryExtensionsLoaded = true;
		} catch (error) {
			console.error('Failed to fetch registry extensions:', error);
		}
	}

	async function fetchModelExtensions() {
		try {
			const response = await orchestratorClient.listExtensions({
				scope: scope,
				type: Extension_Type.MODEL
			});
			modelExtensions.set(response.Extensions);
			isModelExtensionsLoaded = true;
		} catch (error) {
			console.error('Failed to fetch model extensions:', error);
		}
	}

	async function fetchInstanceExtensions() {
		try {
			const response = await orchestratorClient.listExtensions({
				scope: scope,
				type: Extension_Type.INSTANCE
			});
			instanceExtensions.set(response.Extensions);
			isInstanceExtensionsLoaded = true;
		} catch (error) {
			console.error('Failed to fetch instance extensions:', error);
		}
	}

	async function fetchStorageExtensions() {
		try {
			const response = await orchestratorClient.listExtensions({
				scope: scope,
				type: Extension_Type.STORAGE
			});
			storageExtensions.set(response.Extensions);
			isStorageExtensionsLoaded = true;
		} catch (error) {
			console.error('Failed to fetch storage extensions:', error);
		}
	}

	onMount(() => {
		Promise.all([
			fetchMetricsExtensions(),
			fetchServiceMeshExtensions(),
			fetchRegistryExtensions(),
			fetchModelExtensions(),
			fetchInstanceExtensions(),
			fetchStorageExtensions()
		]);
	});
</script>

<Accordion.Root
	type="multiple"
	class="group w-full overflow-hidden rounded-lg border bg-card text-card-foreground transition-all duration-300 **:data-[slot='accordion-trigger']:p-6"
	value={getAccordionValue()}
>
	{#if isMetricsExtensionsLoaded && $metricsExtensions.filter((metricsExtension) => !metricsExtension.latest).length == 0}
		<Accordion.Item value="metrics">
			<Accordion.Trigger>
				<Thumbnail
					{scope}
					extensionsBundle="metrics"
					extensions={metricsExtensions}
					updator={fetchMetricsExtensions}
				/>
			</Accordion.Trigger>
			<Accordion.Content>
				{#each $metricsExtensions as metricsExtension, index (index)}
					<Node extension={metricsExtension} alignment={index % 2 ? 'right' : 'left'} />
				{/each}
			</Accordion.Content>
		</Accordion.Item>
	{/if}

	{#if isServiceMeshExtensionsLoaded && $serviceMeshExtensions.filter((serviceMeshExtension) => !serviceMeshExtension.latest).length == 0}
		<Accordion.Item value="serviceMesh">
			<Accordion.Trigger>
				<Thumbnail
					{scope}
					extensionsBundle="serviceMesh"
					extensions={serviceMeshExtensions}
					updator={fetchServiceMeshExtensions}
				/>
			</Accordion.Trigger>
			<Accordion.Content>
				{#each $serviceMeshExtensions as serviceMeshExtension, index (index)}
					<Node extension={serviceMeshExtension} alignment={index % 2 ? 'right' : 'left'} />
				{/each}
			</Accordion.Content>
		</Accordion.Item>
	{/if}

	{#if isRegistryExtensionsLoaded && $registryExtensions.filter((registryExtension) => !registryExtension.latest).length == 0}
		<Accordion.Item value="registry">
			<Accordion.Trigger>
				<Thumbnail
					{scope}
					extensionsBundle="registry"
					extensions={registryExtensions}
					updator={fetchRegistryExtensions}
				/>
			</Accordion.Trigger>
			<Accordion.Content>
				{#each $registryExtensions as registryExtension, index (index)}
					<Node extension={registryExtension} alignment={index % 2 ? 'right' : 'left'} />
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

	{#if isStorageExtensionsLoaded && $storageExtensions.filter((storageExtension) => !storageExtension.latest).length == 0}
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
</Accordion.Root>
