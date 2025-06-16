<script lang="ts" module>
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Column, Table } from '@tanstack/table-core';
	import type { Subvolume } from './types';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		dataPool: dataPool,
		usage: usage,
		path: path,
		mode: mode,
		createTime: createTime
	};
</script>

{#snippet _row_picker(table: Table<Subvolume>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet name(column: Column<Subvolume>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">NAME</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet dataPool(column: Column<Subvolume>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">DATA POOL</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet usage(column: Column<Subvolume>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">USAGE</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet path(column: Column<Subvolume>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">PATH</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet mode(column: Column<Subvolume>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">MODE</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet createTime(column: Column<Subvolume>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">CREATE TIME</p>
		<Sorter {column} />
	</div>
{/snippet}
