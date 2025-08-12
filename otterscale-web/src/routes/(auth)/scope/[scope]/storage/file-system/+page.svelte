<script lang="ts" module>
	import { page } from '$app/state';
	import { NFS } from '$lib/components/storage/nfs';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
</script>

<script lang="ts">
	let selectedScopeUuid = $state($activeScope ? $activeScope.uuid : '');
	let selectedFacility = $state('ceph-mon');
	let selectedVolume = $state('ceph-fs');
	let selectedSubvolumeGroup = $state('');

	breadcrumb.set({
		parents: [dynamicPaths.storage(page.params.scope)],
		current: dynamicPaths.storagePool(page.params.scope)
	});
</script>

{#if $activeScope}
	{#key selectedScopeUuid + selectedFacility + selectedVolume + selectedSubvolumeGroup}
		<NFS
			bind:selectedScopeUuid
			bind:selectedFacility
			bind:selectedVolume
			bind:selectedSubvolumeGroup
		/>
	{/key}
{/if}
