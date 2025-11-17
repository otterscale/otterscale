<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import type { Image } from '$lib/api/storage/v1/storage_pb';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import * as Sheet from '$lib/components/ui/sheet';

	import { DataTable } from './data-table';
</script>

<script lang="ts">
	let {
		image,
		scope,
		reloadManager
	}: {
		image: Image;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();
</script>

<div class="flex items-center justify-end gap-1">
	{image.snapshots.length}
	<Sheet.Root>
		<Sheet.Trigger>
			<Icon icon="ph:arrow-square-out" />
		</Sheet.Trigger>
		<Sheet.Content class="min-w-[38vw] p-4">
			<DataTable {image} {scope} {reloadManager} />
		</Sheet.Content>
	</Sheet.Root>
</div>
