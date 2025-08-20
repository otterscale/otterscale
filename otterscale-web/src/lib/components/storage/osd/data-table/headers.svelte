<script lang="ts">
	import type { OSD } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import { Header as RowPicker } from '$lib/components/custom/data-table/data-table-row-pickers';
	import Sorter from '$lib/components/custom/data-table/data-table-sorter.svelte';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		row_picker,
		id,
		name,
		state,
		osdUp,
		osdIn,
		exists,
		deviceClass,
		machine,
		placementGroupCount,
		usage,
		iops,
		actions
	};
</script>

{#snippet row_picker(table: Table<OSD>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<RowPicker {table} />
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

{#snippet osdUp(column: Column<OSD>)}
	<Layout.Header class="hidden" />
{/snippet}

{#snippet osdIn(column: Column<OSD>)}
	<Layout.Header class="hidden" />
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
