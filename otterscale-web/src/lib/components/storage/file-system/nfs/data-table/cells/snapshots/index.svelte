<script lang="ts" module>
	import type { Subvolume } from '$lib/api/storage/v1/storage_pb';
	import * as Sheet from '$lib/components/ui/sheet';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { get } from 'svelte/store';
	import type { NFSStore } from '../../../utils.svelte';
	import { DataTable } from './data-table';
	import { createNFSSnapshotStore, type nfsSnapshotStore } from './utils.svelte';
</script>

<script lang="ts">
	let {
		subvolume
	}: {
		subvolume: Subvolume;
	} = $props();

	const nfsStore: NFSStore = getContext('nfsStore');
	const nfsSnapshotStore: nfsSnapshotStore = createNFSSnapshotStore();

	onMount(() => {
		nfsSnapshotStore.selectedScopeUuid.set(get(nfsStore.selectedScopeUuid));
		nfsSnapshotStore.selectedFacilityName.set(get(nfsStore.selectedFacilityName));
		nfsSnapshotStore.selectedVolumeName.set(get(nfsStore.selectedVolumeName));
		nfsSnapshotStore.selectedSubvolumeGroupName.set(get(nfsStore.selectedSubvolumeGroupName));
		nfsSnapshotStore.selectedSubvolumeName.set(subvolume.name);
	});
</script>

<div class="flex items-center justify-end gap-1">
	{subvolume.snapshots.length}
	<Sheet.Root>
		<Sheet.Trigger>
			<Icon icon="ph:arrow-square-out" />
		</Sheet.Trigger>
		<Sheet.Content class="min-w-[50vw] p-4">
			<DataTable {subvolume} />
		</Sheet.Content>
	</Sheet.Root>
</div>
