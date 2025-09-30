<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import type { VirtualMachine } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_picker,
		name,
		status,
		namespace,
		machineId,
		instanceTypeName,
		clusterIp,
		disk,
		port,
		createTime,
		actions,
	};
</script>

{#snippet row_picker(table: Table<VirtualMachine>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<VirtualMachine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet status(column: Column<VirtualMachine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.status()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet namespace(column: Column<VirtualMachine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.namespace()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet machineId(column: Column<VirtualMachine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.machine()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet instanceTypeName(column: Column<VirtualMachine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.instance_type()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet clusterIp(column: Column<VirtualMachine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.cluster_ip()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet disk(column: Column<VirtualMachine>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.disk()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet port(column: Column<VirtualMachine>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.ports()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet createTime(column: Column<VirtualMachine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.create_time()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet actions()}{/snippet}
