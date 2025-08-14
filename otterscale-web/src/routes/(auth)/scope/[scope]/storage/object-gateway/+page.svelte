<script lang="ts" module>
	import { page } from '$app/state';
	import { ObjectDateway } from '$lib/components/storage/object-gateway';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
</script>

<script lang="ts">
	let selectedScopeUuid = $derived($activeScope ? $activeScope.uuid : '');
	let selectedFacility = $state('ceph-mon');

	breadcrumb.set({
		parents: [dynamicPaths.storage(page.params.scope)],
		current: dynamicPaths.storagePool(page.params.scope)
	});
</script>

{#if $activeScope}
	{#key selectedScopeUuid + selectedFacility}
		<ObjectDateway bind:selectedScopeUuid bind:selectedFacility />
	{/key}
{/if}
