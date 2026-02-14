<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { Ban, Boxes } from '@lucide/svelte';
	import { getContext } from 'svelte';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import {
		type APIResource,
		type DiscoveryRequest,
		ResourceService
	} from '$lib/api/resource/v1/resource_pb';
	import KindViewer from '$lib/components/kind-viewer/kind-viewer.svelte';
	import Picker from '$lib/components/kind-viewer/kind-viewer-picker.svelte';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import * as Item from '$lib/components/ui/item';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import * as Table from '$lib/components/ui/table/index.js';
	import { Toggle } from '$lib/components/ui/toggle/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { breadcrumbs } from '$lib/stores';

	breadcrumbs.set([
		{
			title: 'manifest',
			url: resolve('/')
		}
	]);

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

	let clustered = $state(false);
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
				{#each Array(3)}
					<Skeleton class="size-7" />
				{/each}
			</div>
			<div class="rounded-lg border">
				<Table.Root class="w-full">
					<Table.Header>
						<Table.Row>
							{#each Array(5)}
								<Table.Head class="p-4">
									<Skeleton class="h-7 w-full" />
								</Table.Head>
							{/each}
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each Array(13)}
							<Table.Row class="border-none">
								{#each Array(5)}
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
					{#each Array(3)}
						<Skeleton class="size-10" />
					{/each}
				</div>
			</div>
		</div>
	{:then apiResources}
		{@const apiResourceOptions = apiResources
			.filter((apiResource) => apiResource.resource.indexOf('/') < 0)
			.map((apiResource) => ({
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
							{!group ? 'core' : group}/{version}
						</Item.Description>
					</Item.Content>
				</Item.Root>
				<div class="flex items-center gap-2">
					<Tooltip.Provider>
						<Tooltip.Root>
							<Tooltip.Trigger class={page.data.user.roles.includes('admin') ? 'flex' : 'hidden'}>
								<Toggle
									bind:pressed={clustered}
									aria-label="switch clustered"
									variant="outline"
									class="data-[state=on]:*:bg-transparent"
								>
									<Boxes size={16} />
									All Namespaces
								</Toggle>
							</Tooltip.Trigger>
							<Tooltip.Content>Press to view resources across all namespaces.</Tooltip.Content>
						</Tooltip.Root>
					</Tooltip.Provider>
					<Picker
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
				</div>
			</div>
			{#key clustered + resource + namespace}
				{#if selectedAPIResource}
					{@const apiResource = selectedAPIResource}
					<KindViewer
						{apiResource}
						{cluster}
						{group}
						{version}
						{kind}
						{resource}
						{namespace}
						{clustered}
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
				<Alert.Root variant="destructive" class="border-none bg-destructive/5">
					<Alert.Title class="font-bold">{error?.name}</Alert.Title>
					<Alert.Description class="text-start">
						{error?.rawMessage}
					</Alert.Description>
				</Alert.Root>
				<div class="flex gap-4">
					<Button variant="outline" onclick={() => history.back()}>Go Back</Button>
					<Button href="/">Go Home</Button>
				</div>
			</Empty.Content>
		</Empty.Root>
	{/await}
{/key}
