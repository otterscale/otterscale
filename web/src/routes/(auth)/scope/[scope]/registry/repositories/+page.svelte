<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Repositories from '$lib/components/registry/repositories/index.svelte';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.models(),
			url: resolve('/(auth)/scope/[scope]/registry', { scope: page.params.scope! })
		},
		{
			title: m.llm(),
			url: resolve('/(auth)/scope/[scope]/registry/repositories', { scope: page.params.scope! })
		}
	]);
</script>

{#key page.params.scope!}
	<Repositories scope={page.params.scope!} />
{/key}
