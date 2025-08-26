<script lang="ts" module>
	import type { User } from '$lib/api/storage/v1/storage_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
	import type { Column, Table } from '@tanstack/table-core';

	export const messages = {
		id: m.id(),
		name: m.user(),
		suspended: m.suspended(),
		keys: m.key()
	};

	export const headers = {
		row_picker,
		id,
		name,
		suspended,
		keys,
		actions
	};
</script>

{#snippet row_picker(table: Table<User>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet id(column: Column<User>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.id()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<User>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet suspended(column: Column<User>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.suspended()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet keys(column: Column<User>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.key()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions(column: Column<User>)}{/snippet}
