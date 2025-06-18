<script lang="ts" module>
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Column, Table } from '@tanstack/table-core';
	import type { BlockImage } from './types';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		pool: pool,
		namespace: namespace,
		usage: usage,
		objects: objects,
		parent: parent,
		mirroring: mirroring,
		nextScheduledSnapshot: nextScheduledSnapshot
	};
</script>

{#snippet _row_picker(table: Table<BlockImage>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet name(column: Column<BlockImage>)}
	<Layout.Header>
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet pool(column: Column<BlockImage>)}
	<Layout.Header>
		<Layout.HeaderViewer>POOL</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet namespace(column: Column<BlockImage>)}
	<Layout.Header>
		<Layout.HeaderViewer>NAMESPACE</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet usage(column: Column<BlockImage>)}
	<Layout.Header>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>USAGE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet objects(column: Column<BlockImage>)}
	<Layout.Header>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer class="items-end">
			<p>OBJECTS</p>
			<p>SIZE</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet parent(column: Column<BlockImage>)}
	<Layout.Header>
		<Layout.HeaderViewer>PARENT</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet mirroring(column: Column<BlockImage>)}
	<Layout.Header>
		<Layout.HeaderViewer>MIRRORING</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet nextScheduledSnapshot(column: Column<BlockImage>)}
	<Layout.Header>
		<Layout.HeaderViewer>NEXT SCHEDULED SNAPSHOT</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}
