<script lang="ts" module>
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Column, Table } from '@tanstack/table-core';
	import type { Image_Snapshot } from '$gen/api/storage/v1/storage_pb';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		protect: protect
	};
</script>

{#snippet _row_picker(table: Table<Image_Snapshot>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet name(column: Column<Image_Snapshot>)}
	<Layout.Header>
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet protect(column: Column<Image_Snapshot>)}
	<Layout.Header>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>PROTECTED</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
