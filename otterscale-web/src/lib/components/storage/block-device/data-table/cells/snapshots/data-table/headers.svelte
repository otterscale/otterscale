<script lang="ts" module>
	import type { Image_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/header.svelte';
	import Sorter from '$lib/components/custom/data-table/data-table-sorter.svelte';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		_row_picker,
		name,
		protect,
		usage,
		actions
	};
</script>

{#snippet _row_picker(table: Table<Image_Snapshot>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<TableRowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Image_Snapshot>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet protect(column: Column<Image_Snapshot>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>PROTECTED</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet usage(column: Column<Image_Snapshot>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>USAGE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions(column: Column<Image_Snapshot>)}{/snippet}
