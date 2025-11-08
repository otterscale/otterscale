<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Metal } from '$lib/components/machines/metal/index';
	import { m } from '$lib/paraglide/messages';
	import { activeScope, breadcrumbs } from '$lib/stores';

	breadcrumbs.set([
		{ title: m.machines(), url: resolve('/(auth)/machines') },
		{ title: m.metal(), url: resolve('/(auth)/machines/metal') },
		{
			title: page.params.id!,
			url: resolve('/(auth)/machines/metal/[id]', {
				id: page.params.id!
			})
		}
	]);
</script>

{#if $activeScope}
	<Metal />
{/if}
