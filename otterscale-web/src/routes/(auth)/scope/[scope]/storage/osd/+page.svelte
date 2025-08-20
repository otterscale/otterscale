<script lang="ts">
	import { page } from '$app/state';
	import { OSD } from '$lib/components/storage/osd';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	let selectedScopeUuid = $derived($activeScope ? $activeScope.uuid : '');
	let selectedFacility = $state('ceph-mon');

	breadcrumb.set({
		parents: [dynamicPaths.storage(page.params.scope)],
		current: dynamicPaths.storageOSD(page.params.scope)
	});
</script>

{#if $activeScope}
	{#key selectedScopeUuid + selectedFacility}
		<OSD bind:selectedScopeUuid bind:selectedFacility />
	{/key}
{/if}
