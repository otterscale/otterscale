<script lang="ts" module>
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
	import * as Code from '$lib/components/custom/code';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import * as Avatar from '$lib/components/ui/avatar';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Item from '$lib/components/ui/item/index.js';
	import * as Sheet from '$lib/components/ui/sheet';
	import { formatCapacity } from '$lib/formatter';
</script>

<script lang="ts">
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import { Popover } from 'bits-ui';

	import CopyButton from '$lib/components/custom/copy-button/copy-button.svelte';
	import Label from '$lib/components/ui/label/label.svelte';

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
				<Popover.Root>
					<Popover.Trigger class="hover:underline">
						<Item.Root class="w-full">
							<Item.Media class="rounded-full bg-primary/10 p-2 text-primary">
								{#if manifest.config.case === 'chart'}
									<Icon icon="logos:helm" class="size-6" />
								{:else if manifest.config.case === 'image'}
									<Icon icon="logos:docker-icon" class="size-6" />
								{/if}
							</Item.Media>
							<Item.Content class="flex flex-col items-start">
								<Item.Title>
									{manifest.repositoryName}
								</Item.Title>
								<Item.Description class="text-xs">{manifest.tag}</Item.Description>
							</Item.Content>
							<Item.Actions>
								<span class="flex items-center gap-1 text-muted-foreground">
									<Icon icon="ph:file" />
									<p>{sizeValue} {sizeUnit}</p>
								</span>
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
						</Item.Root>
					</Popover.Trigger>
					<Popover.Content align="start" side="left" class="m-4 max-w-sm">
						{#if manifest.config.case === 'chart'}
							{@const chart: Chart = manifest.config.value}
							<Card.Root class="min-w-96">
								<Card.Header>
									<Card.Title>
										<Item.Root class="p-0">
											<Item.Media variant="icon">
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
									<Card.Description>
										{chart.description}
									</Card.Description>
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
										<Label>Versions</Label>
										<div class="grid grid-cols-1 gap-2 md:grid-cols-2">
											<Item.Root class="p-0">
												<Item.Media variant="icon">
													<Icon icon="ph:git-commit" class="size-4" />
												</Item.Media>
												<Item.Content>
													<Item.Description class="text-xs">Application</Item.Description>
													<Item.Title class="text-xs">{chart.appVersion}</Item.Title>
												</Item.Content>
											</Item.Root>

											<Item.Root class="p-0">
												<Item.Media variant="icon">
													<Icon icon="ph:git-commit" class="size-4" />
												</Item.Media>
												<Item.Content>
													<Item.Description class="text-xs">Kubernetes</Item.Description>
													<Item.Title class="text-xs">{chart.kubeVersion}</Item.Title>
												</Item.Content>
											</Item.Root>
										</div>
									{/if}

									{#if chart.dependencies.length > 0}
										<Label>Dependencies</Label>
										<div class="grid grid-cols-1 gap-2 md:grid-cols-2">
											{#each chart.dependencies as dependency, index (index)}
												<a href={dependency.repository} class="hover:underline" target="_blank">
													<Item.Root class="p-0">
														<Item.Media variant="icon">
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
											{/each}
										</div>
									{/if}

									{#if chart.maintainers.length > 0}
										<Label>Maintainers</Label>
										<div class="grid grid-cols-1 gap-2 md:grid-cols-2">
											{#each chart.maintainers as maintainer, index (index)}
												<a href={maintainer.url} class="hover:underline" target="_blank">
													<Item.Root class="p-0">
														<Item.Media variant="icon">
															<Icon icon="ph:user" class="size-4" />
														</Item.Media>
														<Item.Content>
															<Item.Description class="text-xs">Maintainer</Item.Description>
															<Item.Title class="text-xs">{maintainer.name}</Item.Title>
														</Item.Content>
													</Item.Root>
												</a>
											{/each}
										</div>
									{/if}

									{#if chart.keywords || chart.tags}
										<Label>Tags & Keywords</Label>
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
								</Card.Content>
								<Card.Footer class="flex flex-wrap gap-2">
									{#if chart.home}
										<Button size="icon-sm" variant="outline" href={chart.home} target="_blank">
											<Icon icon="ph:house-bold" class="size-4" />
										</Button>
									{/if}
									{#each chart.sources as source, index (index)}
										<Button size="icon-sm" variant="outline" href={source} target="_blank">
											<Icon icon="ph:code-bold" class="size-4" />
										</Button>
									{/each}
								</Card.Footer>
							</Card.Root>
						{:else if manifest.config.case === 'image'}
							{@const image: Image = manifest.config.value}
							<Card.Root class="min-w-96">
								<Card.Header>
									<Card.Title>
										<Item.Root class="p-0">
											<Item.Media class="rounded-full bg-primary/10 p-2 text-primary">
												<Icon icon="logos:docker-icon" class="size-6" />
											</Item.Media>
											<Item.Content>
												<Item.Title>
													{manifest.repositoryName}
												</Item.Title>
												<Item.Description>
													{manifest.tag}
												</Item.Description>
											</Item.Content>
										</Item.Root>
									</Card.Title>
									<Card.Description></Card.Description>
									<Card.Action>
										{#if manifest.repositoryName.startsWith('otterscale/')}
											<Icon icon="ph:star-fill" class="size-4 fill-yellow-400 text-yellow-400" />
										{/if}
									</Card.Action>
								</Card.Header>
								<Card.Content class="space-y-4">
									<div class="grid grid-cols-1 gap-2 md:grid-cols-2">
										{#if image.author}
											<Item.Root class="p-0">
												<Item.Media variant="icon">
													<Icon icon="ph:user" class="size-4" />
												</Item.Media>
												<Item.Content>
													<Item.Description class="text-xs">Author</Item.Description>
													<Item.Title class="text-xs">{image.author}</Item.Title>
												</Item.Content>
											</Item.Root>
										{/if}
										{#if image.createdAt}
											{@const createTime = timestampDate(image.createdAt)}
											<Item.Root class="p-0">
												<Item.Media variant="icon">
													<Icon icon="ph:calendar" class="size-4" />
												</Item.Media>
												<Item.Content>
													<Item.Description class="text-xs">Created At</Item.Description>
													<Item.Title class="text-xs">
														{createTime.toLocaleString('en-GB', {
															year: 'numeric',
															month: '2-digit',
															day: '2-digit',
															hour: '2-digit',
															minute: '2-digit',
															hour12: false
														})}
													</Item.Title>
												</Item.Content>
											</Item.Root>
										{/if}
									</div>

									{#if image.platform}
										<Label>Platform</Label>
										<div class="grid grid-cols-1 gap-2 md:grid-cols-2">
											<Item.Root class="p-0">
												<Item.Media variant="icon">
													<Icon icon="ph:squares-four" class="size-4" />
												</Item.Media>
												<Item.Content>
													<Item.Description class="text-xs">
														{image.platform.architecture}
													</Item.Description>
													{#if image.platform.variant}
														<Item.Title class="text-xs">
															{image.platform.variant}
														</Item.Title>
													{/if}
												</Item.Content>
											</Item.Root>
											<Item.Root class="p-0">
												<Item.Media variant="icon">
													<Icon icon="ph:cpu" class="size-4" />
												</Item.Media>
												<Item.Content>
													<Item.Description class="text-xs">
														{image.platform.os}
													</Item.Description>
													{#if image.platform.osVersion}
														<Item.Title class="text-xs">
															{image.platform.osVersion}
														</Item.Title>
													{/if}
												</Item.Content>
											</Item.Root>
										</div>
									{/if}

									{#if image.config}
										<Label>Configuration</Label>
										<div class="grid grid-cols-1 gap-2 md:grid-cols-2">
											<Item.Root class="p-0">
												<Item.Media variant="icon">
													<Icon icon="ph:briefcase" class="size-4" />
												</Item.Media>
												<Item.Content>
													<Item.Description class="text-xs">Working Directory</Item.Description>
													<Item.Title class="text-xs">{image.config.workingDir || '/'}</Item.Title>
												</Item.Content>
											</Item.Root>

											<Item.Root class="p-0">
												<Item.Media variant="icon">
													<Icon icon="ph:plug" class="size-4" />
												</Item.Media>
												<Item.Content class="h-full">
													<Item.Description class="text-xs">Exposed Ports</Item.Description>
													<Item.Title class="text-xs">
														{#each image.config.exposedPorts as port (port)}
															<p>{port}</p>
														{/each}
													</Item.Title>
												</Item.Content>
											</Item.Root>

											<!-- <Item.Root class="p-0">
												<Item.Media>
													<Icon icon="ph:list" class="size-4" />
												</Item.Media>
												<Item.Content>
													<Item.Description class="text-xs">Environments</Item.Description>
													<Item.Title class="text-xs">
														{#each image.config.environments as environment (environment)}
															<p>{environment}</p>
														{/each}
													</Item.Title>
												</Item.Content>
											</Item.Root> -->

											<Item.Root class="col-span-1 p-0 md:col-span-2">
												<Item.Media variant="icon">
													<Icon icon="ph:sign-in" class="size-4" />
												</Item.Media>
												<Item.Content class="h-full">
													<Item.Description class="text-xs">Entrypoint</Item.Description>
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
												<Item.Media variant="icon">
													<Icon icon="ph:terminal" class="size-4" />
												</Item.Media>
												<Item.Content class="h-full">
													<Item.Description class="text-xs">CMD</Item.Description>
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
													<Item.Media variant="icon">
														<Icon icon="ph:folder" class="size-4" />
													</Item.Media>
													<Item.Content class="h-full">
														<Item.Description class="text-xs">Volumes</Item.Description>
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
							</Card.Root>
						{/if}
					</Popover.Content>
				</Popover.Root>
			{/each}
		</Sheet.Content>
	</Sheet.Root>
{/if}
