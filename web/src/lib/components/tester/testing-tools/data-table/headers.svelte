<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import type { Table, Column } from '@tanstack/table-core';
	import { type TestResult } from '$gen/api/bist/v1/bist_pb';

	export const headers = {
		_row_picker: _row_picker,
		type: type,
		name: name,
		input: input,
	};
</script>

{#snippet _row_picker(table: Table<TestResult>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}


{#snippet type(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">TYPE</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet name(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">NAME</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet input(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">INPUT</p>
		<Sorter {column} />
	</div>
{/snippet}
