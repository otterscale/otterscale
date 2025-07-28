<script lang="ts" module>
	import type { TableRow } from './type';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker,
		id,
		name,
		isVerified
	};
</script>

{#snippet _row_picker(row: Row<TableRow>)}
	<TableRowPicker {row} />
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
