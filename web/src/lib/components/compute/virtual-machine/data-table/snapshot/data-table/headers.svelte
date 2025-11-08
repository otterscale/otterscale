<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import type { VirtualMachine_Snapshot } from '$lib/api/instance/v1/instance_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_picker,
		name,
		namespace,
		sourceName,
		phase,
		ready,
		createTime,
		actions
	};
</script>

{#snippet row_picker(table: Table<VirtualMachine_Snapshot>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<VirtualMachine_Snapshot>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet namespace(column: Column<VirtualMachine_Snapshot>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.namespace()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet sourceName(column: Column<VirtualMachine_Snapshot>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.source_name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet phase(column: Column<VirtualMachine_Snapshot>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.phase()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet ready(column: Column<VirtualMachine_Snapshot>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.ready()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet createTime(column: Column<VirtualMachine_Snapshot>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.create_time()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet actions()}{/snippet}
