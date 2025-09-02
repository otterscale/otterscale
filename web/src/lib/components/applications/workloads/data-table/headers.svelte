<script lang="ts" module>
	import type { Application } from '$lib/api/application/v1/application_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		row_picker,
		name,
		type,
		namespace,
		health,
		service,
		pod,
		replica,
		container,
		volume,
		nodeport,
	};
</script>

{#snippet row_picker(table: Table<Application>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Application>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet type(column: Column<Application>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.type()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet namespace(column: Column<Application>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.namespace()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet health(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.health()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet service(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.service()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet pod(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.pod()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet replica(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.replica()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet container(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.container()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet volume(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.volume()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet nodeport()}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.nodeport()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
