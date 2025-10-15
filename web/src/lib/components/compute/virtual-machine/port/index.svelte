<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import { DataTable } from './data-table';

	import type { VirtualMachine } from '$lib/api/instance/v1/instance_pb';
	import * as Sheet from '$lib/components/ui/sheet';
</script>

<script lang="ts">
	let {
		virtualMachine,
	}: {
		virtualMachine: VirtualMachine;
	} = $props();
</script>

<div class="flex items-center justify-end gap-1">
	{virtualMachine.services.reduce((total, service) => total + service.ports.length, 0)}
	<Sheet.Root>
		<Sheet.Trigger>
			<Icon icon="ph:arrow-square-out" />
		</Sheet.Trigger>
		<Sheet.Content class="min-w-[38vw] p-4">
			<DataTable {virtualMachine} />
		</Sheet.Content>
	</Sheet.Root>
</div>
