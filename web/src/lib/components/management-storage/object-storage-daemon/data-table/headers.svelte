<script lang="ts" module>
	import type { OSD } from '$gen/api/storage/v1/storage_pb';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/header.svelte';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		_row_picker: _row_picker,
		id: id,
		name: name,
		stateUp: stateUp,
		stateIn: stateIn,
		exists: exists,
		deviceClass: deviceClass,
		machine: machine,
		placementGroupCount: placementGroupCount,
		usage: usage,
		iopsRead: iopsRead,
		iopsWrite: iopsWrite
	};
</script>

{#snippet _row_picker(table: Table<OSD>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<TableRowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet id(column: Column<OSD>)}
	<Layout.Header>
		<Layout.HeaderViewer>ID</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<OSD>)}
	<Layout.Header>
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet stateIn(column: Column<OSD>)}
	<Layout.Header>
		<Layout.HeaderViewer>IN</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet stateUp(column: Column<OSD>)}
	<Layout.Header>
		<Layout.HeaderViewer>UP</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet exists(column: Column<OSD>)}
	<Layout.Header>
		<Layout.HeaderViewer>EXISTS</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet deviceClass(column: Column<OSD>)}
	<Layout.Header>
		<Layout.HeaderViewer>DEVICE</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet machine(column: Column<OSD>)}
	<Layout.Header>
		<Layout.HeaderViewer>MACHINE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet placementGroupCount(column: Column<OSD>)}
	<Layout.Header>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>PLACEMENT GROUP</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet usage(column: Column<OSD>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>USAGE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet iopsRead()}
	<Layout.Header class="justify-center">
		<Layout.HeaderViewer>READ</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet iopsWrite()}
	<Layout.Header class="justify-center">
		<Layout.HeaderViewer>WRITE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
