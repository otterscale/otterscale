<script lang="ts" module>
	import type { Subvolume } from '$gen/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { Writable } from 'svelte/store';
	import Delete from './delete.svelte';
	import Edit from './edit.svelte';
	import Grant from './grant-export-access.svelte';
	import Revoke from './revoke-export-access.svelte';
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		selectedVolume,
		selectedSubvolumeGroup,
		subvolume,
		data = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		selectedVolume: string;
		selectedSubvolumeGroup: string;
		subvolume: Subvolume;
		data: Writable<Subvolume[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Grant
			{selectedScope}
			{selectedFacility}
			{selectedSubvolumeGroup}
			{selectedVolume}
			{subvolume}
			bind:data
		/>
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Revoke
			{selectedScope}
			{selectedFacility}
			{selectedSubvolumeGroup}
			{selectedVolume}
			{subvolume}
			bind:data
		/>
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Edit
			{selectedScope}
			{selectedFacility}
			{selectedSubvolumeGroup}
			{selectedVolume}
			{subvolume}
			bind:data
		/>
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete
			{selectedScope}
			{selectedFacility}
			{selectedSubvolumeGroup}
			{selectedVolume}
			{subvolume}
			bind:data
		/>
	</Layout.ActionItem>
</Layout.Actions>
