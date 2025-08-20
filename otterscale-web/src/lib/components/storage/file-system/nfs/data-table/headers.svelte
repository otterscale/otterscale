<script lang="ts" module>
	import type { Subvolume } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import { Header as RowPicker } from '$lib/components/custom/data-table/data-table-row-pickers';
	import Sorter from '$lib/components/custom/data-table/data-table-sorter.svelte';
	import type { Column, Table } from '@tanstack/table-core';

	export const messages = {
		name: 'name',
		poolName: 'poolName',
		usage: 'usage',
		path: 'path',
		mode: 'mode',
		createTime: 'createTime',
		exportSubvolume: 'exportSubvolume',
		snapshots: 'snapshots',
		actions: 'actions'
	};

	export const headers = {
		row_picker,
		name,
		poolName,
		usage,
		path,
		mode,
		createTime,
		exportSubvolume,
		snapshots,
		actions
	};
</script>

{#snippet row_picker(table: Table<Subvolume>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Subvolume>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet path(column: Column<Subvolume>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>PATH</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet mode(column: Column<Subvolume>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>MODE</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet poolName(column: Column<Subvolume>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>POOL NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet usage(column: Column<Subvolume>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>USAGE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet exportSubvolume(column: Column<Subvolume>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>EXPORT</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet createTime(column: Column<Subvolume>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>CREATE TIME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet snapshots(column: Column<Subvolume>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>SNAPSHOT</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions(column: Column<Subvolume>)}{/snippet}
