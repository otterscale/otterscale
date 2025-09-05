<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';

	import Actions from './cell-actions.svelte';

	import type { Subvolume_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Badge } from '$lib/components/ui/badge';
	import { formatTimeAgo } from '$lib/formatter';

	export const cells = {
		row_picker,
		name,
		createTime,
		hasPendingClones,
		actions,
	};
</script>

{#snippet row_picker(row: Row<Subvolume_Snapshot>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<Subvolume_Snapshot>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet hasPendingClones(row: Row<Subvolume_Snapshot>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">{row.original.hasPendingClones}</Badge>
	</Layout.Cell>
{/snippet}

{#snippet createTime(row: Row<Subvolume_Snapshot>)}
	<Layout.Cell class="items-start">
		{#if row.original.createdAt}
			{formatTimeAgo(timestampDate(row.original.createdAt))}
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet actions(row: Row<Subvolume_Snapshot>)}
	<Layout.Cell class="items-start">
		<Actions snapshot={row.original} />
	</Layout.Cell>
{/snippet}
