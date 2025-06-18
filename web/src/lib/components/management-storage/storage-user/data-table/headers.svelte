<script lang="ts" module>
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Column, Table } from '@tanstack/table-core';
	import type { User } from './types';

	export const headers = {
		_row_picker: _row_picker,
		username: username,
		tenant: tenant,
		fullName: fullName,
		emailAddress: emailAddress,
		suspended: suspended,
		maximumBuckets: maximumBuckets,
		capacityLimit: capacityLimit,
		objectLimit: objectLimit
	};
</script>

{#snippet _row_picker(table: Table<User>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet username(column: Column<User>)}
	<Layout.Header>
		<Layout.HeaderViewer>USERNAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet tenant(column: Column<User>)}
	<Layout.Header>
		<Layout.HeaderViewer>TENANT</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet fullName(column: Column<User>)}
	<Layout.Header>
		<Layout.HeaderViewer>FULL NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet emailAddress(column: Column<User>)}
	<Layout.Header>
		<Layout.HeaderViewer>EMAIL ADDRESS</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet suspended(column: Column<User>)}
	<Layout.Header>
		<Layout.HeaderViewer>SUSPENDED</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet maximumBuckets(column: Column<User>)}
	<Layout.Header>
		<Layout.HeaderViewer>MAX BUCKETS</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet capacityLimit(column: Column<User>)}
	<Layout.Header>
		<Layout.HeaderViewer>CAPACITY LIMIT</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet objectLimit(column: Column<User>)}
	<Layout.Header>
		<Layout.HeaderViewer>OBJECT LIMIT</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}
