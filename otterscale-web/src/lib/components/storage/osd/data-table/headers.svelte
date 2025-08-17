<script lang="ts" module>
	import type { OSD } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/header.svelte';
	import Sorter from '$lib/components/custom/data-table/data-table-sorter.svelte';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		_row_picker,
		id,
		name,
		state,
		stateUp,
		stateIn,
		exists,
		deviceClass,
		machine,
		placementGroupCount,
		usage,
		iops,
		actions
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
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>ID</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet state(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>STATE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet stateIn(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>IN</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet stateUp(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>UP</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet exists(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>EXISTS</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet deviceClass(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>DEVICE</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet machine(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>MACHINE</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet placementGroupCount(column: Column<OSD>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>PLACEMENT GROUP</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet usage(column: Column<OSD>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>USAGE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet iops(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>IOPS</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions(column: Column<OSD>)}{/snippet}
