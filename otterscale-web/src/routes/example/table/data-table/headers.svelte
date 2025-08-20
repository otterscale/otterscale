<script lang="ts" module>
	import type { TableRow } from './type';
	import Sorter from '$lib/components/custom/data-table/data-table-sorter.svelte';
	import { Header as RowPicker } from '$lib/components/custom/data-table/data-table-row-pickers';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		row_picker,
		id,
		name,
		isVerified
	};
</script>

{#snippet row_picker(table: Table<TableRow>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet id(column: Column<TableRow>)}
	<Layout.Header>
		<Layout.HeaderViewer>ID</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<TableRow>)}
	<Layout.Header>
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet isVerified(column: Column<TableRow>)}
	<Layout.Header>
		<Layout.HeaderViewer>VERIFIED</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}
