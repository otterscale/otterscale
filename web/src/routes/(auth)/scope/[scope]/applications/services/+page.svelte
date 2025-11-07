<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Services } from '$lib/components/applications/services';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs, currentKubernetes } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.applications(),
			url: resolve('/(auth)/scope/[scope]/applications', { scope: page.params.scope! }),
		},
		{
			title: m.services(),
			url: resolve('/(auth)/scope/[scope]/applications/services', { scope: page.params.scope! }),
		},
	]);
</script>

{#if $currentKubernetes}
	{@const scope = $currentKubernetes.scope}
	{@const facility = $currentKubernetes.name}

	<Services {scope} {facility} />
{/if}
