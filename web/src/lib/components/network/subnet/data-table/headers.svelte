<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import type { Network } from '$lib/api/network/v1/network_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_picker,
		fabric,
		vlan,
		dhcpOn,
		subnet,
		ipAddresses,
		ipRanges,
		statistics,
		actions,
	};
</script>

{#snippet row_picker(table: Table<Network>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet fabric(column: Column<Network>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.fabric()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet vlan(column: Column<Network>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.vlan()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet dhcpOn(column: Column<Network>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.dhcp_on()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet subnet(column: Column<Network>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.subnet()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet ipAddresses(column: Column<Network>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.ip_address()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet ipRanges(column: Column<Network>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.ip_range()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet statistics(column: Column<Network>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.statistics()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions()}{/snippet}
