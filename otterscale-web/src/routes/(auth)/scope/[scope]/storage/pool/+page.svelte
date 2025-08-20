<script lang="ts">
	import { page } from '$app/state';
	import { Pool } from '$lib/components/storage/pool';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	let selectedScopeUuid = $derived($activeScope ? $activeScope.uuid : '');
	let selectedFacility = $state('ceph-mon');

	breadcrumb.set({
		parents: [dynamicPaths.storage(page.params.scope)],
		current: dynamicPaths.storagePool(page.params.scope)
	});
</script>

{#if $activeScope}
	{#key selectedScopeUuid + selectedFacility}
		<Pool bind:selectedScopeUuid bind:selectedFacility />
	{/key}
{/if}
