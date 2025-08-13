<script lang="ts" module>
	import type { SubvolumeGroup } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import type { Writable } from 'svelte/store';
	import Delete from './delete.svelte';
	import Edit from './edit.svelte';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		selectedFacility,
		selectedVolume,
		subvolumeGroup,
		subvolumeGroups: data = $bindable()
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		selectedVolume: string;
		subvolumeGroup: SubvolumeGroup;
		subvolumeGroups: Writable<SubvolumeGroup[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {selectedScopeUuid} {selectedFacility} {selectedVolume} {subvolumeGroup} bind:data />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {selectedScopeUuid} {selectedFacility} {selectedVolume} {subvolumeGroup} bind:data />
	</Layout.ActionItem>
</Layout.Actions>
