<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import type { Table, Column } from '@tanstack/table-core';
	import type { BlockImageSnapshot } from './types';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		size: size,
		used: used,
		state: state,
		createTime: createTime
	};
</script>

{#snippet _row_picker(table: Table<BlockImageSnapshot>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet name(column: Column<BlockImageSnapshot>)}
	<div class="flex items-center justify-between gap-1">
		<p class="text-xs font-light">NAME</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet size(column: Column<BlockImageSnapshot>)}
	<div class="flex items-center justify-between gap-1">
		<Sorter {column} />
		<p class="text-xs font-light">SIZE</p>
	</div>
{/snippet}

{#snippet used(column: Column<BlockImageSnapshot>)}
	<div class="flex items-center justify-between gap-1">
		<Sorter {column} />
		<p class="text-xs font-light">USED</p>
	</div>
{/snippet}

{#snippet state(column: Column<BlockImageSnapshot>)}
	<div class="flex items-center justify-between gap-1">
		<p class="text-xs font-light">STATE</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet createTime(column: Column<BlockImageSnapshot>)}
	<div class="flex items-center justify-between gap-1">
		<p class="text-xs font-light">CREATE TIME</p>
		<Sorter {column} />
	</div>
{/snippet}
