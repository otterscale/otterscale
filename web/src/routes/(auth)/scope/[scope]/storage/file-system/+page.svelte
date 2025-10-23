<script lang="ts" module>
	import { page } from '$app/state';
	import { FileSystem } from '$lib/components/storage/file-system/index';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
</script>

<script lang="ts">
	let selectedScope = $derived($activeScope ? $activeScope.name : '');
	let selectedFacility = $state('ceph-mon');
	let selectedVolume = $state('ceph-fs');
	let selectedSubvolumeGroupName = $state('');

	breadcrumb.set({
		parents: [dynamicPaths.storage(page.params.scope)],
		current: dynamicPaths.storageFileSystem(page.params.scope),
	});
</script>

{#if $activeScope}
	{#key selectedScope + selectedFacility + selectedVolume + selectedSubvolumeGroupName}
		<FileSystem bind:selectedScope bind:selectedFacility bind:selectedVolume bind:selectedSubvolumeGroupName />
	{/key}
{/if}
