<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import Actions from './cell-actions.svelte';

	import { type VirtualMachineSnapshot } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Badge } from '$lib/components/ui/badge';

	export const cells = {
		row_picker,
		name,
		namespace,
		description,
		statusPhase,
		lastConditionMessage,
		lastConditionReason,
		actions,
	};
</script>

{#snippet row_picker(row: Row<VirtualMachineSnapshot>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<VirtualMachineSnapshot>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet namespace(row: Row<VirtualMachineSnapshot>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.namespace}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet statusPhase(row: Row<VirtualMachineSnapshot>)}
	<Layout.Cell class="items-start">
		{row.original.statusPhase}
	</Layout.Cell>
{/snippet}

{#snippet description(row: Row<VirtualMachineSnapshot>)}
	<Layout.Cell class="items-end">
		{row.original.description}
	</Layout.Cell>
{/snippet}

{#snippet lastConditionMessage(row: Row<VirtualMachineSnapshot>)}
	<Layout.Cell class="items-end">
		{row.original.lastConditionMessage}
	</Layout.Cell>
{/snippet}

{#snippet lastConditionReason(row: Row<VirtualMachineSnapshot>)}
	<Layout.Cell class="items-end">
		{row.original.lastConditionReason}
	</Layout.Cell>
{/snippet}

{#snippet actions(row: Row<VirtualMachineSnapshot>)}
	<Layout.Cell class="items-start">
		<Actions snapshot={row.original} />
	</Layout.Cell>
{/snippet}
