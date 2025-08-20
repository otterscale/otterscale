<script lang="ts" module>
	import type { Application } from '$lib/api/application/v1/application_pb';
	import { Sorter, Layout, RowPickers } from '$lib/components/custom/data-table';
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
		nodeport
	};
</script>

{#snippet row_picker(table: Table<Application>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<RowPickers.Header {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Application>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet type(column: Column<Application>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>TYPE</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet namespace(column: Column<Application>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>NAMESPACE</Layout.HeaderViewer>
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
		<Layout.HeaderViewer>HEALTH</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet service(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>SERVICE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet pod(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>POD</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet replica(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>REPLICA</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet container(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>CONTAINER</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet volume(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>VOLUME</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet nodeport()}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>NODEPORT</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
