<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import type { Table, Column } from '@tanstack/table-core';
	import type { Pool } from './types';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		dataProtection: dataProtection,
		applications: applications,
		PGStatus: PGStatus,
		usage: usage
	};
</script>

{#snippet _row_picker(table: Table<Pool>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet name(column: Column<Pool>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">NAME</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet dataProtection(column: Column<Pool>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">DATA PROTECTION</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet applications(column: Column<Pool>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">APPLICATIONS</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet PGStatus(column: Column<Pool>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">PG STATUS</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet usage(column: Column<Pool>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">USAGE</p>
		<Sorter {column} />
	</div>
{/snippet}
