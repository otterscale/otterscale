<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Dashboard from '$lib/components/models/dashboard/index.svelte';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.models(),
			url: resolve('/(auth)/scope/[scope]/models', { scope: page.params.scope! })
		}
	]);
</script>

{#key page.params.scope!}
	<Dashboard scope={page.params.scope!} />
{/key}
