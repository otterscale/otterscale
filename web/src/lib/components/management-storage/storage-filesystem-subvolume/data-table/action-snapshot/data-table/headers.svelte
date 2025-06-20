<script lang="ts" module>
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Column, Table } from '@tanstack/table-core';
	import type { Snapshot } from './types';
	import * as Layout from '$lib/components/custom/data-table/layout';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		createTime: createTime,
		pendingClones: pendingClones
	};
</script>

{#snippet _row_picker(table: Table<Snapshot>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}
{#snippet name(column: Column<Snapshot>)}
	<Layout.Header>
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet createTime(column: Column<Snapshot>)}
	<Layout.Header>
		<Layout.HeaderViewer>CREATE TIME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet pendingClones(column: Column<Snapshot>)}
	<Layout.Header>
		<Layout.HeaderViewer>PENDING CLONE</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}
