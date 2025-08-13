<script lang="ts" module>
	import type { User_Key } from '$lib/api/storage/v1/storage_pb';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/header.svelte';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		_row_picker: _row_picker,
		access: access
	};
</script>

{#snippet _row_picker(table: Table<User_Key>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<TableRowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet access(column: Column<User_Key>)}
	<Layout.Header>
		<Layout.HeaderViewer>ACCESSOR</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}
