<script lang="ts">
	import Icon from '@iconify/svelte';
	import { onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { Single as SingleInput } from '$lib/components/custom/input';
	import {
		type HuggingFaceModel,
		type ModelTag,
		type SortType
	} from '$lib/components/settings/model-artifact/types';
	import {
		fetchModels,
		fetchModelTypes
	} from '$lib/components/settings/model-artifact/utils.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Command from '$lib/components/ui/command';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Popover from '$lib/components/ui/popover';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { formatBigNumber, formatTimeAgo } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	import type { ModeSource } from '../types';

	let {
		value = $bindable(),
		modelSource = $bindable()
	}: { value: string; modelSource?: ModeSource } = $props();

	let selectedModel: HuggingFaceModel | null = $state(null);

	const defaultLimit = 30;
	let limit = $state(defaultLimit);

	const defaultSort = 'downloads' as SortType;
	let sort = $state<SortType>(defaultSort);

	let huggingFaceModels = $state([] as HuggingFaceModel[]);
	let isHuggingFaceModelsLoaded = $state(false);
	async function fetchModelsByTags(tags: ModelTag[]) {
		const response = await fetchModels(
			'RedHatAI',
			tags.map((tag) => tag.id),
			sort,
			limit
		);
		huggingFaceModels = response;
		isHuggingFaceModelsLoaded = true;
	}

	const defaultLicenseTags = [] as ModelTag[];
	const licenseTags = writable(defaultLicenseTags);
	let selectedLicenseTags = $state(defaultLicenseTags);
	let isLicenseTagsLoaded = $state(false);
	async function fetchLicenseTags() {
		const response = await fetchModelTypes('license');
		licenseTags.set(response);
		isLicenseTagsLoaded = true;
	}
	function handleLicenseTagSelect(selectedLicenseTag: ModelTag) {
		if (selectedLicenseTags.find((licenseTag) => licenseTag.id === selectedLicenseTag.id)) {
			selectedLicenseTags = selectedLicenseTags.filter(
				(licenseTag) => licenseTag.id !== selectedLicenseTag.id
			);
		} else {
			selectedLicenseTags = [...selectedLicenseTags, selectedLicenseTag];
		}
	}

	const defaultLibraryTags = [] as ModelTag[];
	const libraryTags = writable(defaultLibraryTags);
	let selectedLibraryTags = $state(defaultLibraryTags);
	let isLibraryTagsLoaded = $state(false);
	async function fetchLibraryTags() {
		const response = await fetchModelTypes('library');
		libraryTags.set(response);
		isLibraryTagsLoaded = true;
	}
	function handleLibraryTagSelect(selectedLibraryTag: ModelTag) {
		if (selectedLibraryTags.find((libraryTag) => libraryTag.id === selectedLibraryTag.id)) {
			selectedLibraryTags = selectedLibraryTags.filter(
				(libraryTag) => libraryTag.id !== selectedLibraryTag.id
			);
		} else {
			selectedLibraryTags = [...selectedLibraryTags, selectedLibraryTag];
		}
	}

	const defaultPipelineTags = [] as ModelTag[];
	const pipelineTags = writable(defaultPipelineTags);
	let selectedPipelineTags = $state(defaultPipelineTags);
	let isPipelineTagsLoaded = $state(false);
	async function fetchPipelineTags() {
		const response = await fetchModelTypes('pipeline_tag');
		pipelineTags.set(response);
		isPipelineTagsLoaded = true;
	}
	function handlePipelineTagSelect(selectedPipelineTag: ModelTag) {
		if (selectedPipelineTags.find((pipelineTag) => pipelineTag.id === selectedPipelineTag.id)) {
			selectedPipelineTags = selectedPipelineTags.filter(
				(pipelineTag) => pipelineTag.id !== selectedPipelineTag.id
			);
		} else {
			selectedPipelineTags = [...selectedPipelineTags, selectedPipelineTag];
		}
	}

	const modelTags = $derived([
		...selectedLicenseTags,
		...selectedLibraryTags,
		...selectedPipelineTags
	]);

	function isTagSelected(selectedModelTag: ModelTag) {
		return modelTags.find((modelTag) => modelTag.id === selectedModelTag.id);
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	function reset() {
		selectedLicenseTags = defaultLicenseTags;
		selectedLibraryTags = defaultLibraryTags;
		selectedPipelineTags = defaultPipelineTags;
		limit = defaultLimit;
		sort = defaultSort;
	}

	$effect(() => {
		if (modelTags && limit && sort) {
			fetchModelsByTags(modelTags);
		}
	});
	onMount(async () => {
		await fetchPipelineTags();
		await fetchLicenseTags();
		await fetchLibraryTags();
		await fetchModelsByTags([]);
	});
	onDestroy(() => {
		reset();
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger
		class={cn('w-fit', buttonVariants({ variant: 'outline' }))}
		onclick={() => {
			modelSource = 'cloud' as ModeSource;
		}}
	>
		<div class="font-base flex items-center gap-2 text-sm text-primary">
			<Icon icon="ph:cloud" />
			cloud
		</div>
	</Dialog.Trigger>
	<Dialog.Content
		class="h-[77vh] min-w-[77vw] overflow-y-auto"
		data-slot="huggingface-models-store"
	>
		<div class="h-full space-y-8">
			<Dialog.Header>
				<Dialog.Title>{m.models_store_title()}</Dialog.Title>
				<Dialog.Description>
					{m.models_store_description()}
				</Dialog.Description>
			</Dialog.Header>
			{#if selectedModel}
				<!-- Selected Model -->
				<div class="flex w-full items-center justify-between rounded-lg bg-muted p-4">
					<div class="flex w-full items-center gap-2 text-center">
						<span class="rounded-full bg-muted-foreground/50 p-2 text-card">
							<Icon icon="ph:robot" class="size-8" />
						</span>
						<div class="space-y-1">
							<p class="text-sm">{selectedModel.id}</p>
							<div class="hidden items-center text-xs text-muted-foreground md:flex">
								<span class="flex items-center gap-1">
									<Icon icon="ph:clock" />
									<p>{formatTimeAgo(new Date(selectedModel.createdAt))}</p>
								</span>
								<Icon icon="ph:dot-bold" />
								<span class="flex items-center gap-1">
									<Icon icon="ph:download-simple" />
									<p>{formatBigNumber(selectedModel.downloads)}</p>
								</span>
								<Icon icon="ph:dot-bold" />
								<span class="flex items-center gap-1">
									<Icon icon="ph:heart" />
									<p>{formatBigNumber(selectedModel.likes)}</p>
								</span>
							</div>
						</div>
					</div>
					<Button
						class="ml-auto h-8"
						onclick={() => {
							if (selectedModel) {
								value = selectedModel.id;
							}
							close();
						}}
					>
						{m.confirm()}
					</Button>
				</div>
			{/if}
			<div class="flex h-12 items-center gap-1">
				{#if isPipelineTagsLoaded}
					<!-- Pipeline Filter -->
					<Popover.Root>
						<Popover.Trigger>
							<Button variant="outline" size="sm" class="flex h-8 items-center gap-2">
								<Icon icon="ph:cylinder" class="size-4" />
								{m.pipeline()}
								{#if selectedPipelineTags.length > 0}
									<Separator orientation="vertical" class="h-4" />
									{selectedPipelineTags.length}
									<Icon icon="ph:tag" class="size-4" />
								{/if}
								<Icon icon="ph:caret-down" class="size-4" />
							</Button>
						</Popover.Trigger>
						<Popover.Content class="p-0">
							<Command.Root>
								<Command.Input placeholder="Search" />
								<Command.List>
									<Command.Empty>{m.no_result()}</Command.Empty>
									<Command.Group>
										{#each $pipelineTags as tag (tag)}
											<Command.Item
												onclick={() => {
													handlePipelineTagSelect(tag);
												}}
											>
												<Icon
													icon="ph:check"
													class={cn(isTagSelected(tag) ? 'visible' : 'invisible', 'h-4 w-4')}
												/>
												<div class="flex flex-col items-start justify-start gap-1">
													<h6 class="text-sm">{tag.label}</h6>
													{#if tag.subType}
														<span class="flex items-center gap-1">
															<Icon icon="ph:tag" class="size-4" />
															<p class="text-xs text-muted-foreground">
																{tag.subType}
															</p>
														</span>
													{/if}
												</div>
											</Command.Item>
										{/each}
									</Command.Group>
								</Command.List>
							</Command.Root>
						</Popover.Content>
					</Popover.Root>
				{/if}
				{#if isLibraryTagsLoaded}
					<!-- Library Filter -->
					<Popover.Root>
						<Popover.Trigger>
							<Button variant="outline" size="sm" class="flex h-8 items-center gap-2">
								<Icon icon="ph:package" class="size-4" />
								{m.library()}
								{#if selectedLibraryTags.length > 0}
									<Separator orientation="vertical" class="h-4" />
									{selectedLibraryTags.length}
									<Icon icon="ph:tag" class="size-4" />
								{/if}
								<Icon icon="ph:caret-down" class="size-4" />
							</Button>
						</Popover.Trigger>
						<Popover.Content class="p-0">
							<Command.Root>
								<Command.Input placeholder="Search" />
								<Command.List>
									<Command.Empty>{m.no_result()}</Command.Empty>
									<Command.Group>
										{#each $libraryTags as tag (tag)}
											<Command.Item
												onclick={() => {
													handleLibraryTagSelect(tag);
												}}
											>
												<Icon
													icon="ph:check"
													class={cn(isTagSelected(tag) ? 'visible' : 'invisible', 'h-4 w-4')}
												/>
												<div class="flex flex-col items-start justify-start gap-1">
													<h6 class="text-sm">{tag.label}</h6>
													{#if tag.subType}
														<span class="flex items-center gap-1">
															<Icon icon="ph:tag" class="size-4" />
															<p class="text-xs text-muted-foreground">
																{tag.subType}
															</p>
														</span>
													{/if}
												</div>
											</Command.Item>
										{/each}
									</Command.Group>
								</Command.List>
							</Command.Root>
						</Popover.Content>
					</Popover.Root>
				{/if}
				{#if isLicenseTagsLoaded}
					<!-- License Filter -->
					<Popover.Root>
						<Popover.Trigger>
							<Button variant="outline" size="sm" class="flex h-8 items-center gap-2">
								<Icon icon="ph:identification-badge" class="size-4" />
								{m.license()}
								{#if selectedLicenseTags.length > 0}
									<Separator orientation="vertical" class="h-4" />
									{selectedLicenseTags.length}
									<Icon icon="ph:tag" class="size-4" />
								{/if}
								<Icon icon="ph:caret-down" class="size-4" />
							</Button>
						</Popover.Trigger>
						<Popover.Content class="p-0">
							<Command.Root>
								<Command.Input placeholder="Search" />
								<Command.List>
									<Command.Empty>{m.no_result()}</Command.Empty>
									<Command.Group>
										{#each $licenseTags as tag (tag)}
											<Command.Item
												onclick={() => {
													handleLicenseTagSelect(tag);
												}}
											>
												<Icon
													icon="ph:check"
													class={cn(isTagSelected(tag) ? 'visible' : 'invisible', 'h-4 w-4')}
												/>
												<div class="flex flex-col items-start justify-start gap-1">
													<h6 class="text-sm">{tag.label}</h6>
													{#if tag.subType}
														<span class="flex items-center gap-1">
															<Icon icon="ph:tag" class="size-4" />
															<p class="text-xs text-muted-foreground">
																{tag.subType}
															</p>
														</span>
													{/if}
												</div>
											</Command.Item>
										{/each}
									</Command.Group>
								</Command.List>
							</Command.Root>
						</Popover.Content>
					</Popover.Root>
				{/if}
				<!-- Sort -->
				<Popover.Root>
					<Popover.Trigger>
						<Button variant="outline" size="sm" class="flex h-8 items-center gap-2">
							<Icon icon="ph:funnel" class="size-4" />
							<p class="capitalize">{sort}</p>
							<Icon icon="ph:caret-down" class="size-4" />
						</Button>
					</Popover.Trigger>
					<Popover.Content class="p-0">
						<Command.Root>
							<Command.List>
								<Command.Group>
									<Command.Item
										onclick={() => {
											sort = 'downloads';
										}}
									>
										<Icon icon={sort === 'downloads' ? 'ph:download-fill' : 'ph:download'} />
										<p>{m.downloads()}</p>
									</Command.Item>
									<Command.Item
										onclick={() => {
											sort = 'likes';
										}}
									>
										<Icon icon={sort === 'likes' ? 'ph:heart-fill' : 'ph:heart'} />
										<p>{m.likes()}</p>
									</Command.Item>
								</Command.Group>
							</Command.List>
						</Command.Root>
					</Popover.Content>
				</Popover.Root>
				<SingleInput.General type="number" class="h-8 w-24" bind:value={limit} min={0} step={6} />
			</div>
			{#if huggingFaceModels.length > 0}
				<!-- Models -->
				<div class="grid gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
					{#each huggingFaceModels as model (model.id)}
						<Card.Root
							class={cn(
								'h-40 hover:bg-muted/50',
								selectedModel?.id === model.id ? 'bg-muted' : 'bg-transparent'
							)}
							onclick={() => {
								selectedModel = model;
								document
									.querySelector('[data-slot="huggingface-models-store"]')
									?.scrollTo({ top: 0, behavior: 'smooth' });
							}}
						>
							<Card.Content>
								<div class="flex flex-col gap-2">
									<div class="flex items-center gap-2">
										<span class="h-fit w-fit rounded-full bg-muted p-3">
											<Icon icon="ph:robot" class="size-8" />
										</span>
										<div class="space-y-2">
											<p class="text-sm">{model.id}</p>
											<div class="ml-auto flex items-center text-xs text-muted-foreground">
												<span class="flex items-center gap-1">
													<Icon icon="ph:clock" />
													<p>{formatTimeAgo(new Date(model.createdAt))}</p>
												</span>
												<Icon icon="ph:dot-bold" />
												<span class="flex items-center gap-1">
													<Icon icon="ph:download-simple" />
													<p>{formatBigNumber(model.downloads)}</p>
												</span>
												<Icon icon="ph:dot-bold" />
												<span class="flex items-center gap-1">
													<Icon icon="ph:heart" />
													<p>{formatBigNumber(model.likes)}</p>
												</span>
											</div>
										</div>
									</div>
								</div>
							</Card.Content>
							<Card.Footer class="mt-auto">
								<span class="flex gap-1">
									{#each model.tags.slice(0, 3) as tag (tag)}
										<Badge>{tag}</Badge>
									{/each}
									{#if model.tags.length > 3}
										<Badge variant="outline">+{model.tags.length - 3}</Badge>
									{/if}
								</span>
							</Card.Footer>
						</Card.Root>
					{/each}
				</div>
			{:else if isHuggingFaceModelsLoaded}
				<!-- Empty -->
				<div class="space-y-8 p-8">
					<div class="text-center">
						<h3 class="text-lg font-semibold">{m.no_models_found()}</h3>
						<p class="text-sm text-muted-foreground">
							{m.no_models_matching_filters()}
						</p>
					</div>
					<div class="flex flex-col items-center justify-center gap-4">
						<Icon icon="ph:robot-fill" class="size-32 animate-pulse text-muted-foreground/50" />
						<Button
							variant="destructive"
							onclick={() => {
								reset();
							}}
						>
							{m.reset()}
						</Button>
					</div>
				</div>
			{/if}
		</div>
	</Dialog.Content>
</Dialog.Root>
