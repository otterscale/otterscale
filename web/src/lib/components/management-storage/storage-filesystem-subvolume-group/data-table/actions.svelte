<script lang="ts" module>
	import type { SubvolumeGroup } from '$gen/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { Writable } from 'svelte/store';
	import ActionDelete from './delete.svelte';
	import ActionEdit from './edit.svelte';
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		selectedVolume,
		subvolumeGroup,
		data = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		selectedVolume: string;
		subvolumeGroup: SubvolumeGroup;
		data: Writable<SubvolumeGroup[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<ActionEdit {selectedScope} {selectedFacility} {selectedVolume} {subvolumeGroup} bind:data />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionDelete {selectedScope} {selectedFacility} {selectedVolume} {subvolumeGroup} bind:data />
	</Layout.ActionItem>
</Layout.Actions>
