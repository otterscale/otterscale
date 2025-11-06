<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { onMount } from 'svelte';

	import { formatBigNumber, formatTimeAgo } from '$lib/formatter';

	type HuggingFaceModel = {
		id: string;
		tags: string[];
		downloads: number;
		likes: number;
		createdAt: string;
	};
</script>

<script lang="ts">
	let huggingFaceModels = $state([] as HuggingFaceModel[]);
	let isLoading = $state(true);

	async function fetchModels() {
		const response = await fetch(
			'https://huggingface.co/api/models?author=RedHatAI&sort=downloads&direction=-1&limit=10',
		);

		if (!response.ok) throw new Error(`Failed to fetch models: ${response.status} ${response.statusText}`);

		huggingFaceModels = await response.json();

		isLoading = false;
	}

	onMount(() => {
		fetchModels();
	});
</script>

{#if !isLoading}
	{#each huggingFaceModels as model}
		<div class="my-2 font-mono">
			<h6 class="text-base">{model.id}</h6>
			<div class="text-muted-foreground flex items-center text-xs">
				<span class="flex items-center gap-1">
					<Icon icon="ph:clock" />
					<p>{formatTimeAgo(new Date(model.createdAt))}</p>
				</span>
				<Icon icon="ph:dot-bold" class="size-6" />
				<span class="flex items-center gap-1">
					<Icon icon="ph:download-simple" />
					<p>{formatBigNumber(model.downloads)}</p>
				</span>
				<Icon icon="ph:dot-bold" class="size-6" />
				<span class="flex items-center gap-1">
					<Icon icon="ph:heart" />
					<p>{formatBigNumber(model.likes)}</p>
				</span>
			</div>
		</div>
	{/each}
{/if}
