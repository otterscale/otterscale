<script lang="ts" module>
	import type { Subvolume } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import type { Writable } from 'svelte/store';
	import Delete from './delete.svelte';
	import Edit from './edit.svelte';
	import Grant from './grant-export-access.svelte';
	import Revoke from './revoke-export-access.svelte';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		selectedFacility,
		selectedVolume,
		selectedSubvolumeGroup,
		subvolume,
		subvolumes: data = $bindable()
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		selectedVolume: string;
		selectedSubvolumeGroup: string;
		subvolume: Subvolume;
		subvolumes: Writable<Subvolume[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Grant
			{selectedScopeUuid}
			{selectedFacility}
			{selectedSubvolumeGroup}
			{selectedVolume}
			{subvolume}
			bind:data
		/>
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Revoke
			{selectedScopeUuid}
			{selectedFacility}
			{selectedSubvolumeGroup}
			{selectedVolume}
			{subvolume}
			bind:data
		/>
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Edit
			{selectedScopeUuid}
			{selectedFacility}
			{selectedSubvolumeGroup}
			{selectedVolume}
			{subvolume}
			bind:data
		/>
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete
			{selectedScopeUuid}
			{selectedFacility}
			{selectedSubvolumeGroup}
			{selectedVolume}
			{subvolume}
			bind:data
		/>
	</Layout.ActionItem>
</Layout.Actions>
