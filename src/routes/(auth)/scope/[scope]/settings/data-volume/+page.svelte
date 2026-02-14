<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { DataVolume } from '$lib/components/settings/data-volume';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/scope/[scope]/settings', { scope: page.params.scope! })
		},
		{
			title: m.data_volume(),
			url: resolve('/(auth)/scope/[scope]/settings/data-volume', { scope: page.params.scope! })
		}
	]);
</script>

{#key page.params.scope!}
	<DataVolume scope={page.params.scope!} />
{/key}
