<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Extensions } from '$lib/components/settings/extensions';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs, currentKubernetes } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/scope/[scope]/scope-based-settings', { scope: page.params.scope! }),
		},
		{
			title: m.extensions(),
			url: resolve('/(auth)/scope/[scope]/scope-based-settings/extensions', { scope: page.params.scope! }),
		},
	]);
</script>

{#if $currentKubernetes}
	{@const scope = $currentKubernetes.scope}
	{@const facility = $currentKubernetes.name}

	<Extensions {scope} {facility} />
{/if}
