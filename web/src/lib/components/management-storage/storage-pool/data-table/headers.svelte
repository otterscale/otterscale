<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import type { Table, Column } from '@tanstack/table-core';
	import type { Pool } from './types';
	import * as Layout from '$lib/components/custom/data-table/layout';

	export const headers = {
		_row_picker: _row_picker,
		name: name,
		dataProtection: dataProtection,
		applications: applications,
		PGStatus: PGStatus,
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

{#snippet dataProtection(column: Column<Pool>)}
	<Layout.Header>
		<Layout.HeaderViewer>DATA PROTECTION</Layout.HeaderViewer>
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

{#snippet PGStatus(column: Column<Pool>)}
	<Layout.Header>
		<Layout.HeaderViewer>PG STATUS</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
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
	<Layout.Header>
		<Layout.HeaderViewer>READ</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet writeBytes(column: Column<Pool>)}
	<Layout.Header>
		<Layout.HeaderViewer>WRITE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
