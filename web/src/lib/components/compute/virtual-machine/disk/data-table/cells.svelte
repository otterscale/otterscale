<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	// import Actions from './cell-actions.svelte';

	import { type VirtualMachineDisk, VirtualMachineDisk_type } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Badge } from '$lib/components/ui/badge';

	export const cells = {
		row_picker,
		name,
		type,
		bus,
		source,
		// actions,
	};
</script>

{#snippet row_picker(row: Row<VirtualMachineDisk>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<VirtualMachineDisk>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet type(row: Row<VirtualMachineDisk>)}
	<Layout.Cell class="items-end">
		<Badge variant="outline">
			{VirtualMachineDisk_type[row.original.diskType]}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet bus(row: Row<VirtualMachineDisk>)}
	<Layout.Cell class="items-end">
		<Badge variant="outline">
			{row.original.bus}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet source(row: Row<VirtualMachineDisk>)}
	<Layout.Cell class="items-end">
		{row.original.source}
	</Layout.Cell>
{/snippet}

<!-- {#snippet actions(row: Row<VirtualMachineDisk>)}
	<Layout.Cell class="items-start">
		<Actions snapshot={row.original} />
	</Layout.Cell>
{/snippet} -->
