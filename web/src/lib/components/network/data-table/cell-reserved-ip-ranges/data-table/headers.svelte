<script lang="ts" module>
	import type { Network_IPRange } from '$lib/api/network/v1/network_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
	import type { Column, Table } from '@tanstack/table-core';

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
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}
{#snippet type(column: Column<Network_IPRange>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.type()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet startIp(column: Column<Network_IPRange>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.start_ip()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet endIp(column: Column<Network_IPRange>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.end_ip()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet comment(column: Column<Network_IPRange>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.comment()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}
