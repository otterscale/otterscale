<script lang="ts" module>
	import type { Subvolume, Subvolume_Snapshot } from '$gen/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { Writable } from 'svelte/store';
	import Delete from './delete.svelte';
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		selectedVolume,
		selectedSubvolumeGroup,
		subvolume,
		snapshot,
		subvolumes: data = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		selectedVolume: string;
		selectedSubvolumeGroup: string;
		subvolume: Subvolume;
		snapshot: Subvolume_Snapshot;
		subvolumes: Writable<Subvolume[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Delete
			{selectedScope}
			{selectedFacility}
			{selectedVolume}
			{selectedSubvolumeGroup}
			{subvolume}
			{snapshot}
			bind:data
		/>
	</Layout.ActionItem>
</Layout.Actions>
