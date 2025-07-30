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
		accessMode: accessMode,
		jobCount: jobCount,
		runTime: runTime,
		blockSize: blockSize,
		fileSize: fileSize,
		ioDepth: ioDepth,
		bandwidth: bandwidth,
		iops: iops,
		latencyMin: latencyMin,
		latencyMax: latencyMax,
		latencyMean: latencyMean,
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

{#snippet accessMode(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">ACCESS MODE</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet jobCount(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">JOB COUNT</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet runTime(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">RUN TIME</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet blockSize(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">BLOCK SIZE</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet fileSize(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">FILE SIZE</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet ioDepth(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">IO DEPTH</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet createdBy(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">CREATOR</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet bandwidth(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">BANDWIDTH</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet iops(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">IOPS</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet latencyMin(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">LATENCY ( MIN )</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet latencyMax(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">LATENCY ( MAX )</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet latencyMean(column: Column<TestResult>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">LATENCY ( MEAN )</p>
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
