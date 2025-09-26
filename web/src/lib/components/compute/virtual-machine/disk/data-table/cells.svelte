<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import Actions from './cell-actions.svelte';

	import {
		type DataVolumeSource,
		VirtualMachineDisk_type,
		VirtualMachineDisk_bus,
	} from '$lib/api/kubevirt/v1/kubevirt_pb';
	import type { EnhancedDisk } from '$lib/components/compute/virtual-machine/units/type';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Badge } from '$lib/components/ui/badge';
	import * as Tooltip from '$lib/components/ui/tooltip';

	export const cells = {
		row_picker,
		name,
		type,
		bus,
		source,
		sourceType,
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

{#snippet type(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-end">
		<Badge variant="outline">
			{VirtualMachineDisk_type[row.original.diskType]}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet bus(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-end">
		<Badge variant="outline">
			{VirtualMachineDisk_bus[row.original.busType]}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet sourceType(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-end">
		{row.original.sourceData.case === 'dataVolume' ? (row.original.sourceData.value as DataVolumeSource).type : ''}
	</Layout.Cell>
{/snippet}

{#snippet source(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-end">
		{@const sourceText =
			row.original.sourceData.case === 'source'
				? row.original.sourceData.value
				: row.original.sourceData.case === 'dataVolume'
					? (row.original.sourceData.value as DataVolumeSource).source
					: ''}

		<Tooltip.Root>
			<Tooltip.Trigger>
				<p class="max-w-[210px] truncate">
					{sourceText}
				</p>
			</Tooltip.Trigger>
			<Tooltip.Content>
				{sourceText}
			</Tooltip.Content>
		</Tooltip.Root>
	</Layout.Cell>
{/snippet}

{#snippet size(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-end">
		{row.original.sourceData.case === 'dataVolume'
			? (row.original.sourceData.value as DataVolumeSource).sizeBytes
			: ''}
	</Layout.Cell>
{/snippet}

{#snippet actions(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		<Actions virtualMachineDisk={row.original} />
	</Layout.Cell>
{/snippet}
