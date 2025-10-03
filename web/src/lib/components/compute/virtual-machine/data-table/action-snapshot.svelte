<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import { DataTable } from './snapshot/data-table';

	import type { VirtualMachine } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import * as Sheet from '$lib/components/ui/sheet';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { virtualMachine }: { virtualMachine: VirtualMachine } = $props();
</script>

<div class="flex items-center justify-end gap-1">
	<Sheet.Root>
		<Sheet.Trigger class="flex items-center gap-1">
			<Icon icon="ph:camera" />
			{m.snapshot()}
		</Sheet.Trigger>

		<Sheet.Content class="min-w-[70vw] p-4">
			<DataTable {virtualMachine} />
		</Sheet.Content>
	</Sheet.Root>
</div>
