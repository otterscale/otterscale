<script lang="ts">
	import type { JsonObject } from '@bufbuild/protobuf';
	import { createClient, type Transport } from '@connectrpc/connect';
	import lodash from 'lodash';
	import { getContext } from 'svelte';

	import { goto } from '$app/navigation';
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
	const kind = $derived(page.params.kind ?? '');
	const group = $derived(page.url.searchParams.get('group') ?? '');
	const version = $derived(page.url.searchParams.get('version') ?? '');
	const resource = $derived(page.url.searchParams.get('resource') ?? '');
	const namespace = $derived(page.url.searchParams.get('namespace') ?? '');

	const transport: Transport = getContext('transport');
	const client = createClient(ResourceService, transport);

	let apiResources = $state<APIResource[]>([]);
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
		console.log(apiResources);
		return apiResources;
	}
	const selectedAPIResource = $derived(
		apiResources.find(
			(apiResource) =>
				apiResource &&
				apiResource.group === group &&
				apiResource.version === version &&
				apiResource.kind === kind &&
				apiResource.resource === resource
		)
	);

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
			icon: 'ph:user',
			label: apiResource.resource,
			value: apiResource.resource,
			description: `${apiResource.group}/${apiResource.version}/${apiResource.kind}`
		}))}
		<div class="space-y-4">
			<div class="flex items-end justify-between gap-4">
				<Item.Root class="p-0">
					<Item.Content class="text-left">
						<Item.Title class="text-xl font-bold">
							{kind}
						</Item.Title>
						<Item.Description class="text-base">
							{group}/{version}
						</Item.Description>
					</Item.Content>
				</Item.Root>
				<div class="flex flex-row-reverse items-center gap-2">
					<ResourcePicker
						class="w-fit"
						resource="resource"
						value={resource}
						options={apiResourceOptions}
						onSelect={(option) => {
							page.url.searchParams.set('resource', option.value);
							// eslint-disable-next-line svelte/no-navigation-without-resolve
							goto(page.url.href);
						}}
					/>
					{#if selectedAPIResource && selectedAPIResource.namespaced}
						{#await fetchNamespaces(cluster) then namespaces}
							{@const namespaceOptions = namespaces
								.sort((previous: JsonObject | undefined, next: JsonObject | undefined) =>
									lodash
										.get(previous, 'metadata.name', '')
										.localeCompare(lodash.get(next, 'metadata.name', ''))
								)
								.map((namespace: JsonObject | undefined) => ({
									icon: 'ph:cube',
									label: lodash.get(namespace, 'metadata.name', ''),
									value: lodash.get(namespace, 'metadata.name', ''),
									description: lodash.get(namespace, 'status.phase', '')
								}))}
							<ResourcePicker
								class="w-fit"
								resource="namespace"
								value={namespace}
								options={namespaceOptions}
								onSelect={(option) => {
									page.url.searchParams.set('namespace', option.value);
									// eslint-disable-next-line svelte/no-navigation-without-resolve
									goto(page.url.href);
								}}
							/>
						{/await}
					{/if}
				</div>
			</div>
			{#key resource + namespace}
				<ResourceManager {cluster} {group} {version} {kind} {resource} {namespace} />
			{/key}
		</div>
	{:catch}
		Error
	{/await}
{/key}
