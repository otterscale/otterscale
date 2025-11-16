<script lang="ts" module>
	import { ReloadManager } from '$lib/components/custom/reloader';
	import Icon from '@iconify/svelte';

	import type { Network_Subnet } from '$lib/api/network/v1/network_pb';
	import * as Sheet from '$lib/components/ui/sheet';

	import { DataTable } from './data-table';
</script>

<script lang="ts">
	let {
		subnet,
		reloadManager
	}: {
		subnet: Network_Subnet;
		reloadManager: ReloadManager;
	} = $props();
</script>

<div class="flex items-center justify-end gap-1">
	{subnet.ipRanges.length}
	<Sheet.Root>
		<Sheet.Trigger>
			<Icon icon="ph:arrow-square-out" />
		</Sheet.Trigger>
		<Sheet.Content class="min-w-[50vw] p-4">
			<DataTable {subnet} {reloadManager} />
		</Sheet.Content>
	</Sheet.Root>
</div>
