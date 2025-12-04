<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';

	import type { Bucket } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity, formatTimeAgo } from '$lib/formatter';

	import Actions from './cell-actions.svelte';

	export const cells = {
		row_picker,
		name,
		owner,
		usage,
		createTime,
		actions
	};
</script>

{#snippet row_picker(row: Row<Bucket>)}
	<Table.Cell alignClass="items-center">
		<Cells.RowPicker {row} />
	</Table.Cell>
{/snippet}

{#snippet name(row: Row<Bucket>)}
	<Table.Cell alignClass="items-start">
		{row.original.name}
	</Table.Cell>
{/snippet}

{#snippet owner(row: Row<Bucket>)}
	<Table.Cell alignClass="items-start">
		<Badge variant="outline">{row.original.owner}</Badge>
	</Table.Cell>
{/snippet}

{#snippet usage(row: Row<Bucket>)}
	{@const { value, unit } = formatCapacity(row.original.usedBytes)}
	<Table.Cell alignClass="items-end gap-0">
		{value}
		{unit}
		<Layout.SubCell>{row.original.usedObjects} unit(s)</Layout.SubCell>
	</Table.Cell>
{/snippet}

{#snippet createTime(row: Row<Bucket>)}
	<Table.Cell alignClass="items-start">
		{#if row.original.createdAt}
			{formatTimeAgo(timestampDate(row.original.createdAt))}
		{/if}
	</Table.Cell>
{/snippet}

{#snippet actions(data: { row: Row<Bucket>; scope: string; reloadManager: ReloadManager })}
	<Table.Cell alignClass="items-start">
		<Actions bucket={data.row.original} scope={data.scope} reloadManager={data.reloadManager} />
	</Table.Cell>
{/snippet}
