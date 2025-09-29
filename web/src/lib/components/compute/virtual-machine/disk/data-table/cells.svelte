<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import Actions from './cell-actions.svelte';

	import { VirtualMachine_Disk_Bus } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import type { EnhancedDisk } from '$lib/components/compute/virtual-machine/units/type';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity } from '$lib/formatter';

	export const cells = {
		row_picker,
		name,
		bus,
		bootOrder,
		dataVolume,
		type,
		phase,
		boot,
		size,
		actions,
	};
</script>

{#snippet row_picker(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet bus(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{VirtualMachine_Disk_Bus[row.original.bus]}
		</Badge>
	</Layout.Cell>
{/snippet}

, , , boot,

{#snippet bootOrder(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.bootOrder}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet dataVolume(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		{row.original.volume?.name ?? ''}
	</Layout.Cell>
{/snippet}

{#snippet type(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		{row.original.volume?.source?.type ?? ''}
	</Layout.Cell>
{/snippet}

{#snippet phase(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		{row.original.phase}
	</Layout.Cell>
{/snippet}

{#snippet boot(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.bootImage}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet size(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-end">
		{#if row.original.sizeBytes}
			{@const { value, unit } = formatCapacity(row.original.sizeBytes)}
			{value}
			{unit}
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet actions(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		<Actions virtualMachineDisk={row.original} />
	</Layout.Cell>
{/snippet}
