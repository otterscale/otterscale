<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import type { Table, Column } from '@tanstack/table-core';
	import type { Namespace } from './types';

	export const headers = {
		_row_picker: _row_picker,
		pool: pool,
		namespace: namespace,
		totalImages: totalImages
	};
</script>

{#snippet _row_picker(table: Table<Namespace>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet pool(column: Column<Namespace>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">POOL</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet namespace(column: Column<Namespace>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">NAMESPACE</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet totalImages(column: Column<Namespace>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">TOTAL IMAGE</p>
		<Sorter {column} />
	</div>
{/snippet}
