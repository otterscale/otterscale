<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import type { Table, Column } from '@tanstack/table-core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { Bucket } from './types';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		owner: owner,
		usedCapacity: usedCapacity,
		capacityLimit: capacityLimit,
		objects: objects,
		objectLimit: objectLimit
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

{#snippet usedCapacity(column: Column<Bucket>)}
	<Layout.Header>
		<Layout.HeaderViewer>USED CAPACITY</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet capacityLimit(column: Column<Bucket>)}
	<Layout.Header>
		<Layout.HeaderViewer>CAPACITY LIMIT</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet objects(column: Column<Bucket>)}
	<Layout.Header>
		<Layout.HeaderViewer>OBJECTS</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet objectLimit(column: Column<Bucket>)}
	<Layout.Header>
		<Layout.HeaderViewer>OBJECT LIMIT</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}
