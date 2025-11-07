<script lang="ts" module>
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { BlockDevice } from '$lib/components/storage/block-device';
	import { m } from '$lib/paraglide/messages';
	import { activeScope, breadcrumbs } from '$lib/stores';
</script>

<script lang="ts">
	let selectedScope = $derived($activeScope ? $activeScope.name : '');
	let selectedFacility = $state('ceph-mon');

	breadcrumbs.set([
		{ title: m.storage(), url: resolve('/(auth)/scope/[scope]/storage', { scope: page.params.scope! }) },
		{
			title: m.block_device(),
			url: resolve('/(auth)/scope/[scope]/storage/block-device', { scope: page.params.scope! }),
		},
	]);
</script>

{#if $activeScope}
	{#key selectedScope + selectedFacility}
		<BlockDevice bind:selectedScope bind:selectedFacility />
	{/key}
{/if}
