<script lang="ts" module>
	import type { Subvolume_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import { Badge } from '$lib/components/ui/badge';
	import { formatTimeAgo } from '$lib/formatter';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';
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
	<Cells.RowPicker {row} />
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

{#snippet actions(row: Row<Subvolume_Snapshot>)}
	<Actions snapshot={row.original} />
{/snippet}
