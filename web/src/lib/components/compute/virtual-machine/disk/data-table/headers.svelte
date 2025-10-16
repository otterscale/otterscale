<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import type { EnhancedDisk } from '$lib/components/compute/virtual-machine/units/type';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_picker,
		name,
		bus,
		bootOrder,
		dataVolume,
		type,
		phase,
		boot,
		size,
		actions,
		// data, // start script
	};
</script>

{#snippet row_picker(table: Table<EnhancedDisk>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<EnhancedDisk>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet bus(column: Column<EnhancedDisk>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.bus()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet bootOrder(column: Column<EnhancedDisk>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.boot_order()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet dataVolume(column: Column<EnhancedDisk>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.data_volume()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet type(column: Column<EnhancedDisk>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.type()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet phase(column: Column<EnhancedDisk>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.phase()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet boot(column: Column<EnhancedDisk>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.boot()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet size(column: Column<EnhancedDisk>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.size()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions()}{/snippet}
