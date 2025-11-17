<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { Single as SingleInput } from '$lib/components/custom/input';
	import { type HuggingFaceModel, type ModelTag, type SortType } from '$lib/components/settings/model-artifact/types';
	import { fetchModels, fetchModelTypes } from '$lib/components/settings/model-artifact/utils.svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Command from '$lib/components/ui/command';
	import * as Popover from '$lib/components/ui/popover';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { formatBigNumber, formatTimeAgo } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let { value = $bindable() }: { value: HuggingFaceModel | null } = $props();

	let selectedModel: HuggingFaceModel | null = $state(null);

	let huggingFaceModels = $state([] as HuggingFaceModel[]);
	let isModelsLoaded = $state(false);
	async function fetchModelsByTags(tags: ModelTag[]) {
		isModelsLoaded = false;
		await fetchModels(
			'RedHatAI',
			tags.map((tag) => tag.id),
			sort,
			limit,
		)
			.then((response) => {
				huggingFaceModels = response;
				isModelsLoaded = true;
			})
			.catch((error) => {
				console.error('Error fetching models on create modal mount:', error);
			});
	}

	const defaultLicenseTags = [] as ModelTag[];
	const licenseTags = writable(defaultLicenseTags);
	let selectedLicenseTags = $state(defaultLicenseTags);
	let isLicenseTagsLoaded = $state(false);
	async function fetchLicenseTags() {
		fetchModelTypes('license')
			.then((response) => {
				licenseTags.set(response);
				console.log(licenseTags);
			})
			.catch((error) => {
				console.error('Error fetching model tags:', error);
			});
		isLicenseTagsLoaded = true;
	}
	function handleLicenseTagSelect(modelTag: ModelTag) {
		if (selectedLicenseTags.find((t) => t.id === modelTag.id)) {
			selectedLicenseTags = selectedLicenseTags.filter((t) => t.id !== modelTag.id);
		} else {
			selectedLicenseTags = [...selectedLicenseTags, modelTag];
		}
	}

	const defaultModelTags = [] as ModelTag[];
	const libraryTags = writable(defaultModelTags);
	let selectedLibraryTags = $state(defaultModelTags);
	let isLibraryTagsLoaded = $state(false);
	async function fetchLibraryTags() {
		fetchModelTypes('library')
			.then((response) => {
				libraryTags.set(response);
				console.log(libraryTags);
			})
			.catch((error) => {
				console.error('Error fetching model tags:', error);
			});
		isLibraryTagsLoaded = true;
	}
	function handleLibraryTagSelect(modelTag: ModelTag) {
		if (selectedLibraryTags.find((t) => t.id === modelTag.id)) {
			selectedLibraryTags = selectedLibraryTags.filter((t) => t.id !== modelTag.id);
		} else {
			selectedLibraryTags = [...selectedLibraryTags, modelTag];
		}
	}

	const defaultPipelineTags = [] as ModelTag[];
	const pipelineTags = writable(defaultPipelineTags);
	let selectedPipelineTags = $state(defaultPipelineTags);
	let isPipelineTagsLoaded = $state(false);
	async function fetchPipelineTags() {
		fetchModelTypes('pipeline_tag')
			.then((response) => {
				pipelineTags.set(response);
				console.log(pipelineTags);
			})
			.catch((error) => {
				console.error('Error fetching model tags:', error);
			});
		isPipelineTagsLoaded = true;
	}
	function handlePipelineTagSelect(modelTag: ModelTag) {
		if (selectedPipelineTags.find((t) => t.id === modelTag.id)) {
			selectedPipelineTags = selectedPipelineTags.filter((t) => t.id !== modelTag.id);
		} else {
			selectedPipelineTags = [...selectedPipelineTags, modelTag];
		}
	}

	const modelTags = $derived([...selectedLicenseTags, ...selectedLibraryTags, ...selectedPipelineTags]);

	const defaultLimit = 30;
	let limit = $state(defaultLimit);

	const defaultSort = 'downloads' as SortType;
	let sort = $state<SortType>(defaultSort);

	function isTagSelected(modelTag: ModelTag) {
		return modelTags.find((t) => t.id === modelTag.id);
	}

	function reset() {
		selectedLicenseTags = defaultLicenseTags;
		selectedLibraryTags = defaultModelTags;
		selectedPipelineTags = defaultPipelineTags;
		limit = defaultLimit;
		sort = defaultSort;
	}

	$effect(() => {
		if (modelTags && limit && sort) {
			fetchModelsByTags(modelTags);
		}
	});

	let open = $state(false);
	function close() {
		open = false;
	}

	onMount(async () => {
		await fetchPipelineTags();
		await fetchLicenseTags();
		await fetchLibraryTags();
		await fetchModelsByTags([]);
	});
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class={buttonVariants({ variant: 'outline' })}>
		<Icon icon="ph:plus" />
	</AlertDialog.Trigger>
	<AlertDialog.Content class="max-h-[50vh] min-h-[77vh] min-w-[77vw] space-y-8 overflow-y-auto p-8">
		<div class="flex max-h-12 items-start gap-1">
			{#if isPipelineTagsLoaded}
				<Popover.Root>
					<Popover.Trigger>
						<Button variant="outline" size="sm" class="flex h-8 items-center gap-2">
							<Icon icon="ph:cylinder" class="size-4" />
							Pipeline
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
									{#each $pipelineTags as tag}
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
														<p class="text-muted-foreground text-xs">
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
				<Popover.Root>
					<Popover.Trigger>
						<Button variant="outline" size="sm" class="flex h-8 items-center gap-2">
							<Icon icon="ph:package" class="size-4" />
							Library
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
									{#each $libraryTags as tag}
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
														<p class="text-muted-foreground text-xs">
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
				<Popover.Root>
					<Popover.Trigger>
						<Button variant="outline" size="sm" class="flex h-8 items-center gap-2">
							<Icon icon="ph:identification-badge" class="size-4" />
							License
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
									{#each $licenseTags as tag}
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
														<p class="text-muted-foreground text-xs">
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
									<p>Downloads</p>
								</Command.Item>
								<Command.Item
									onclick={() => {
										sort = 'likes';
									}}
								>
									<Icon icon={sort === 'likes' ? 'ph:heart-fill' : 'ph:heart'} />
									<p>Likes</p>
								</Command.Item>
							</Command.Group>
						</Command.List>
					</Command.Root>
				</Popover.Content>
			</Popover.Root>
			<SingleInput.General type="number" class="h-8 w-24" bind:value={limit} min={0} step={6} />
			{#if selectedModel}
				<Button
					class="ml-auto h-8"
					onclick={() => {
						value = selectedModel;
						close();
					}}
				>
					Confirm
				</Button>
			{/if}
		</div>
		{#if selectedModel}
			<div
				class="text-muted-foreground bg-muted flex h-40 flex-col items-center justify-center gap-4 rounded-lg text-center text-sm"
			>
				<Icon icon="ph:robot" class="size-24" />
				{selectedModel.id}
			</div>
		{/if}
		{#if huggingFaceModels.length === 0 && isModelsLoaded}
			<div
				class="text-muted-foreground flex h-full flex-col items-center justify-center gap-4 rounded-lg bg-red-50 text-center text-sm"
			>
				<Icon icon="ph:robot-fill" class="size-60 animate-pulse" />
				<p>There is no model matching the selected filters.</p>
				<Button
					variant="destructive"
					onclick={() => {
						reset();
					}}
				>
					Reset
				</Button>
			</div>
		{:else}
			<div class="grid gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
				{#each huggingFaceModels as model}
					<Card.Root
						class={cn('hover:bg-muted/50', selectedModel?.id === model.id ? 'bg-muted' : 'bg-transparent')}
						onclick={() => {
							selectedModel = model;
						}}
					>
						<Card.Content>
							<div class="flex items-center gap-2">
								<span class="bg-muted h-fit w-fit rounded-full p-3">
									<Icon icon="ph:robot" class="size-8" />
								</span>
								<div class="space-y-2">
									<p class="text-sm">{model.id}</p>
									<span class="flex gap-1">
										{#each model.tags.slice(0, 3) as tag}
											<Badge>{tag}</Badge>
										{/each}
										{#if model.tags.length > 3}
											<Badge variant="outline">+{model.tags.length - 3}</Badge>
										{/if}
									</span>
									<div class="text-muted-foreground ml-auto flex items-center text-xs">
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
						</Card.Content>
					</Card.Root>
				{/each}
			</div>
		{/if}
	</AlertDialog.Content>
</AlertDialog.Root>
