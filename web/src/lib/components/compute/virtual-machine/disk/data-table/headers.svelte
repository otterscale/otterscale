<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import type { VirtualMachineDisk } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_picker,
		name,
		type,
		bus,
		source,
		sourceType,
		size,
		actions,
	};
</script>

{#snippet row_picker(table: Table<VirtualMachineDisk>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<VirtualMachineDisk>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet type(column: Column<VirtualMachineDisk>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.type()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet bus(column: Column<VirtualMachineDisk>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>bus</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet sourceType(column: Column<VirtualMachineDisk>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>Source Type</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet source(column: Column<VirtualMachineDisk>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.source()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet size(column: Column<VirtualMachineDisk>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>Size</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions()}{/snippet}
