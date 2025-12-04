<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';

	import type { VirtualMachine_Restore } from '$lib/api/instance/v1/instance_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { ReloadManager } from '$lib/components/custom/reloader';
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
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<VirtualMachine_Restore>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet namespace(row: Row<VirtualMachine_Restore>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.namespace}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet targetName(row: Row<VirtualMachine_Restore>)}
	<Layout.Cell class="items-start">
		{row.original.targetName}
	</Layout.Cell>
{/snippet}

{#snippet snapshotName(row: Row<VirtualMachine_Restore>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.snapshotName}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet complete(row: Row<VirtualMachine_Restore>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.complete}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet createTime(row: Row<VirtualMachine_Restore>)}
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

{#snippet actions(data: {
	row: Row<VirtualMachine_Restore>;
	scope: string;
	reloadManager: ReloadManager;
})}
	<Layout.Cell class="items-start">
		<Actions
			virtualMachineRestore={data.row.original}
			scope={data.scope}
			reloadManager={data.reloadManager}
		/>
	</Layout.Cell>
{/snippet}
