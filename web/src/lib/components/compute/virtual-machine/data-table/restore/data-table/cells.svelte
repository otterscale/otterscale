<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';

	import type { VirtualMachine_Restore } from '$lib/api/instance/v1/instance_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatTimeAgo } from '$lib/formatter';

	import Actions from './cell-actions.svelte';

	export const cells = {
		row_picker,
		name,
		namespace,
		targetName,
		snapshotName,
		complete,
		createTime,
		actions
	};
</script>

{#snippet row_picker(row: Row<VirtualMachine_Restore>)}
	<Table.Cell alignClass="items-center">
		<Cells.RowPicker {row} />
	</Table.Cell>
{/snippet}

{#snippet name(row: Row<VirtualMachine_Restore>)}
	<Table.Cell alignClass="items-start">
		{row.original.name}
	</Table.Cell>
{/snippet}

{#snippet namespace(row: Row<VirtualMachine_Restore>)}
	<Table.Cell alignClass="items-start">
		<Badge variant="outline">
			{row.original.namespace}
		</Badge>
	</Table.Cell>
{/snippet}

{#snippet targetName(row: Row<VirtualMachine_Restore>)}
	<Table.Cell alignClass="items-start">
		{row.original.targetName}
	</Table.Cell>
{/snippet}

{#snippet snapshotName(row: Row<VirtualMachine_Restore>)}
	<Table.Cell alignClass="items-start">
		<Badge variant="outline">
			{row.original.snapshotName}
		</Badge>
	</Table.Cell>
{/snippet}

{#snippet complete(row: Row<VirtualMachine_Restore>)}
	<Table.Cell alignClass="items-start">
		<Badge variant="outline">
			{row.original.complete}
		</Badge>
	</Table.Cell>
{/snippet}

{#snippet createTime(row: Row<VirtualMachine_Restore>)}
	<Table.Cell alignClass="items-start">
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
	</Table.Cell>
{/snippet}

{#snippet actions(data: {
	row: Row<VirtualMachine_Restore>;
	scope: string;
	reloadManager: ReloadManager;
})}
	<Table.Cell alignClass="items-start">
		<Actions
			virtualMachineRestore={data.row.original}
			scope={data.scope}
			reloadManager={data.reloadManager}
		/>
	</Table.Cell>
{/snippet}
