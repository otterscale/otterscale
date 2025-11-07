<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Metal } from '$lib/components/machines/metal/index';
	import { m } from '$lib/paraglide/messages';
	import { activeScope, breadcrumbs } from '$lib/stores';

	breadcrumbs.set([
		{ title: m.machines(), url: resolve('/(auth)/scope/[scope]/machines', { scope: page.params.scope! }) },
		{ title: m.metal(), url: resolve('/(auth)/scope/[scope]/machines/metal', { scope: page.params.scope! }) },
		{
			title: page.params.id!,
			url: resolve('/(auth)/scope/[scope]/machines/metal/[id]', {
				scope: page.params.scope!,
				id: page.params.id!,
			}),
		},
	]);
</script>

{#if $activeScope}
	<Metal />
{/if}
