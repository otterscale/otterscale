<script lang="ts" module>
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { Column, Table } from '@tanstack/table-core';

	export const messages = {
		fqdn_ip: 'fqdn_ip',
		powerState: 'powerState',
		status: 'status',
		cores_arch: 'cores_arch',
		ram: 'ram',
		disk: 'disk',
		storage: 'storage',
		tags: 'tags',
		scope: 'scope',
		actions: 'actions'
	};

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
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet fqdn_ip(column: Column<Machine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>
			<p>FQDN</p>
			<p>IP</p>
		</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet powerState(column: Column<Machine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>POWER</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet status(column: Column<Machine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>STATUS</Layout.HeaderViewer>
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
			<p>CORES</p>
			<p>ARCH</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet ram(column: Column<Machine>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>RAM</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet disk(column: Column<Machine>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>DISK</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet storage(column: Column<Machine>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>STORAGE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet tags(column: Column<Machine>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>TAGS</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet scope(column: Column<Machine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>
			<p>SCOPE</p>
			<p>LAST COMMISSIONED</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions()}{/snippet}
