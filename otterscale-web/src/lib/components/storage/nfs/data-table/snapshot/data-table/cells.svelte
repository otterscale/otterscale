<script lang="ts" module>
	import type { Subvolume_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { formatTimeAgo } from '$lib/formatter';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		createTime: createTime,
		hasPendingClones: hasPendingClones
	};
</script>

{#snippet _row_picker(row: Row<Subvolume_Snapshot>)}
	<TableRowPicker {row} />
{/snippet}

{#snippet name(row: Row<Subvolume_Snapshot>)}
	{row.original.name}
{/snippet}

{#snippet hasPendingClones(row: Row<Subvolume_Snapshot>)}
	<Badge variant="outline">{row.original.hasPendingClones}</Badge>
{/snippet}

{#snippet createTime(row: Row<Subvolume_Snapshot>)}
	{#if row.original.createdAt}
		{formatTimeAgo(timestampDate(row.original.createdAt))}
	{/if}
{/snippet}
