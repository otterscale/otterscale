<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import type { Table, Column } from '@tanstack/table-core';
	import type { FlexibleIOTest } from './types';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		rwMode: rwMode,
		fileSize: fileSize,
		numberJobs: numberJobs,
		blockSize: blockSize,
		runtime: runtime,
		createTime: createTime,
		modifyTime: modifyTime,
	};
</script>

{#snippet _row_picker(table: Table<FlexibleIOTest>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet name(column: Column<FlexibleIOTest>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">NAME</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet rwMode(column: Column<FlexibleIOTest>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">READ WRITE MODE</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet fileSize(column: Column<FlexibleIOTest>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">FIZE SIZE</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet numberJobs(column: Column<FlexibleIOTest>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">NUMBER JOBS</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet blockSize(column: Column<FlexibleIOTest>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">BLOCK SIZE</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet runtime(column: Column<FlexibleIOTest>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">Runtime</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet createTime(column: Column<FlexibleIOTest>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">createTime</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet modifyTime(column: Column<FlexibleIOTest>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">modifyTime</p>
		<Sorter {column} />
	</div>
{/snippet}
