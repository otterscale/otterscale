<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import type { Secret } from '$lib/api/application/v1/application_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_picker,
		name,
		namespace,
		type,
		immutable,
		labels,
		created_at
	};
</script>

{#snippet row_picker(table: Table<Secret>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Secret>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet namespace(column: Column<Secret>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.namespace()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet type(column: Column<Secret>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.type()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet immutable()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.immutable()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet labels()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.labels()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet created_at(column: Column<Secret>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.create_time()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
