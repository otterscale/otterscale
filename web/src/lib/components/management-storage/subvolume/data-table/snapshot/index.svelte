<script lang="ts" module>
	import type { Subvolume } from '$gen/api/storage/v1/storage_pb';
	import * as Sheet from '$lib/components/ui/sheet';
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';
	import { DataTable } from './data-table';
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		selectedVolume,
		selectedSubvolumeGroup,
		subvolume,
		subvolumes = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		selectedVolume: string;
		selectedSubvolumeGroup: string;
		subvolume: Subvolume;
		subvolumes: Writable<Subvolume[]>;
	} = $props();
</script>

<div class="flex items-center justify-end gap-1">
	{subvolume.snapshots.length}
	<Sheet.Root>
		<Sheet.Trigger>
			<Icon icon="ph:arrow-square-out" />
		</Sheet.Trigger>
		<Sheet.Content class="min-w-[50vw]">
			<DataTable
				{selectedScope}
				{selectedFacility}
				{selectedVolume}
				{selectedSubvolumeGroup}
				{subvolume}
				bind:subvolumes
			/>
		</Sheet.Content>
	</Sheet.Root>
</div>
