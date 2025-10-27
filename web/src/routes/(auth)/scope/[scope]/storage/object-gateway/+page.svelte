<script lang="ts" module>
	import { page } from '$app/state';
	import { ObjectDateway } from '$lib/components/storage/object-gateway';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
</script>

<script lang="ts">
	let selectedScope = $derived($activeScope ? $activeScope.name : '');
	let selectedFacility = $state('ceph-mon');

	breadcrumb.set({
		parents: [dynamicPaths.storage(page.params.scope)],
		current: dynamicPaths.storageObjectGateway(page.params.scope),
	});
</script>

{#if $activeScope}
	{#key selectedScope + selectedFacility}
		<ObjectDateway bind:selectedScope bind:selectedFacility />
	{/key}
{/if}
