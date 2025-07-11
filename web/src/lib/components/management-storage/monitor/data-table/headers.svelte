<script lang="ts" module>
	import type { MON } from '$gen/api/storage/v1/storage_pb';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/header.svelte';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		_row_picker: _row_picker,
		leader: leader,
		name: name,
		rank: rank,
		publicAddress: publicAddress,
		machine: machine
	};
</script>

{#snippet _row_picker(table: Table<MON>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<TableRowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet leader(column: Column<MON>)}
	<Layout.Header>
		<Layout.HeaderViewer>LEADER</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<MON>)}
	<Layout.Header>
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet rank(column: Column<MON>)}
	<Layout.Header>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>RANK</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet publicAddress(column: Column<MON>)}
	<Layout.Header>
		<Layout.HeaderViewer>PUBLIC ADDRESS</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet machine(column: Column<MON>)}
	<Layout.Header>
		<Layout.HeaderViewer>MACHINE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
