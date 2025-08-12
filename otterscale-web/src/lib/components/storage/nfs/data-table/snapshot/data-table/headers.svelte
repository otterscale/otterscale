<script lang="ts" module>
	import type { Subvolume_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/header.svelte';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		createTime: createTime,
		hasPendingClones: hasPendingClones
	};
</script>

{#snippet _row_picker(table: Table<Subvolume_Snapshot>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<TableRowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Subvolume_Snapshot>)}
	<Layout.Header>
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet hasPendingClones(column: Column<Subvolume_Snapshot>)}
	<Layout.Header>
		<Layout.HeaderViewer>PENDING CLONE</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet createTime(column: Column<Subvolume_Snapshot>)}
	<Layout.Header>
		<Layout.HeaderViewer>CREATE TIME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}
