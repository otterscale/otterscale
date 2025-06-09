<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import type { Table, Column } from '@tanstack/table-core';
	import type { FileSystem } from './types';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		enabled: enabled,
		permission: permission,
		createTime: createTime
	};
</script>

{#snippet _row_picker(table: Table<FileSystem>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet name(column: Column<FileSystem>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">NAME</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet enabled(column: Column<FileSystem>)}
	<p class="text-end text-xs font-light">ENABLED</p>
{/snippet}

{#snippet permission(column: Column<FileSystem>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">PERMISSION</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet createTime(column: Column<FileSystem>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">CREATE TIME</p>
		<Sorter {column} />
	</div>
{/snippet}
