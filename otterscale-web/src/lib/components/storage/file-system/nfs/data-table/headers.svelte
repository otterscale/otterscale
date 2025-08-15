<script lang="ts" module>
	import type { Subvolume } from '$lib/api/storage/v1/storage_pb';
	import Sorter from '$lib/components/custom/data-table/data-table-sorter.svelte';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/header.svelte';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		_row_picker,
		name,
		poolName,
		usage,
		path,
		mode,
		createTime,
		Export,
		snapshot
	};
</script>

{#snippet _row_picker(table: Table<Subvolume>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<TableRowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Subvolume>)}
	<Layout.Header>
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet path(column: Column<Subvolume>)}
	<Layout.Header>
		<Layout.HeaderViewer>PATH</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet mode(column: Column<Subvolume>)}
	<Layout.Header>
		<Layout.HeaderViewer>MODE</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet poolName(column: Column<Subvolume>)}
	<Layout.Header>
		<Layout.HeaderViewer>POOL NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet usage(column: Column<Subvolume>)}
	<Layout.Header>
		<Layout.HeaderViewer>USAGE</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet Export(column: Column<Subvolume>)}
	<Layout.Header>
		<Layout.HeaderViewer>EXPORT</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet createTime(column: Column<Subvolume>)}
	<Layout.Header>
		<Layout.HeaderViewer>CREATE TIME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet snapshot()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>SNAPSHOT</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
