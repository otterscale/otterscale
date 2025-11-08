<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Store } from '$lib/components/applications/store';
	import { m } from '$lib/paraglide/messages';
	import { activeScope, breadcrumbs, currentKubernetes } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.applications(),
			url: resolve('/(auth)/scope/[scope]/applications', { scope: page.params.scope! })
		},
		{
			title: m.store(),
			url: resolve('/(auth)/scope/[scope]/applications/store', { scope: page.params.scope! })
		}
	]);
</script>

{#if $activeScope && $currentKubernetes}
	{@const scope = $currentKubernetes.scope}
	{@const facility = $currentKubernetes.name}
	<Store {scope} {facility} />
{/if}
