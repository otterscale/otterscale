<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import type { Table, Column } from '@tanstack/table-core';
	import { type TestResult } from '$gen/api/bist/v1/bist_pb';

	export const headers = {
		_row_picker: _row_picker,
		uid: uid,
		name: name,
		status: status,
		createdBy: createdBy,
		startedAt: startedAt, 
		completedAt: completedAt, 
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


{#snippet uid(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">UID</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet name(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">NAME</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet status(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">STATUS</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet createdBy(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">CREATED BY</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet startedAt(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">STARTED AT</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet completedAt(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">COMPLETED AT</p>
		<Sorter {column} />
	</div>
{/snippet}
