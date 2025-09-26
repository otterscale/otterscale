<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import { DataTable } from './data-table';

	import type { VirtualMachine, VirtualMachine_Disk } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import * as Sheet from '$lib/components/ui/sheet';
</script>

<script lang="ts">
	let {
		virtualMachine,
		virtualMachineDisks,
	}: {
		virtualMachine: VirtualMachine;
		virtualMachineDisks: VirtualMachine_Disk[];
	} = $props();
</script>

<div class="flex items-center justify-end gap-1">
	{virtualMachineDisks.length}
	<Sheet.Root>
		<Sheet.Trigger>
			<Icon icon="ph:arrow-square-out" />
		</Sheet.Trigger>
		<Sheet.Content class="min-w-[70vw] p-4">
			<DataTable {virtualMachine} {virtualMachineDisks} />
		</Sheet.Content>
	</Sheet.Root>
</div>
