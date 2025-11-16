<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import type { Subvolume } from '$lib/api/storage/v1/storage_pb';
	import * as Sheet from '$lib/components/ui/sheet';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table';
</script>

<script lang="ts">
	let {
		subvolume,
		scope,
		volume,
		group,
		reloadManager
	}: {
		subvolume: Subvolume;
		scope: string;
		volume: string;
		group: string;
		reloadManager: ReloadManager;
	} = $props();
</script>

<div class="flex items-center justify-end gap-1">
	{subvolume.snapshots.length}
	<Sheet.Root>
		<Sheet.Trigger>
			<Icon icon="ph:arrow-square-out" />
		</Sheet.Trigger>
		<Sheet.Content class="min-w-[50vw] p-4">
			<DataTable {subvolume} {scope} {volume} {group} {reloadManager}  />
		</Sheet.Content>
	</Sheet.Root>
</div>
