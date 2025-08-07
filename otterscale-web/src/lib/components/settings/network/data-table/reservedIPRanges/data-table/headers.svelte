<script lang="ts" module>
	import type { Network_IPRange } from '$lib/api/network/v1/network_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/header.svelte';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		_row_picker: _row_picker,
		type,
		startIp,
		endIp,
		comment
	};
</script>

{#snippet _row_picker(table: Table<Network_IPRange>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<TableRowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet type(column: Column<Network_IPRange>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>TYPE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet startIp(column: Column<Network_IPRange>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>START</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet endIp(column: Column<Network_IPRange>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>END</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet comment(column: Column<Network_IPRange>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>COMMENT</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
