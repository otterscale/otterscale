<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';

	import Actions from './cell-actions.svelte';

	import type { VirtualMachine_Snapshot } from '$lib/api/instance/v1/instance_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Badge } from '$lib/components/ui/badge';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatTimeAgo } from '$lib/formatter';

	export const cells = {
		row_picker,
		name,
		namespace,
		sourceName,
		phase,
		ready,
		createTime,
		actions,
	};
</script>

{#snippet row_picker(row: Row<VirtualMachine_Snapshot>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<VirtualMachine_Snapshot>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet namespace(row: Row<VirtualMachine_Snapshot>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.namespace}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet sourceName(row: Row<VirtualMachine_Snapshot>)}
	<Layout.Cell class="items-start">
		{row.original.sourceName}
	</Layout.Cell>
{/snippet}

{#snippet phase(row: Row<VirtualMachine_Snapshot>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.phase}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet ready(row: Row<VirtualMachine_Snapshot>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.readyToUse}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet createTime(row: Row<VirtualMachine_Snapshot>)}
	<Layout.Cell class="items-start">
		{#if row.original.createdAt}
			<Tooltip.Provider>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{formatTimeAgo(timestampDate(row.original.createdAt))}
					</Tooltip.Trigger>
					<Tooltip.Content>
						{timestampDate(row.original.createdAt)}
					</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet actions(row: Row<VirtualMachine_Snapshot>)}
	<Layout.Cell class="items-start">
		<Actions virtualMachineSnapshot={row.original} />
	</Layout.Cell>
{/snippet}
