<script lang="ts" module>
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Column, Table } from '@tanstack/table-core';
	import type { Snapshot } from './types';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		createTime: createTime,
		pendingClones: pendingClones
	};
</script>

{#snippet _row_picker(table: Table<Snapshot>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet name(column: Column<Snapshot>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">NAME</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet createTime(column: Column<Snapshot>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">CREATE TIME</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet pendingClones(column: Column<Snapshot>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">PENDING CLONE</p>
		<Sorter {column} />
	</div>
{/snippet}
