<script lang="ts">
	import { page } from '$app/state';
	import { ModelArtifact } from '$lib/components/settings/model-artifact';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb, currentKubernetes } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.scopeBasedSettings(page.params.scope)],
		current: { title: m.model(), url: '' },
	});
</script>

{#if $currentKubernetes}
	{@const scope = $currentKubernetes.scope}
	{@const facility = $currentKubernetes.name}
	{@const namespace = page.params.namespace ?? ''}

	<ModelArtifact {scope} {facility} {namespace} />
{/if}
