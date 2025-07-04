<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import type { Table, Column } from '@tanstack/table-core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { Bucket } from '$gen/api/storage/v1/storage_pb';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		owner: owner,
		usage: usage,
		createTime: createTime
	};
</script>

{#snippet _row_picker(table: Table<Bucket>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet name(column: Column<Bucket>)}
	<Layout.Header>
		<Layout.HeaderViewer>BUCKET NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet owner(column: Column<Bucket>)}
	<Layout.Header>
		<Layout.HeaderViewer>OWNER</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet usage(column: Column<Bucket>)}
	<Layout.Header>
		<Layout.HeaderViewer class="w-full">
			<p class="text-end">USAGE</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet createTime(row: Column<Bucket>)}
	<Layout.Header>
		<Layout.HeaderViewer>CREATE TIME</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
