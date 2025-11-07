<script lang="ts" module>
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Pool } from '$lib/components/storage/pool';
	import { m } from '$lib/paraglide/messages';
	import { activeScope, breadcrumbs } from '$lib/stores';
</script>

<script lang="ts">
	let selectedScope = $derived($activeScope ? $activeScope.name : '');
	let selectedFacility = $state('ceph-mon');

	breadcrumbs.set([
		{ title: m.storage(), url: resolve('/(auth)/scope/[scope]/storage', { scope: page.params.scope! }) },
		{ title: m.pool(), url: resolve('/(auth)/scope/[scope]/storage/pool', { scope: page.params.scope! }) },
	]);
</script>

{#if $activeScope}
	{#key selectedScope + selectedFacility}
		<Pool bind:selectedScope bind:selectedFacility />
	{/key}
{/if}
