<script lang="ts" module>
	import type { User_Key } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import { Header as RowPicker } from '$lib/components/custom/data-table/data-table-row-pickers';
	import Sorter from '$lib/components/custom/data-table/data-table-sorter.svelte';
	import type { Column, Table } from '@tanstack/table-core';

	export const messages = {
		accessKey: 'accessKey',
		actions: 'actions'
	};

	export const headers = {
		row_picker,
		accessKey,
		actions
	};
</script>

{#snippet row_picker(table: Table<User_Key>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet accessKey(column: Column<User_Key>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>ACCESS KEY</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet actions(column: Column<User_Key>)}{/snippet}
