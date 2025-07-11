<script lang="ts" module>
	import type { Volume } from '$gen/api/storage/v1/storage_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import { formatTimeAgo } from '$lib/formatter';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		createTime: createTime
	};
</script>

{#snippet _row_picker(row: Row<Volume>)}
	<TableRowPicker {row} />
{/snippet}

{#snippet name(row: Row<Volume>)}
	{row.original.name}
{/snippet}

{#snippet createTime(row: Row<Volume>)}
	{#if row.original.createdAt}
		{formatTimeAgo(timestampDate(row.original.createdAt))}
	{/if}
{/snippet}
