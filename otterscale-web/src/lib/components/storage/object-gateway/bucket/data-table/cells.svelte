<script lang="ts" module>
	import type { Bucket } from '$lib/api/storage/v1/storage_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity, formatTimeAgo } from '$lib/formatter';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		owner: owner,
		usage: usage,
		createTime: createTime
	};
</script>

{#snippet _row_picker(row: Row<Bucket>)}
	<TableRowPicker {row} />
{/snippet}

{#snippet name(row: Row<Bucket>)}
	{row.original.name}
{/snippet}

{#snippet owner(row: Row<Bucket>)}
	<Badge variant="outline">{row.original.owner}</Badge>
{/snippet}

{#snippet usage(row: Row<Bucket>)}
	{@const { value, unit } = formatCapacity(Number(row.original.usedBytes) / (1024 * 1024))}
	<div class="flex flex-col items-end">
		<div class="flex items-end">
			{value}
			{unit}
		</div>
		<p class="font-light">{row.original.usedObjects} unit(s)</p>
	</div>
{/snippet}

{#snippet createTime(row: Row<Bucket>)}
	{#if row.original.createdAt}
		{formatTimeAgo(timestampDate(row.original.createdAt))}
	{/if}
{/snippet}
