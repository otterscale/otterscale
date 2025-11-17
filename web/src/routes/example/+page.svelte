<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import type { ModelTag } from '$lib/components/settings/model-artifact/types';
	import { fetchModelTypes } from '$lib/components/settings/model-artifact/utils.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
</script>

<script lang="ts">
	const pipelineTags = writable([] as ModelTag[]);
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

	const licenseTags = writable([] as ModelTag[]);
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

	onMount(async () => {
		await fetchPipelineTags();
		await fetchLicenseTags();
	});
</script>

<main class="flex items-center gap-1">
	{#if isPipelineTagsLoaded}
		<div class="flex max-h-[500px] flex-col gap-4 overflow-y-auto border p-4">
			{#each $pipelineTags as tag}
				<div class="flex items-center gap-2">
					<Icon icon="ph:cylinder" class="size-5" />
					<div class="flex flex-col gap-1">
						<h6 class="text-sm">{tag.label}</h6>
						<span class="flex items-center gap-1">
							<Badge variant="outline">
								{tag.type}
							</Badge>
							{#if tag.subType}
								<Badge variant="outline">
									{tag.subType}
								</Badge>
							{/if}
						</span>
					</div>
				</div>
			{/each}
		</div>
	{/if}

	{#if isLicenseTagsLoaded}
		<div class="flex max-h-[500px] flex-col gap-4 overflow-y-auto border p-4">
			{#each $licenseTags as tag}
				<div class="flex items-center gap-2">
					<Icon icon="ph:identification-badge" class="size-5" />
					<div class="flex flex-col gap-1">
						<h6 class="text-sm">{tag.label}</h6>
						<span class="flex items-center gap-1">
							<Badge variant="outline">
								{tag.type}
							</Badge>
							{#if tag.subType}
								<Badge variant="outline">
									{tag.subType}
								</Badge>
							{/if}
						</span>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</main>
