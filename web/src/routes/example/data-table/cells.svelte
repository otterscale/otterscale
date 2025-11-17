<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import { Cells } from '$lib/components/custom/data-table/core';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	import Actions from './actions.svelte';
	import type { TableRow } from './type';

	export const cells = {
		row_picker,
		id,
		name,
		isVerified,
		actions
	};
</script>

{#snippet row_picker(row: Row<TableRow>)}
	<Cells.RowPicker {row} />
{/snippet}

{#snippet id(row: Row<TableRow>)}
	{row.original.id}
{/snippet}

{#snippet name(row: Row<TableRow>)}
	<Badge variant="outline">{row.original.name}</Badge>
{/snippet}

{#snippet isVerified(row: Row<TableRow>)}
	<Badge variant="outline">{row.original.isVerified}</Badge>
{/snippet}

{#snippet actions(data: { row: Row<TableRow>; reloadManager: ReloadManager })}
	<Actions row={data.row.original} reloadManager={data.reloadManager} />
{/snippet}
