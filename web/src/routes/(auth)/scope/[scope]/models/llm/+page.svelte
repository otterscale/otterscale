<script lang="ts">
	import { page } from '$app/state';
	import LLM from '$lib/components/models/large-language-model/index.svelte';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb, currentKubernetes } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.models(page.params.scope)],
		current: dynamicPaths.modelsLLM(page.params.scope),
	});
</script>

{#if $currentKubernetes}
	{#key $currentKubernetes?.scope + $currentKubernetes?.name}
		<LLM scopeUuid={$currentKubernetes?.scope} facilityName={$currentKubernetes?.name} />
	{/key}
{/if}
