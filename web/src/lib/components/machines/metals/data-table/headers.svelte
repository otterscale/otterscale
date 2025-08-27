<script lang="ts" module>
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		row_picker,
		fqdn_ip,
		powerState,
		status,
		cores_arch,
		ram,
		disk,
		storage,
		tags,
		scope,
		actions
	};
</script>

{#snippet row_picker(table: Table<Machine>)}
	<Layout.Header class="items-start">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet fqdn_ip(column: Column<Machine>)}
	<Layout.Header class="items-start">
		<Layout.HeaderViewer>
			{m.fqdn()}
			<Layout.SubHeaderViewer>
				{m.ip()}
			</Layout.SubHeaderViewer>
		</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet powerState(column: Column<Machine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.power()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet status(column: Column<Machine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.status()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet cores_arch(column: Column<Machine>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>
			{m.core()}
			<Layout.SubHeaderViewer>
				{m.architecture()}
			</Layout.SubHeaderViewer>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet ram(column: Column<Machine>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.ram()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet disk(column: Column<Machine>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.disk()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet storage(column: Column<Machine>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.storage()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet tags(column: Column<Machine>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.tags()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet scope(column: Column<Machine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>
			{m.scope()}
			<Layout.SubHeaderViewer>
				{m.last_commissioned()}
			</Layout.SubHeaderViewer>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions()}{/snippet}
