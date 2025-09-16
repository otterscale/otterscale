<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import { DataTable } from './data-table';

	import type { VirtualMachineDisk } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import * as Sheet from '$lib/components/ui/sheet';
</script>

<script lang="ts">
	let {
		virtualMachines,
	}: {
		virtualMachines: VirtualMachineDisk[];
	} = $props();
</script>

<div class="flex items-center justify-end gap-1">
	{virtualMachines.length}
	<Sheet.Root>
		<Sheet.Trigger>
			<Icon icon="ph:arrow-square-out" />
		</Sheet.Trigger>
		<Sheet.Content class="min-w-[70vw] p-4">
			<DataTable {virtualMachines} />
		</Sheet.Content>
	</Sheet.Root>
</div>
