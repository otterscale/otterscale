<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import type { Table, Column } from '@tanstack/table-core';
	import type { OSD } from './types';

	export const headers = {
		_row_picker: _row_picker,
		id: id,
		host: host,
		status: status,
		deviceClass: deviceClass,
		pgs: pgs,
		size: size,
		flags: flags,
		usage: usage
	};
</script>

{#snippet _row_picker(table: Table<OSD>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet id(column: Column<OSD>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">ID</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet host(column: Column<OSD>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">HOST</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet status(column: Column<OSD>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">STATUS</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet deviceClass(column: Column<OSD>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">DEVICE CLASS</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet pgs(column: Column<OSD>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">PLACEMENT GROUP</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet size(column: Column<OSD>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">SIZE</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet flags(column: Column<OSD>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">FLAG</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet usage(column: Column<OSD>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">USAGE</p>
		<Sorter {column} />
	</div>
{/snippet}
