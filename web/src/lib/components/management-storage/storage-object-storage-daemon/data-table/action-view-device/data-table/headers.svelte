<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import type { Table, Column } from '@tanstack/table-core';
	import type { Device } from './types';

	export const headers = {
		_row_picker: _row_picker,
		id: id,
		name: name,
		state: state,
		lifeExpectancy: lifeExpectancy,
		daemons: daemons
	};
</script>

{#snippet _row_picker(table: Table<Device>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet id(column: Column<Device>)}
	<div class="flex items-center justify-between gap-1">
		<p class="text-xs font-light">ID</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet name(column: Column<Device>)}
	<div class="flex items-center justify-between gap-1">
		<p class="text-xs font-light">NAME</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet state(column: Column<Device>)}
	<div class="flex items-center justify-between gap-1">
		<p class="text-xs font-light">STATE</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet lifeExpectancy(column: Column<Device>)}
	<div class="flex items-center justify-between gap-1">
		<p class="text-xs font-light">LIFE EXPECTANCY</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet daemons(column: Column<Device>)}
	<div class="flex items-center justify-between gap-1">
		<p class="text-xs font-light">DAEMON</p>
		<Sorter {column} />
	</div>
{/snippet}
