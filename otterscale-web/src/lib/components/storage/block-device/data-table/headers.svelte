<script lang="ts" module>
	import type { Image } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import { Header as RowPicker } from '$lib/components/custom/data-table/data-table-row-pickers';
	import Sorter from '$lib/components/custom/data-table/data-table-sorter.svelte';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		row_picker,
		name,
		poolName,
		usage,
		snapshots,
		actions
	};
</script>

{#snippet row_picker(table: Table<Image>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Image>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet poolName(column: Column<Image>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>POOL NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet usage(column: Column<Image>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>USAGE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet snapshots(column: Column<Image>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>SNAPSHOTS</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions(column: Column<Image>)}{/snippet}
