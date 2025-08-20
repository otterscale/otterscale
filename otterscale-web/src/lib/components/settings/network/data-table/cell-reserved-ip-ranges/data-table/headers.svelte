<script lang="ts" module>
	import type { Network_IPRange } from '$lib/api/network/v1/network_pb';
	import { Layout, RowPickers, Sorter } from '$lib/components/custom/data-table';
	import type { Column, Table } from '@tanstack/table-core';

	export const messages = {
		type: 'type',
		startIp: 'startIp',
		endIp: 'endIp',
		comment: 'comment'
	};

	export const headers = {
		row_picker,
		type,
		startIp,
		endIp,
		comment
	};
</script>

{#snippet row_picker(table: Table<Network_IPRange>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<RowPickers.Header {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet type(column: Column<Network_IPRange>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>TYPE</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet startIp(column: Column<Network_IPRange>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>START</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet endIp(column: Column<Network_IPRange>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>END</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet comment(column: Column<Network_IPRange>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>COMMENT</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}
