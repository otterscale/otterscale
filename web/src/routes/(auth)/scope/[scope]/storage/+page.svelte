<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Dashboard from '$lib/components/storage/dashboard/index.svelte';
	import { m } from '$lib/paraglide/messages';
	import { activeScope, breadcrumbs, currentCeph } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.storage(),
			url: resolve('/(auth)/scope/[scope]/storage', { scope: page.params.scope! })
		}
	]);
</script>

{#if $activeScope && $currentCeph}
	<Dashboard scope={$activeScope} ceph={$currentCeph} />
{/if}
