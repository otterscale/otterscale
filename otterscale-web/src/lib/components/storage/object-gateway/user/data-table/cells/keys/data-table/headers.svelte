<script lang="ts" module>
	import type { User_Key } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/header.svelte';
	import Sorter from '$lib/components/custom/data-table/data-table-sorter.svelte';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		_row_picker,
		access,
		actions
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
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>ACCESSOR</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet actions(column: Column<User_Key>)}{/snippet}
