<script lang="ts" module>
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Column, Table } from '@tanstack/table-core';
	import type { Snapshot } from './types';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		size: size,
		used: used,
		state: state,
		createTime: createTime
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

{#snippet size(column: Column<Snapshot>)}
	<Layout.Header>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>SIZE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet used(column: Column<Snapshot>)}
	<Layout.Header>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>USED</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet state(column: Column<Snapshot>)}
	<Layout.Header>
		<Layout.HeaderViewer>STATE</Layout.HeaderViewer>
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
