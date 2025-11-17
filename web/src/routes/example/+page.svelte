<script lang="ts" module>
	import Models from './models.svelte';

	import type { HuggingFaceModel } from '$lib/components/settings/model-artifact/types';
</script>

<script lang="ts">
	let selectedModel: HuggingFaceModel | null = $state(null);
</script>

<main class="flex flex-col gap-2">
	{selectedModel?.id}
	<Models bind:value={selectedModel} />
</main>
