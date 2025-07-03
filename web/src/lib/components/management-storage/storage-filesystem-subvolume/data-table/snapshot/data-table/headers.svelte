<script lang="ts" module>
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Column, Table } from '@tanstack/table-core';
	import type { Subvolume_Snapshot } from '$gen/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		createTime: createTime,
		hasPendingClones: hasPendingClones
	};
</script>

{#snippet _row_picker(table: Table<Subvolume_Snapshot>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet name(column: Column<Subvolume_Snapshot>)}
	<Layout.Header>
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet hasPendingClones(column: Column<Subvolume_Snapshot>)}
	<Layout.Header>
		<Layout.HeaderViewer>PENDING CLONE</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet createTime(column: Column<Subvolume_Snapshot>)}
	<Layout.Header>
		<Layout.HeaderViewer>CREATE TIME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}
