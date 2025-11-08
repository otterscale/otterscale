<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import type { Service } from '../types.d';

	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_picker,
		name,
		type,
		clusterIp,
		ports,
		endpoints,
		actions
	};
</script>

{#snippet row_picker(table: Table<Service>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Service>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet type(column: Column<Service>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.type()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet clusterIp(column: Column<Service>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.cluster_ip()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet ports(column: Column<Service>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.ports()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet endpoints()}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.endpoints()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions()}{/snippet}
