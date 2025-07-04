<script lang="ts" module>
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Column, Table } from '@tanstack/table-core';
	import type { User } from '$gen/api/storage/v1/storage_pb';

	export const headers = {
		_row_picker: _row_picker,
		id: id,
		name: name,
		suspended: suspended,
		keys: keys
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

{#snippet id(column: Column<User>)}
	<Layout.Header>
		<Layout.HeaderViewer>ID</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<User>)}
	<Layout.Header>
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
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

{#snippet keys()}
	<Layout.Header>
		<Layout.HeaderViewer class="w-full">
			<p class="text-end">KEYS</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
