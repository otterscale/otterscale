<script lang="ts" module>
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Column, Table } from '@tanstack/table-core';
	import type { User_Key } from '$gen/api/storage/v1/storage_pb';

	export const headers = {
		_row_picker: _row_picker,
		access: access
	};
</script>

{#snippet _row_picker(table: Table<User_Key>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}

{#snippet access(column: Column<User_Key>)}
	<Layout.Header>
		<Layout.HeaderViewer>ACCESSOR</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}
