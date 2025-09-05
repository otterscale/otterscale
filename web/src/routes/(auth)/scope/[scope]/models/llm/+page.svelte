<script lang="ts">
	import { page } from '$app/state';
	import LLM from '$lib/components/models/large-language-model/index.svelte';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.models(page.params.scope)],
		current: dynamicPaths.modelsLLM(page.params.scope),
	});
</script>

<LLM />
