<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import type { Table, Column } from '@tanstack/table-core';
	import { type TestResult } from '$gen/api/bist/v1/bist_pb';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		status: status,
		target: target,
		operation: operation,
		duration: duration,
		objectSize: objectSize, 
		objectCount: objectCount,
		throughputFastest: throughputFastest,
		throughputSlowest: throughputSlowest,
		throughputMedian: throughputMedian,
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

{#snippet target(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">TARGET</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet createdBy(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">CREATOR</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet operation(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">OPERATION</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet duration(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">DURATION</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet objectSize(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">OBJECT SIZE</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet objectCount(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">OBJECT COUNT</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet throughputFastest(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">THROUGHPUT ( FASTEST )</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet throughputSlowest(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">THROUGHPUT ( SLOWEST )</p>
		<Sorter {column} />
	</div>
{/snippet}


{#snippet throughputMedian(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">THROUGHPUT ( MEDIAN )</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet startedAt(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">STARTED</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet completedAt(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">COMPLETED</p>
		<Sorter {column} />
	</div>
{/snippet}
