<script lang="ts" module>
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { ObjectDateway } from '$lib/components/storage/object-gateway';
	import { m } from '$lib/paraglide/messages';
	import { activeScope, breadcrumbs } from '$lib/stores';
</script>

<script lang="ts">
	let selectedScope = $derived($activeScope ? $activeScope.name : '');
	let selectedFacility = $state('ceph-mon');

	breadcrumbs.set([
		{
			title: m.storage(),
			url: resolve('/(auth)/scope/[scope]/storage', { scope: page.params.scope! })
		},
		{
			title: m.object_gateway(),
			url: resolve('/(auth)/scope/[scope]/storage/object-gateway', { scope: page.params.scope! })
		}
	]);
</script>

{#if $activeScope}
	{#key selectedScope + selectedFacility}
		<ObjectDateway bind:selectedScope bind:selectedFacility />
	{/key}
{/if}
