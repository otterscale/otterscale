<script lang="ts" module>
	import type { Pool } from '$lib/api/storage/v1/storage_pb';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/header.svelte';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		type: type,
		applications: applications,
		placement_group_state: placement_group_state,
		usage: usage,
		iops: iops
	};
</script>

{#snippet _row_picker(table: Table<Pool>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<TableRowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Pool>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet type(column: Column<Pool>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>TYPE</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet applications(column: Column<Pool>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>APPLICATIONS</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet placement_group_state(column: Column<Pool>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>PLACEMENT GROUP STATE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet usage(column: Column<Pool>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>USAGE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet iops()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>IOPS</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
