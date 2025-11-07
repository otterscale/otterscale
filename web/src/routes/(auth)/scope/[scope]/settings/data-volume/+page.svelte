<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { DataVolume } from '$lib/components/settings/data-volume';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs, currentKubernetes } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/scope/[scope]/settings', { scope: page.params.scope! }),
		},
		{
			title: m.data_volume(),
			url: resolve('/(auth)/scope/[scope]/settings/data-volume', { scope: page.params.scope! }),
		},
	]);
</script>

{#if $currentKubernetes}
	{@const scope = $currentKubernetes.scope}
	{@const facility = $currentKubernetes.name}
	{@const namespace = page.params.namespace ?? ''}

	<DataVolume {scope} {facility} {namespace} />
{/if}
