<script lang="ts" module>
	import type { User } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/header.svelte';
	import Sorter from '$lib/components/custom/data-table/data-table-sorter.svelte';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		_row_picker,
		id,
		name,
		suspended,
		keys,
		actions
	};
</script>

{#snippet _row_picker(table: Table<User>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<TableRowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet id(column: Column<User>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>ID</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<User>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet suspended(column: Column<User>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>SUSPENDED</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet keys(column: Column<User>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>KEYS</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions(column: Column<User>)}
	<Layout.Header>
		<Layout.HeaderViewer>ACTIONS</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
