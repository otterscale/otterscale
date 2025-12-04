<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';

	import type { Subvolume, Subvolume_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { formatTimeAgo } from '$lib/formatter';

	import Actions from './cell-actions.svelte';

	export const cells = {
		row_picker,
		name,
		createTime,
		hasPendingClones,
		actions
	};
</script>

{#snippet row_picker(row: Row<Subvolume_Snapshot>)}
	<Table.Cell alignClass="items-center">
		<Cells.RowPicker {row} />
	</Table.Cell>
{/snippet}

{#snippet name(row: Row<Subvolume_Snapshot>)}
	<Table.Cell alignClass="items-start">
		{row.original.name}
	</Table.Cell>
{/snippet}

{#snippet hasPendingClones(row: Row<Subvolume_Snapshot>)}
	<Table.Cell alignClass="items-start">
		<Badge variant="outline">{row.original.hasPendingClones}</Badge>
	</Table.Cell>
{/snippet}

{#snippet createTime(row: Row<Subvolume_Snapshot>)}
	<Table.Cell alignClass="items-start">
		{#if row.original.createdAt}
			{formatTimeAgo(timestampDate(row.original.createdAt))}
		{/if}
	</Table.Cell>
{/snippet}

{#snippet actions(data: {
	row: Row<Subvolume_Snapshot>;
	subvolume: Subvolume;
	scope: string;
	volume: string;
	group: string;
	reloadManager: ReloadManager;
})}
	<Table.Cell alignClass="items-start">
		<Actions
			snapshot={data.row.original}
			subvolume={data.subvolume}
			scope={data.scope}
			volume={data.volume}
			group={data.group}
			reloadManager={data.reloadManager}
		/>
	</Table.Cell>
{/snippet}
