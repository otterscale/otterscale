<script lang="ts">
	import { page } from '$app/state';
	import { Plugins } from '$lib/components/settings/plugins';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb, currentKubernetes } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.settings(page.params.scope)],
		current: { title: m.package_repository(), url: '' },
	});
</script>

{#if $currentKubernetes}
	{@const scope = $currentKubernetes.scope}
	{@const facility = $currentKubernetes.name}

	<Plugins {scope} {facility} />
{/if}
