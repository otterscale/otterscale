<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import type { Application_Service_Port } from '$lib/api/application/v1/application_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_picker,
		name,
		protocol,
		port,
		nodePort
	};
</script>

{#snippet row_picker(table: Table<Application_Service_Port>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Application_Service_Port>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet protocol(column: Column<Application_Service_Port>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.protocol()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet port(column: Column<Application_Service_Port>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.host_and_container()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet nodePort(column: Column<Application_Service_Port>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.node_port()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}
