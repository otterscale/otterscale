<script lang="ts">
	import type { JsonObject } from '@bufbuild/protobuf';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { Ban } from '@lucide/svelte';
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
	import * as Alert from '$lib/components/ui/alert/index.js';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import * as Item from '$lib/components/ui/item';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import * as Table from '$lib/components/ui/table/index.js';

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
	{#await fetchAPIResources(cluster, group, version, kind)}
		<div class="space-y-4">
			<div class="space-y-4">
				<Skeleton class="h-13 w-1/3" />
				<Skeleton class="h-7 w-1/5" />
			</div>
			<div class="flex items-center gap-2">
				<Skeleton class="h-7 w-full" />
				{#each Array(3) as i (i)}
					<Skeleton class="size-7" />
				{/each}
			</div>
			<div class="rounded-lg border">
				<Table.Root class="w-full">
					<Table.Header>
						<Table.Row>
							{#each Array(5) as i (i)}
								<Table.Head class="p-4">
									<Skeleton class="h-7 w-full" />
								</Table.Head>
							{/each}
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each Array(13) as i (i)}
							<Table.Row class="border-none">
								{#each Array(5) as i (i)}
									<Table.Cell>
										<Skeleton class="h-7 w-full" />
									</Table.Cell>
								{/each}
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</div>
			<div class="flex items-center justify-between gap-4">
				<Skeleton class="h-7 w-1/5" />
				<div class="flex items-center gap-2">
					{#each Array(3) as i (i)}
						<Skeleton class="size-10" />
					{/each}
				</div>
			</div>
		</div>
	{:then apiResources}
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
									String(lodash.get(previous, 'metadata.name') ?? '').localeCompare(
										String(lodash.get(next, 'metadata.name') ?? '')
									)
								)
								.map((namespace: JsonObject | undefined) => ({
									icon: 'ph:cube',
									label: String(lodash.get(namespace, 'metadata.name') ?? ''),
									value: String(lodash.get(namespace, 'metadata.name') ?? ''),
									description: String(lodash.get(namespace, 'status.phase') ?? '')
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
				{#if selectedAPIResource}
					{@const apiResource = selectedAPIResource}
					<ResourceManager
						{apiResource}
						{cluster}
						{group}
						{version}
						{kind}
						{resource}
						{namespace}
					/>
				{/if}
			{/key}
		</div>
	{:catch error}
		<Empty.Root>
			<Empty.Header>
				<Empty.Media class="rounded-full bg-muted p-4">
					<Ban size={36} />
				</Empty.Media>
				<Empty.Title class="text-2xl font-bold">Failed to load data</Empty.Title>
				<Empty.Description>
					An error occurred while fetching data. Please check your connection or try again later.
				</Empty.Description>
			</Empty.Header>
			<Empty.Content>
				<Alert.Root
					variant="destructive"
					class="border border-destructive bg-destructive/10 text-left"
				>
					<Alert.Title>Message</Alert.Title>
					<Alert.Description>
						{error?.message ?? error?.toString() ?? JSON.stringify(error)}
					</Alert.Description>
				</Alert.Root>
				<Button href="/">Go Back</Button>
			</Empty.Content>
		</Empty.Root>
	{/await}
{/key}
