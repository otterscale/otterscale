<script lang="ts">
	import { page } from '$app/state';
	import Dashboard from '$lib/components/models/dashboard/index.svelte';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb, currentKubernetes } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [],
		current: dynamicPaths.models(page.params.scope),
	});
</script>

{#if $currentKubernetes}
	{#key $currentKubernetes?.scope + $currentKubernetes?.name}
		<Dashboard scope={$currentKubernetes?.scope} facility={$currentKubernetes?.name} />
	{/key}
{/if}
