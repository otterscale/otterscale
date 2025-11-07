<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Metals } from '$lib/components/machines/metals/index';
	import { m } from '$lib/paraglide/messages';
	import { activeScope, breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{ title: m.machines(), url: resolve('/(auth)/scope/[scope]/machines', { scope: page.params.scope! }) },
		{
			title: m.metal(),
			url: resolve('/(auth)/scope/[scope]/machines/metal', { scope: page.params.scope! }),
		},
	]);
</script>

{#if $activeScope}
	{#key $activeScope.name}
		<Metals />
	{/key}
{/if}
