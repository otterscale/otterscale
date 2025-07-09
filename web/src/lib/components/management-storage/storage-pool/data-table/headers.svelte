<script lang="ts" module>
	import type { Pool } from '$gen/api/storage/v1/storage_pb';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		applications: applications,
		placement_group_state: placement_group_state,
		usage: usage,
		readBytes: readBytes,
		writeBytes: writeBytes
	};
</script>

{#snippet _row_picker(table: Table<Pool>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet name(column: Column<Pool>)}
	<Layout.Header>
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet applications(column: Column<Pool>)}
	<Layout.Header>
		<Layout.HeaderViewer>APPLICATIONS</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet placement_group_state(column: Column<Pool>)}
	<Layout.Header>
		<Layout.HeaderViewer>PLACEMENT GROUP STATE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet usage(column: Column<Pool>)}
	<Layout.Header>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>USAGE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet readBytes(column: Column<Pool>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderViewer>READ</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet writeBytes(column: Column<Pool>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderViewer>WRITE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
