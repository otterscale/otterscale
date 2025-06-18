<script lang="ts" module>
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Column, Table } from '@tanstack/table-core';
	import type { Role } from './types';

	export const headers = {
		_row_picker: _row_picker,
		roleName: roleName,
		path: path,
		arn: arn,
		createDate: createDate,
		maximumSessionDuration: maximumSessionDuration
	};
</script>

{#snippet _row_picker(table: Table<Role>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet roleName(column: Column<Role>)}
	<Layout.Header>
		<Layout.HeaderViewer>ROLE NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet path(column: Column<Role>)}
	<Layout.Header>
		<Layout.HeaderViewer>PATH</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet arn(column: Column<Role>)}
	<Layout.Header>
		<Layout.HeaderViewer>ARN</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet createDate(column: Column<Role>)}
	<Layout.Header>
		<Layout.HeaderViewer>CREATE DATE</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet maximumSessionDuration(column: Column<Role>)}
	<Layout.Header>
		<Layout.HeaderViewer>MAXIMUM SESSION DURATION</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}
