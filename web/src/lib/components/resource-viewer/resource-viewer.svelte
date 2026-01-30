<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { Ban, Braces } from '@lucide/svelte';
	import File from '@lucide/svelte/icons/file';
	import Layers from '@lucide/svelte/icons/layers';
	import { getContext, onDestroy } from 'svelte';
	import { stringify } from 'yaml';

	import { type GetRequest, ResourceService } from '$lib/api/resource/v1/resource_pb';
	import * as Code from '$lib/components/custom/code';
	import { typographyVariants } from '$lib/components/typography/index.ts';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as Item from '$lib/components/ui/item';
	import Label from '$lib/components/ui/label/label.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';

	import type { ViewerType } from './viewers';
	import { getResourceViewer } from './viewers';

	let {
		cluster,
		namespace,
		group,
		version,
		kind,
		resource,
		name
	}: {
		cluster: string;
		namespace: string;
		group: string;
		version: string;
		kind: string;
		resource: string;
		name: string;
	} = $props();

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);

	let getAbortController: AbortController | null = null;
	async function GetResource(): Promise<any> {
		getAbortController = new AbortController();
		try {
			const response = await resourceClient.get(
				{
					cluster,
					namespace,
					group,
					version,
					resource,
					name
				} as GetRequest,
				{ signal: getAbortController?.signal }
			);

			return response.object;
		} finally {
			if (getAbortController) getAbortController = null;
		}
	}

	onDestroy(() => {
		if (getAbortController) {
			getAbortController.abort();
		}
	});
</script>

{#await GetResource()}
	<Field.Group class="pb-8">
		<Field.Set>
			<Item.Root>
				<Item.Media>
					<Skeleton class="size-10" />
				</Item.Media>
				<Item.Content>
					<Item.Description>
						<Skeleton class="h-7 w-1/6" />
					</Item.Description>
					<Item.Title>
						<Skeleton class="h-5 w-[10vw]" />
					</Item.Title>
					<div class="grid grid-cols-3">
						{#each Array(3)}
							<Item.Root class="p-0">
								<Item.Content>
									<Item.Description>
										<Skeleton class="h-5 w-1/6" />
									</Item.Description>
									<Item.Title>
										<Skeleton class="h-3 w-[10vw]" />
									</Item.Title>
								</Item.Content>
							</Item.Root>
						{/each}
					</div>
				</Item.Content>
				<Item.Actions>
					<Skeleton class="size-10" />
				</Item.Actions>
				<Item.Footer class="flex flex-col items-start space-y-4">
					{#each Array(2)}
						<Item.Root class="p-0">
							<Item.Media>
								<Skeleton class="h-5 w-10" />
							</Item.Media>
							<Item.Content>
								<Item.Title>
									{#each Array(3)}
										<Skeleton class="h-3 w-30" />
									{/each}
								</Item.Title>
							</Item.Content>
						</Item.Root>
					{/each}
				</Item.Footer>
			</Item.Root>
		</Field.Set>
		<Field.Set>
			{#each Array(13).keys() as index (index)}
				{#if index % 2 === 0}
					{#if index % 3 !== 0}
						{#if index % 5 === 0}
							{#if index % 7 !== 0}
								{#if index % 11 === 0}
									<Skeleton class="h-1 w-full" />
								{:else}
									<Skeleton class="h-11 w-5/6" />
								{/if}
							{:else}
								<Skeleton class="h-7 w-4/5" />
							{/if}
						{:else}
							<Skeleton class="h-5 w-3/4" />
						{/if}
					{:else}
						<Skeleton class="h-3 w-2/3" />
					{/if}
				{:else}
					<Skeleton class="h-2 w-1/2" />
				{/if}
			{/each}
		</Field.Set>
	</Field.Group>
{:then object}
	{@const Inspector: ViewerType = getResourceViewer(resource)}
	<Field.Group class="pb-8">
		<Field.Set>
			<!-- Header -->
			<Item.Root class="w-full p-0">
				<Item.Media variant="image" class="bg-muted-foreground/50 p-2">
					<Layers size={48} />
				</Item.Media>
				<Item.Content>
					<Item.Description>
						<Badge variant="outline">{object?.kind}</Badge>
						{object?.apiVersion}
					</Item.Description>
					<Item.Title class={typographyVariants({ variant: 'h3' })}>
						{object?.metadata?.name}
					</Item.Title>
					<Separator class="invisible" />
					<div class="grid gap-2 md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-6">
						{#if object?.metadata}
							{@const clusterData = { name: 'Cluster', information: cluster }}
							{@const namespaceData = { name: 'Namespace', information: namespace }}
							{@const creationTimestampData = {
								name: 'Creation Timestamp',
								information: new Date(object.metadata?.creationTimestamp).toLocaleString('sv-SE')
							}}
							{@const generationData = {
								name: 'Generation',
								information: object.metadata?.generation
							}}
							{@const resourceVersionData = {
								name: 'Resource Version',
								information: object.metadata?.resourceVersion
							}}
							{#each [clusterData, namespaceData, creationTimestampData, generationData, resourceVersionData] as data, index (index)}
								{#if data.information}
									<Item.Root class="p-0">
										<Item.Content>
											<Item.Description>
												{data.name}
											</Item.Description>
											<Item.Title>
												{data.information}
											</Item.Title>
										</Item.Content>
									</Item.Root>
								{/if}
							{/each}
						{/if}
					</div>
				</Item.Content>
				<Item.Actions>
					<Sheet.Root>
						<Sheet.Trigger>
							<Button variant="outline" size="icon-lg">
								<File />
							</Button>
						</Sheet.Trigger>
						<Sheet.Content
							side="right"
							class="flex h-full max-w-[62vw] min-w-[50vw] flex-col gap-0 overflow-y-auto p-4"
						>
							<Sheet.Header class="shruk-0 space-y-4">
								<Sheet.Title>
									{name}
									<p class="text-muted-foreground">{group}/{version}/{kind}/{resource}</p>
								</Sheet.Title>
								<Sheet.Description></Sheet.Description>
							</Sheet.Header>
							{#if object}
								<Code.Root
									code={stringify(object)}
									lang="yaml"
									class="no-shiki-limit m-4 border-none bg-muted"
								/>
							{:else}
								<Empty.Root class="m-4 bg-muted/50">
									<Empty.Header>
										<Empty.Media variant="icon">
											<Braces size={36} />
										</Empty.Media>
										<Empty.Title>No Data</Empty.Title>
										<Empty.Description>
											No data is currently available for this resource.
											<br />
											To populate this resource, please add properties or values through the resource
											editor.
										</Empty.Description>
									</Empty.Header>
									<Empty.Content></Empty.Content>
								</Empty.Root>
							{/if}
						</Sheet.Content>
					</Sheet.Root>
				</Item.Actions>
				<Item.Footer class="flex flex-col items-start justify-start gap-2">
					<!-- Tags -->
					{@const tags = {
						Labels: object?.metadata?.labels ?? {},
						Annotations: object?.metadata?.annotations ?? {}
					}}
					{#each Object.entries(tags) as [key, values], index (index)}
						{#if Object.keys(values).length > 0}
							<Item.Root class="grid w-full grid-cols-[80px_1fr] p-0">
								<Item.Media class="relative flex w-fit items-center self-start">
									<Label>{key}</Label>
								</Item.Media>
								<Item.Content class="group flex flex-row flex-wrap gap-2">
									{#each Object.entries(values) as [key, value], index (index)}
										<Badge variant="outline" class="max-w-full border">
											<p class="text-muted-foreground">{key}</p>
											<Separator orientation="vertical" class="h-1" />
											<p class="max-w-xs truncate">{value}</p>
										</Badge>
									{/each}
								</Item.Content>
							</Item.Root>
						{/if}
					{/each}
				</Item.Footer>
			</Item.Root>
		</Field.Set>
		<Inspector {object} />
	</Field.Group>
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

<style>
	@reference '../../../app.css';

	:global(.no-shiki-limit pre.shiki:not([data-code-overflow] *):not([data-code-overflow])) {
		overflow-y: visible !important;
		max-height: none !important;
	}
</style>
