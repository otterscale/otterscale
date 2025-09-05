<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import type { TableRow } from './type';

	import { Cells } from '$lib/components/custom/data-table/core';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	export const cells = {
		row_picker,
		id,
		name,
		isVerified,
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
