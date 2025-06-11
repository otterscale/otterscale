<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import type { Table, Column } from '@tanstack/table-core';
	import type { BlockImage } from './types';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		pool: pool,
		namespace: namespace,
		// size: size,
		usage: usage,
		objects: objects,
		parent: parent,
		mirroring: mirroring,
		nextScheduledSnapshot: nextScheduledSnapshot
	};
</script>

{#snippet _row_picker(table: Table<BlockImage>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet name(column: Column<BlockImage>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">NAME</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet pool(column: Column<BlockImage>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">POOL</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet namespace(column: Column<BlockImage>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">NAMESPACE</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet usage(column: Column<BlockImage>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">USAGE</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet objects(column: Column<BlockImage>)}
	<div class="flex items-center gap-1">
		<Sorter {column} />
		<div class="flex flex-col items-end text-xs font-light">
			OBJECTS
			<p class="font-extralight">SIZE</p>
		</div>
	</div>
{/snippet}

{#snippet parent(column: Column<BlockImage>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">PARENT</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet mirroring(column: Column<BlockImage>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">MIRRORING</p>
		<Sorter {column} />
	</div>
{/snippet}

{#snippet nextScheduledSnapshot(column: Column<BlockImage>)}
	<div class="flex items-center gap-1">
		<p class="text-xs font-light">NEXT SCHEDULED SNAPSHOT</p>
		<Sorter {column} />
	</div>
{/snippet}
