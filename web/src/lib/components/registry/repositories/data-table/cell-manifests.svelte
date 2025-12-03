<script lang="ts">
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import {
		type Chart,
		type Image,
		type Manifest,
		RegistryService,
		type Repository
	} from '$lib/api/registry/v1/registry_pb';
	import { fuzzLogosIcon } from '$lib/components/applications/store/commerce-store/utils';
	import CopyButton from '$lib/components/custom/copy-button/copy-button.svelte';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import * as Avatar from '$lib/components/ui/avatar';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Item from '$lib/components/ui/item/index.js';
	import Label from '$lib/components/ui/label/label.svelte';
	import * as Sheet from '$lib/components/ui/sheet';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

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
		<Sheet.Trigger disabled={!repository.manifestCount}>
			<div class="flex items-center gap-1">
				{repository.manifestCount}
				<Icon icon="ph:arrow-square-out" />
			</div>
		</Sheet.Trigger>
		<Sheet.Content side="right" class="min-w-[23vw] p-4">
			<Sheet.Header>
				{@const { value: sizeValue, unit: sizeUnit } = formatCapacity(repository.sizeBytes)}
				<div class="flex items-center gap-2">
					<div class="h-fit w-fit rounded-full bg-muted p-3">
						<Icon icon="ph:package" class="size-8" />
					</div>
					<div class="items-between flex flex-col justify-between gap-1">
						<h1>{repository.name}</h1>
						<p class="text-xs text-muted-foreground">{sizeValue} {sizeUnit}</p>
					</div>
				</div>
			</Sheet.Header>
			{#each $manifests as manifest (manifest.digest)}
				<Dialog.Root>
					<Dialog.Trigger class="hover:underline">
						<Item.Root class="w-full">
							<Item.Media variant="icon">
								{#if manifest.config.case === 'chart'}
									<Icon icon="logos:helm" />
								{:else if manifest.config.case === 'image'}
									<Icon icon="logos:docker-icon" />
								{/if}
							</Item.Media>
							<Item.Content class="flex flex-col items-start">
								<Item.Title>
									<p class="nowrap max-w-40 overflow-hidden text-ellipsis whitespace-nowrap">
										{manifest.repositoryName}
									</p>
								</Item.Title>
								<Item.Description class="text-xs">
									<span class="flex items-center gap-1 text-muted-foreground">
										<Icon icon="ph:tag" />
										<p>{manifest.tag}</p>
									</span>
								</Item.Description>
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
													fetchManifests();
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
						</Item.Root>
					</Dialog.Trigger>
					<Dialog.Content class="max-h-[77vh] overflow-y-auto  p-2">
						{#if manifest.config.case === 'chart'}
							{@const chart: Chart = manifest.config.value}
							<Card.Root class="border-none shadow-none">
								<Card.Header>
									<Card.Title>
										<Item.Root class="p-0">
											<Item.Media>
												<Avatar.Root>
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
												<Item.Title>
													{chart.name}
												</Item.Title>
												<Item.Description>
													{chart.version}
												</Item.Description>
											</Item.Content>
										</Item.Root>
									</Card.Title>
									<Card.Description class="py-2 text-sm">
										{chart.description}
									</Card.Description>
									{#if chart.keywords || chart.tags}
										<div class="flex flex-wrap items-center gap-2">
											{#if chart.tags}
												<span class="flex items-center gap-1 text-xs text-muted-foreground">
													<Icon icon="ph:tag" class="size-3" />
													<p>{chart.tags}</p>
												</span>
											{/if}
											{#each chart.keywords as keyword (keyword)}
												<span class="flex items-center gap-1 text-xs text-muted-foreground">
													<Icon icon="ph:tag" class="size-3" />
													<p>{keyword}</p>
												</span>
											{/each}
										</div>
									{/if}
									<Card.Action>
										{#if chart.repositoryName.startsWith('otterscale/')}
											<Icon icon="ph:star-fill" class="size-4 fill-yellow-400 text-yellow-400" />
										{/if}
										{#if chart.deprecated}
											<Icon icon="ph:prohibit-fill" class="size-4 fill-red-500 text-yellow-400" />
										{/if}
									</Card.Action>
								</Card.Header>
								<Card.Content class="space-y-4">
									{#if chart.appVersion || chart.kubeVersion}
										<Label>{m.versions()}</Label>
										<div class="grid grid-cols-1 gap-2 md:grid-cols-2">
											<Item.Root class="p-0">
												<Item.Media variant="icon" class="bg-transparent">
													<Icon icon="ph:git-commit" class="size-4" />
												</Item.Media>
												<Item.Content>
													<Item.Description class="text-xs">{m.application()}</Item.Description>
													<Item.Title class="text-xs">{chart.appVersion}</Item.Title>
												</Item.Content>
											</Item.Root>

											<Item.Root class="p-0">
												<Item.Media variant="icon" class="bg-transparent">
													<Icon icon="ph:git-commit" class="size-4" />
												</Item.Media>
												<Item.Content>
													<Item.Description class="text-xs">{m.kubernetes()}</Item.Description>
													<Item.Title class="text-xs">{chart.kubeVersion}</Item.Title>
												</Item.Content>
											</Item.Root>
										</div>
									{/if}

									{#if chart.dependencies.length > 0}
										<Label>{m.dependencies()}</Label>
										<div class="grid grid-cols-1 gap-2 md:grid-cols-2">
											{#each chart.dependencies as dependency, index (index)}
												<!-- eslint-disable svelte/no-navigation-without-resolve -->
												<a
													href={dependency.repository}
													class="col-span-1 hover:underline md:col-span-2"
													target="_blank"
												>
													<Item.Root class="p-0 ">
														<Item.Media variant="icon" class="bg-transparent">
															<Icon icon="ph:cube" class="size-4" />
														</Item.Media>
														<Item.Content>
															<Item.Description class="text-xs"
																>{dependency.name}@{dependency.version}</Item.Description
															>
															<Item.Title class="text-xs">{dependency.condition}</Item.Title>
														</Item.Content>
													</Item.Root>
												</a>
												<!-- eslint-enable svelte/no-navigation-without-resolve -->
											{/each}
										</div>
									{/if}

									{#if chart.maintainers.length > 0}
										<Label>{m.maintainers()}</Label>
										<div class="grid grid-cols-1 gap-2 md:grid-cols-2">
											{#each chart.maintainers as maintainer, index (index)}
												<!-- eslint-disable svelte/no-navigation-without-resolve -->
												<a href={maintainer.url} class="hover:underline" target="_blank">
													<Item.Root class="p-0">
														<Item.Media variant="icon" class="bg-transparent">
															<Icon icon="ph:user" class="size-4" />
														</Item.Media>
														<Item.Content>
															<Item.Description class="text-xs">{m.maintainer()}</Item.Description>
															<Item.Title class="text-xs">{maintainer.name}</Item.Title>
														</Item.Content>
													</Item.Root>
												</a>
												<!-- eslint-enable svelte/no-navigation-without-resolve -->
											{/each}
										</div>
									{/if}
								</Card.Content>
								<Card.Footer class="flex flex-wrap gap-4">
									{#if chart.home}
										<Tooltip.Provider>
											<Tooltip.Root>
												<Tooltip.Trigger>
													<Button
														variant="secondary"
														size="icon-sm"
														href={chart.home}
														target="_blank"
													>
														<Icon icon="ph:house-bold" class="size-4" />
													</Button>
												</Tooltip.Trigger>
												<Tooltip.Content>
													{chart.home}
												</Tooltip.Content>
											</Tooltip.Root>
										</Tooltip.Provider>
									{/if}
									{#each chart.sources as source, index (index)}
										<Tooltip.Provider>
											<Tooltip.Root>
												<Tooltip.Trigger>
													<Button variant="secondary" size="icon-sm" href={source} target="_blank">
														<Icon icon="ph:code-bold" class="size-4" />
													</Button>
												</Tooltip.Trigger>
												<Tooltip.Content>
													{source}
												</Tooltip.Content>
											</Tooltip.Root>
										</Tooltip.Provider>
									{/each}
								</Card.Footer>
							</Card.Root>
						{:else if manifest.config.case === 'image'}
							{@const image: Image = manifest.config.value}
							<Card.Root class="border-none shadow-none">
								<Card.Header>
									<Card.Title>
										<Item.Root class="p-0">
											<Item.Media variant="image">
												<Icon icon="logos:docker-icon" class="size-6" />
											</Item.Media>
											<Item.Content>
												<Item.Title class="overflow-hidden text-ellipsis whitespace-nowrap">
													{manifest.repositoryName}
												</Item.Title>
												<Item.Description>
													{manifest.tag}
												</Item.Description>
											</Item.Content>
										</Item.Root>
									</Card.Title>
									<Card.Description class="space-y-2 py-2 text-xs">
										{#if image.author}
											<div class="grow">
												<h3 class="font-semibold text-primary">Author</h3>
												<p class="mt-1 break-all">{image.author}</p>
											</div>
										{/if}
										{#if image.createdAt}
											<div class="grow">
												<h3 class="font-semibold text-primary">Create Time</h3>
												<p class="mt-1 break-all">
													{timestampDate(image.createdAt).toDateString()}
												</p>
											</div>
										{/if}
										{#if image.config && image.config.labels}
											{#each Object.entries(image.config.labels) as [label, value] (label)}
												<div class="grow">
													<h3 class="font-semibold text-primary capitalize">{label}</h3>
													<p class="mt-1 break-all">
														{value}
													</p>
												</div>
											{/each}
										{/if}
									</Card.Description>
									{#if image.platform}
										<div class="flex flex-wrap items-center gap-4 text-xs text-muted-foreground">
											{#if image.platform.architecture}
												<span class="flex items-center gap-1">
													<Icon icon="ph:cpu" class="size-4" />
													<p>{image.platform.architecture}</p>
												</span>
											{/if}
											{#if image.platform.variant}
												<span class="flex items-center gap-1">
													<Icon icon="ph:cpu" class="size-4" />
													<p>{image.platform.variant}</p>
												</span>
											{/if}
											{#if image.platform.os && image.platform.variant}
												<span class="flex items-center gap-1">
													<Icon icon="ph:squares-four" class="size-4" />
													<p>{image.platform.os} {image.platform.variant}</p>
												</span>
											{/if}
										</div>
									{/if}
									<Card.Action>
										{#if manifest.repositoryName.startsWith('otterscale/')}
											<Icon icon="ph:star-fill" class="size-4 fill-yellow-400 text-yellow-400" />
										{/if}
									</Card.Action>
								</Card.Header>
								<Card.Content class="space-y-4">
									{#if image.config}
										<Label>{m.configuration()}</Label>
										<div class="grid grid-cols-1 gap-2 md:grid-cols-2">
											<Item.Root class="p-0">
												<Item.Media variant="icon" class="bg-transparent">
													<Icon icon="ph:briefcase" class="size-4" />
												</Item.Media>
												<Item.Content>
													<Item.Description class="text-xs"
														>{m.working_directory()}</Item.Description
													>
													<Item.Title class="text-xs">{image.config.workingDir || '/'}</Item.Title>
												</Item.Content>
											</Item.Root>

											<Item.Root class="p-0">
												<Item.Media variant="icon" class="bg-transparent">
													<Icon icon="ph:plug" class="size-4" />
												</Item.Media>
												<Item.Content class="h-full">
													<Item.Description class="text-xs">{m.export_ports()}</Item.Description>
													<Item.Title class="text-xs">
														{#each image.config.exposedPorts as port (port)}
															<p>{port}</p>
														{/each}
													</Item.Title>
												</Item.Content>
											</Item.Root>

											<Item.Root class="col-span-1 p-0 md:col-span-2">
												<Item.Media variant="icon" class="bg-transparent">
													<Icon icon="ph:sign-in" class="size-4" />
												</Item.Media>
												<Item.Content class="h-full">
													<Item.Description class="text-xs">{m.entrypoint()}</Item.Description>
													<Item.Title class="flex flex-col gap-1 text-xs">
														{#each image.config.entrypoint as point (point)}
															<span class="group flex w-full items-center gap-1">
																<p class="font-mono">{point}</p>
																<CopyButton
																	class="ml-auto size-4 opacity-0 transition-opacity duration-200 group-hover:opacity-100"
																	text={point}
																/>
															</span>
														{/each}
													</Item.Title>
												</Item.Content>
											</Item.Root>

											<Item.Root class="col-span-1 p-0 md:col-span-2">
												<Item.Media variant="icon" class="bg-transparent">
													<Icon icon="ph:terminal" class="size-4" />
												</Item.Media>
												<Item.Content class="h-full">
													<Item.Description class="text-xs">{m.cmd()}</Item.Description>
													<Item.Title class="text-xs">
														<span class="group flex w-full items-center gap-1">
															<p class="font-mono">
																{image.config.cmd.join(' ')}
															</p>
															<CopyButton
																class="ml-auto size-4 opacity-0 transition-opacity duration-200 group-hover:opacity-100"
																text={image.config.cmd.join(' ')}
															/>
														</span>
													</Item.Title>
												</Item.Content>
											</Item.Root>

											{#if image.config.volumes.length > 0}
												<Item.Root class="col-span-1 p-0 md:col-span-2">
													<Item.Media variant="icon" class="bg-transparent">
														<Icon icon="ph:folder" class="size-4" />
													</Item.Media>
													<Item.Content class="h-full">
														<Item.Description class="text-xs">{m.volumes()}</Item.Description>
														<Item.Title class="flex flex-col gap-1 text-xs">
															{#each image.config.volumes as volume (volume)}
																<span class="group flex w-full items-center gap-1">
																	<p class="font-mono">{volume}</p>
																	<CopyButton
																		class="ml-auto size-4 opacity-0 transition-opacity duration-200 group-hover:opacity-100"
																		text={volume}
																	/>
																</span>
															{/each}
														</Item.Title>
													</Item.Content>
												</Item.Root>
											{/if}
										</div>
									{/if}
								</Card.Content>
								<Card.Footer></Card.Footer>
							</Card.Root>
						{/if}
					</Dialog.Content>
				</Dialog.Root>
			{/each}
		</Sheet.Content>
	</Sheet.Root>
{/if}
