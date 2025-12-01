<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import {
		type Manifest,
		RegistryService,
		type Repository
	} from '$lib/api/registry/v1/registry_pb';
	import { fuzzLogosIcon } from '$lib/components/applications/store/commerce-store/utils';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import * as Accordion from '$lib/components/ui/accordion/index.js';
	import * as Avatar from '$lib/components/ui/avatar';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Item from '$lib/components/ui/item/index.js';
	import * as Sheet from '$lib/components/ui/sheet';
	import { formatCapacity } from '$lib/formatter';
</script>

<script lang="ts">
	let {
		repository,
		scope,
		reloadManager
	}: {
		repository: Repository;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();

	const transport: Transport = getContext('transport');
	const registryClient = createClient(RegistryService, transport);

	const manifests = writable<Manifest[]>([]);
	async function fetchManifests() {
		try {
			const response = await registryClient.listManifests({
				scope,
				repositoryName: repository.name
			});
			manifests.set(response.manifests);
		} catch (error) {
			console.error('Failed to fetch manifests:', error);
		}
	}

	let isLoaded = $state(false);
	onMount(async () => {
		await fetchManifests();
		isLoaded = true;
	});
</script>

{#if isLoaded}
	<Sheet.Root>
		<Sheet.Trigger>
			<div class="flex items-center gap-1">
				{repository.manifestCount}
				<Icon icon="ph:squares-four" />
			</div>
		</Sheet.Trigger>
		<Sheet.Content side="right" class="min-w-[23vw] p-4">
			<Sheet.Header class="text-xl">Manifests</Sheet.Header>
			{#each $manifests as manifest (manifest.digest)}
				{@const { value: sizeValue, unit: sizeUnit } = formatCapacity(manifest.sizeBytes)}
				{#if manifest.config.case === 'chart'}
					{@const chart = manifest.config.value}
					<Accordion.Root type="single">
						<Accordion.Item value="item-1">
							<Accordion.Trigger>
								<div class="flex w-full items-center gap-4 text-left">
									<div class="shrink-0 rounded-lg bg-primary/10 p-2 text-primary">
										<Icon icon="ph:cube" class="size-6" />
									</div>
									<div class="grow">
										<h3 class="text-md font-semibold">Repository</h3>
										<p class="mt-1 text-sm text-muted-foreground">
											{manifest.repositoryName}
										</p>
									</div>
									<div class="hidden shrink-0 items-center gap-3 md:flex">
										<div class="flex items-center gap-1.5 text-xs text-muted-foreground">
											<p>{sizeValue} {sizeUnit}</p>
										</div>
									</div>
								</div>
							</Accordion.Trigger>
							<Accordion.Content class="space-y-4 rounded-lg bg-muted p-4">
								<div class="flex w-full max-w-lg flex-col gap-6">
									<Item.Root class="relative border-none" variant="outline">
										<Item.Media>
											<Avatar.Root class="size-6">
												<Avatar.Image src={chart.icon} class="object-contain" />
												<Avatar.Fallback>
													<Icon
														icon={fuzzLogosIcon(chart.name, 'fluent-emoji-flat:otter')}
														class="size-6"
													/>
												</Avatar.Fallback>
											</Avatar.Root>
										</Item.Media>
										<Item.Content>
											<Item.Title>{chart.name}</Item.Title>
											<Item.Description>{chart.description}</Item.Description>
										</Item.Content>
										<Item.Actions>
											<Button
												size="icon-sm"
												variant="ghost"
												class="text-destructive"
												onclick={() => {
													toast.promise(
														() =>
															registryClient.deleteManifest({
																scope,
																digest: manifest.digest,
																repositoryName: manifest.repositoryName
															}),
														{
															loading: 'Loading...',
															success: () => {
																reloadManager.force();
																return `Delete ${manifest.repositoryName} success`;
															},
															error: (error) => {
																let message = `Fail to delete ${manifest.repositoryName}`;
																toast.error(message, {
																	description: (error as ConnectError).message.toString(),
																	duration: Number.POSITIVE_INFINITY
																});
																return message;
															}
														}
													);
												}}
											>
												<Icon icon="ph:trash" />
											</Button>
										</Item.Actions>
										<span class="absolute top-0 right-0 flex items-center gap-1">
											{#if chart.repositoryName.startsWith('otterscale/')}
												<Icon icon="ph:star-fill" class="size-4 fill-yellow-400 text-yellow-400" />
											{/if}
											{#if chart.deprecated}
												<Icon icon="ph:prohibit-fill" class="size-4 fill-red-500 text-yellow-400" />
											{/if}
										</span>
									</Item.Root>
								</div>
								<!-- <div class="relative flex w-full items-start gap-2">
									<Avatar.Root class="size-6">
										<Avatar.Image src={chart.icon} class="object-contain" />
										<Avatar.Fallback>
											<Icon
												icon={fuzzLogosIcon(chart.name, 'fluent-emoji-flat:otter')}
												class="size-6"
											/>
										</Avatar.Fallback>
									</Avatar.Root>
									<span>
										<h3 class="font-semibold">{chart.name}</h3>
										<p class="flex items-center gap-1 text-sm text-muted-foreground">
											{chart.version}
										</p>
									</span>
									<span class="absolute top-0 right-0 flex items-center gap-1">
										{#if chart.repositoryName.startsWith('otterscale/')}
											<Icon icon="ph:star-fill" class="size-4 fill-yellow-400 text-yellow-400" />
										{/if}
										{#if chart.deprecated}
											<Icon icon="ph:prohibit-fill" class="size-4 fill-red-500 text-yellow-400" />
										{/if}
									</span>
								</div> -->
							</Accordion.Content>
						</Accordion.Item>
					</Accordion.Root>
				{/if}
			{/each}
		</Sheet.Content>
	</Sheet.Root>
{/if}
