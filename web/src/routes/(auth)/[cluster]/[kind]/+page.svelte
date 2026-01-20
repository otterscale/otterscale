<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';

	import { page } from '$app/state';
	import {
		type APIResource,
		type DiscoveryRequest,
		type ListRequest,
		ResourceService
	} from '$lib/api/resource/v1/resource_pb';
	import ResourceManager from '$lib/components/dynamical-table/resource-manager.svelte';
	import ResourcePicker from '$lib/components/dynamical-table/resource-picker.svelte';
	import * as Item from '$lib/components/ui/item';

	const cluster = $derived(page.params.cluster ?? '');
	const group = $derived(page.url.searchParams.get('group') ?? '');
	const version = $derived(page.url.searchParams.get('version') ?? '');
	const kind = $derived(page.params.kind ?? '');

	const transport: Transport = getContext('transport');
	const client = createClient(ResourceService, transport);

	let apiResources = $state<APIResource[]>([]);
	let selectedAPIResourceResource = $state('');
	const selectedAPIResource = $derived(
		apiResources.find(
			(apiResource) =>
				apiResource &&
				apiResource.group === group &&
				apiResource.version === version &&
				apiResource.kind === kind &&
				apiResource.resource === selectedAPIResourceResource
		)
	);
	async function fetchAPIResources(cluster: string, group: string, version: string, kind: string) {
		const response = await client.discovery({
			cluster: cluster
		} as DiscoveryRequest);

		apiResources = response.apiResources.filter(
			(apiResource) =>
				apiResource &&
				apiResource.group === group &&
				apiResource.version === version &&
				apiResource.kind === kind
		);
		return apiResources;
	}

	let selectedNamespaceMetadataName = $state(page.url.searchParams.get('namespace') ?? '');
	async function fetchNamespaces(cluster: string) {
		const response = await client.list({
			cluster: cluster,
			resource: 'namespaces',
			version: 'v1'
		} as ListRequest);
		return response.items.map((item) => item.object);
	}
</script>

{#key cluster + group + version + kind}
	{#await fetchAPIResources(cluster, group, version, kind) then apiResources}
		{@const apiResourceOptions = apiResources.map((apiResource) => ({
			icon: 'lucide:list',
			label: apiResource.resource,
			value: apiResource.resource,
			description: `${apiResource.group}/${apiResource.version}/${apiResource.kind}`
		}))}
		{@const [initialAPIResourceOption] = apiResourceOptions}
		<div class="space-y-4">
			<div class="flex items-end justify-between gap-4">
				<Item.Root class="p-0">
					<Item.Content class="text-left">
						<Item.Title class="text-xl font-bold">
							{kind}
						</Item.Title>
						<Item.Description class="text-base">
							{cluster}
							{group}/{version}
						</Item.Description>
					</Item.Content>
				</Item.Root>
				<div class="flex items-center gap-4">
					<ResourcePicker
						class="w-fit"
						bind:value={selectedAPIResourceResource}
						initialValue={initialAPIResourceOption.value}
						options={apiResourceOptions}
					/>
					{#if selectedAPIResource && selectedAPIResource.namespaced}
						{#await fetchNamespaces(cluster) then namespaces}
							{@const namespaceOptions = namespaces.map((namespace: any) => ({
								icon: 'lucide:list',
								label: namespace?.metadata?.name,
								value: namespace?.metadata?.name,
								description: namespace?.status?.phase
							}))}
							{@const [initialNamespaceOptions] = namespaceOptions}
							<ResourcePicker
								class="w-fit"
								bind:value={selectedNamespaceMetadataName}
								initialValue={initialNamespaceOptions.value}
								options={namespaceOptions}
							/>
						{/await}
					{/if}
				</div>
			</div>
			{#if selectedAPIResource}
				{#key selectedAPIResourceResource + selectedNamespaceMetadataName}
					{@const resource = selectedAPIResourceResource}
					{@const namespace = selectedAPIResource.namespaced
						? selectedNamespaceMetadataName
						: undefined}
					<ResourceManager {cluster} {group} {version} {kind} {resource} {namespace} />
				{/key}
			{/if}
		</div>
	{:catch}
		Error
	{/await}
{/key}
