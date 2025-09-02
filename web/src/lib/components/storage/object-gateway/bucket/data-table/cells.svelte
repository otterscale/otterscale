<script lang="ts" module>
	import type { Bucket } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity, formatTimeAgo } from '$lib/formatter';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { Row } from '@tanstack/table-core';
	import Actions from './cell-actions.svelte';

	export const cells = {
		row_picker,
		name,
		owner,
		usage,
		createTime,
		actions,
	};
</script>

{#snippet row_picker(row: Row<Bucket>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<Bucket>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet owner(row: Row<Bucket>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">{row.original.owner}</Badge>
	</Layout.Cell>
{/snippet}

{#snippet usage(row: Row<Bucket>)}
	{@const { value, unit } = formatCapacity(row.original.usedBytes)}
	<Layout.Cell class="items-end">
		{value}
		{unit}
		<Layout.SubCell>{row.original.usedObjects} unit(s)</Layout.SubCell>
	</Layout.Cell>
{/snippet}

{#snippet createTime(row: Row<Bucket>)}
	<Layout.Cell class="items-start">
		{#if row.original.createdAt}
			{formatTimeAgo(timestampDate(row.original.createdAt))}
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet actions(row: Row<Bucket>)}
	<Layout.Cell class="items-start">
		<Actions bucket={row.original} />
	</Layout.Cell>
{/snippet}
