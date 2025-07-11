<script lang="ts" module>
	import type { User } from '$gen/api/storage/v1/storage_pb';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/header.svelte';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		_row_picker: _row_picker,
		id: id,
		name: name,
		suspended: suspended,
		keys: keys
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
	<Layout.Header>
		<Layout.HeaderViewer>ID</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<User>)}
	<Layout.Header>
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet suspended(column: Column<User>)}
	<Layout.Header>
		<Layout.HeaderViewer>SUSPENDED</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet keys()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>KEYS</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
