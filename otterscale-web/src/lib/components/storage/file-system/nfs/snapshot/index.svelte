<script lang="ts">
	import type { Subvolume } from '$lib/api/storage/v1/storage_pb';
	import * as Sheet from '$lib/components/ui/sheet';
	import Icon from '@iconify/svelte';
	import { setContext } from 'svelte';
	import { DataTable } from './data-table';

	let {
		subvolume
	}: {
		subvolume: Subvolume;
	} = $props();

	setContext('subvolume', subvolume);
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
