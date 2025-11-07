<script lang="ts">
	import { resolve } from '$app/paths';
	import { Metals } from '$lib/components/machines/metals/index';
	import { m } from '$lib/paraglide/messages';
	import { activeScope, breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{ title: m.machines(), url: resolve('/(auth)/machines') },
		{
			title: m.metal(),
			url: resolve('/(auth)/machines/metal'),
		},
	]);
</script>

{#if $activeScope}
	{#key $activeScope.name}
		<Metals />
	{/key}
{/if}
